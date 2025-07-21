package application

import (
	"github.com/vicpoo/websocketVoltaje/Voltaje/domain"
	"github.com/vicpoo/websocketVoltaje/Voltaje/domain/entities"
)

type VoltajeUseCase struct {
	repo domain.VoltajeRepository
}

func NewVoltajeUseCase(repo domain.VoltajeRepository) *VoltajeUseCase {
	return &VoltajeUseCase{repo: repo}
}

func (uc *VoltajeUseCase) SaveVoltajeData(data entities.VoltajeData) error {
	return uc.repo.Save(data)
}

func (uc *VoltajeUseCase) GetAllVoltajeData() ([]entities.VoltajeData, error) {
	return uc.repo.GetAll()
}
