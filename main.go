package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"webssh/config"
	"webssh/db"
	"webssh/handlers"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	handlers.SetJWTSecret(cfg.JWT.Secret)

	if err := db.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	staticFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/static", http.FS(staticFS))
	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("static/", http.FS(staticFS))
	})

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
	}

	authAPI := r.Group("/api")
	authAPI.Use(handlers.AuthMiddleware())
	{
		authAPI.POST("/change-password", handlers.ChangePassword) // 新增
	}

	wsAuth := r.Group("/api")
	wsAuth.Use(handlers.WSAuthMiddleware())
	{
		wsAuth.GET("/ws/ssh", handlers.SSHWebSocket)
	}

	addr := cfg.Server.Addr
	if cfg.Server.CertFile != "" && cfg.Server.KeyFile != "" {
		log.Printf("Starting HTTPS server on %s", addr)
		if err := r.RunTLS(addr, cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
			log.Fatalf("HTTPS server error: %v", err)
		}
	} else {
		log.Printf("Starting HTTP server on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}
}