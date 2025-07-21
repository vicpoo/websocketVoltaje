package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/websocketVoltaje/Voltaje/application"
)

type VoltajeController struct {
	useCase *application.VoltajeUseCase
}

func NewVoltajeController(useCase *application.VoltajeUseCase) *VoltajeController {
	return &VoltajeController{useCase: useCase}
}

func (vc *VoltajeController) GetAll(c *gin.Context) {
	data, err := vc.useCase.GetAllVoltajeData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
