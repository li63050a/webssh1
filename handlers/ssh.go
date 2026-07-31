package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SSHRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func SSHWebSocket(c *gin.Context) {
	wsConn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer wsConn.Close()

	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		return
	}
	var req SSHRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Invalid SSH parameters"))
		return
	}

	sshConfig := &ssh.ClientConfig{
		User: req.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(req.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", req.Host, req.Port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("SSH connection failed: "+err.Error()))
		return
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Session creation failed: "+err.Error()))
		return
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Request pty failed: "+err.Error()))
		return
	}

	stdinPipe, _ := session.StdinPipe()
	stdoutPipe, _ := session.StdoutPipe()
	stderrPipe, _ := session.StderrPipe()

	if err := session.Shell(); err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Start shell failed: "+err.Error()))
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				return
			}
			stdinPipe.Write(msg)
		}
	}()

	go func() {
		defer wg.Done()
		outputReader := io.MultiReader(stdoutPipe, stderrPipe)
		buf := make([]byte, 1024)
		for {
			n, err := outputReader.Read(buf)
			if n > 0 {
				wsConn.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	session.Wait()
	wg.Wait()
}