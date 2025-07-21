package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vicpoo/websocketVoltaje/Voltaje/infrastructure"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Inicializa el hub WebSocket
	hub := infrastructure.NewHub()
	go hub.Run()

	// Servicio de mensajería para voltaje (RabbitMQ)
	messagingService := infrastructure.NewMessagingService(hub)
	defer messagingService.Close()

	// Rutas WebSocket para voltaje
	infrastructure.SetupVoltajeRoutes(r, hub)

	// Inicia el consumidor de mensajes de voltaje
	if err := messagingService.ConsumeVoltajeMessages(); err != nil {
		log.Fatalf("Failed to start Voltaje consumer: %v", err)
	}

	// Captura de señales para apagado controlado
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ejecutar el servidor
	go func() {
		if err := r.Run(":8004"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("🔌 Server started on port 8004")
	log.Println("⚡ Voltaje RabbitMQ consumer started")

	<-sigChan
	log.Println("Shutting down voltaje server...")
}
