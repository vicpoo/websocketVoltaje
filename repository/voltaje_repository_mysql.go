package repository

import (
	"database/sql"
	"fmt"

	"github.com/vicpoo/websocketVoltaje/Voltaje/domain"
	"github.com/vicpoo/websocketVoltaje/Voltaje/domain/entities"
	"github.com/vicpoo/websocketVoltaje/core"
)

type voltajeRepositoryMySQL struct {
	db *sql.DB
}

func NewVoltajeRepositoryMySQL() domain.VoltajeRepository {
	return &voltajeRepositoryMySQL{
		db: core.GetBD(),
	}
}

func (r *voltajeRepositoryMySQL) Save(data entities.VoltajeData) error {
	var sensorID int
	err := r.db.QueryRow("SELECT id FROM sensors WHERE name = ?", data.Sensor).Scan(&sensorID)
	if err != nil {
		return fmt.Errorf("no se encontró el sensor '%s': %v", data.Sensor, err)
	}

	_, err = r.db.Exec(`
		INSERT INTO sensor_readings (
			sensor_id, voltage, current, potencia, recorded_at
		) VALUES (?, ?, ?, ?, FROM_UNIXTIME(?))`,
		sensorID, data.Voltage, data.Current, data.Power, data.Timestamp)

	if err != nil {
		return fmt.Errorf("error al insertar en sensor_readings: %v", err)
	}

	return nil
}

func (r *voltajeRepositoryMySQL) GetAll() ([]entities.VoltajeData, error) {
	rows, err := r.db.Query(`
		SELECT s.name, sr.voltaje, sr.corriente, sr.potencia, UNIX_TIMESTAMP(sr.recorded_at), s.location
		FROM sensor_readings sr
		JOIN sensors s ON sr.sensor_id = s.id
		WHERE sr.voltaje IS NOT NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.VoltajeData
	for rows.Next() {
		var d entities.VoltajeData
		var timestamp int64

		err := rows.Scan(&d.Sensor, &d.Voltage, &d.Current, &d.Power, &timestamp, &d.Location)
		if err != nil {
			return nil, err
		}
		d.Timestamp = timestamp
		result = append(result, d)
	}

	return result, nil
}
