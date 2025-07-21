package domain

import "github.com/vicpoo/websocketVoltaje/Voltaje/domain/entities"

type VoltajeRepository interface {
	Save(data entities.VoltajeData) error
	GetAll() ([]entities.VoltajeData, error)
}
