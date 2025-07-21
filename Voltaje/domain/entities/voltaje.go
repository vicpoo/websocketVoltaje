package entities

import "time"

type VoltajeData struct {
	ID             int       `json:"id"`
	Sensor         string    `json:"sensor"`
	Voltage        float64   `json:"voltaje"`
	Current        float64   `json:"corriente"`
	Power          float64   `json:"potencia"`
	VoltageUnit    string    `json:"unidad_voltaje"`
	CurrentUnit    string    `json:"unidad_corriente"`
	PowerUnit      string    `json:"unidad_potencia"`
	Timestamp      int64     `json:"timestamp"`
	Location       string    `json:"ubicacion"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewVoltajeData(
	sensor string,
	voltage float64,
	current float64,
	power float64,
	voltageUnit string,
	currentUnit string,
	powerUnit string,
	timestamp int64,
	location string,
) *VoltajeData {
	return &VoltajeData{
		Sensor:      sensor,
		Voltage:     voltage,
		Current:     current,
		Power:       power,
		VoltageUnit: voltageUnit,
		CurrentUnit: currentUnit,
		PowerUnit:   powerUnit,
		Timestamp:   timestamp,
		Location:    location,
		CreatedAt:   time.Now(),
	}
}
