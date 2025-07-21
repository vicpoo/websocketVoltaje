package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/vicpoo/websocketVoltaje/Voltaje/application"
	"github.com/vicpoo/websocketVoltaje/Voltaje/domain/entities"
	"github.com/vicpoo/websocketVoltaje/repository"
)

type MessagingService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	hub  *Hub
}

func NewMessagingService(hub *Hub) *MessagingService {
	conn, err := amqp.Dial("amqp://reyhades:reyhades@44.219.123.4:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return nil
	}

	err = ch.ExchangeDeclare(
		"amq.topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
		return nil
	}

	return &MessagingService{
		conn: conn,
		ch:   ch,
		hub:  hub,
	}
}

func (ms *MessagingService) ConsumeVoltajeMessages() error {
	q, err := ms.ch.QueueDeclare(
		"sensor_voltaje", // cola para voltaje
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	err = ms.ch.QueueBind(
		q.Name,
		"sensor_volt", // routing key para voltaje
		"amq.topic",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ms.ch.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Inicializa el repositorio y use case para voltaje
	repo := repository.NewVoltajeRepositoryMySQL()
	useCase := application.NewVoltajeUseCase(repo)

	go func() {
		for msg := range msgs {
			log.Printf("Mensaje sensor_volt recibido: %s", string(msg.Body))

			var payload struct {
				Sensor         string  `json:"sensor"`
				Voltaje        float64 `json:"voltaje"`
				Corriente      float64 `json:"corriente"`
				Potencia       float64 `json:"potencia"`
				UnidadVoltaje  string  `json:"unidad_voltaje"`
				UnidadCorriente string `json:"unidad_corriente"`
				UnidadPotencia string  `json:"unidad_potencia"`
				Timestamp      int64   `json:"timestamp"`
				Ubicacion      string  `json:"ubicacion"`
			}

			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("Error al parsear mensaje JSON: %v", err)
				msg.Nack(false, false)
				continue
			}

			data := entities.NewVoltajeData(
				payload.Sensor,
				payload.Voltaje,
				payload.Corriente,
				payload.Potencia,
				payload.UnidadVoltaje,
				payload.UnidadCorriente,
				payload.UnidadPotencia,
				payload.Timestamp,
				payload.Ubicacion,
			)

			if err := useCase.SaveVoltajeData(*data); err != nil {
				log.Printf("Error al guardar en BD: %v", err)
			} else {
				log.Printf("Datos de voltaje guardados correctamente")
			}

			// Enviar a WebSocket
			ms.hub.broadcast <- msg.Body
			msg.Ack(false)
		}
	}()

	return nil
}

func (ms *MessagingService) Close() {
	if ms.ch != nil {
		ms.ch.Close()
	}
	if ms.conn != nil {
		ms.conn.Close()
	}
}
