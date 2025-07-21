package infrastructure

import (
	"github.com/gin-gonic/gin"
)

func SetupVoltajeRoutes(r *gin.Engine, hub *Hub) {
	// Ruta WebSocket específica para los datos de voltaje
	r.GET("/ws/voltaje", hub.HandleWebSocket)
}
