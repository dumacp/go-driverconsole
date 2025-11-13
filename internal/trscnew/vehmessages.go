package app

import (
	"fmt"
	"strings"

	"github.com/dumacp/go-driverconsole/internal/ui"
)

type VehicleMessages struct {
	VehicleId string
	Message   string
}

// función para agragar mensajes a la cola de mensajes del vehículo
func (a *App) AddVehicleMessage(msg *VehicleMessages) {
	if msg == nil || len(msg.VehicleId) == 0 || len(msg.Message) == 0 {
		return
	}
	if a.vehicleMessages == nil {
		a.vehicleMessages = make([]string, 0, 10)
	}
	if len(a.vehicleMessages) > 9 {
		a.vehicleMessages = a.vehicleMessages[1:] // eliminar el mensaje más antiguo
	}
	a.vehicleMessages = append(a.vehicleMessages, msg.Message)
	fmt.Printf("Mensaje agregado al vehículo %s: %s\n", msg.VehicleId, msg.Message)
}

// función para mostrar los mensajes del vehículo
func (a *App) ShowVehicleMessages() error {
	if len(a.vehicleMessages) == 0 {
		return nil // no hay mensajes para mostrar
	}

	// ordenar el slice de mensajes en orden inverso en un nuevo slice
	messages := make([]string, len(a.vehicleMessages))
	for i := 0; i < len(a.vehicleMessages); i++ {
		messages[i] = a.vehicleMessages[len(a.vehicleMessages)-1-i]
	}
	// completar el slice con espacios con string de espacios en blanco
	for i := len(messages); i < Label2DisplayRegister(ui.NOTIFICATIONS_VEHI_TEXT).Len; i++ {
		spaces := strings.Repeat(" ", Label2DisplayRegister(ui.NOTIFICATIONS_VEHI_TEXT).Size)
		messages = append(messages, spaces)
	}

	// mostrar los mensajes en la consola
	if err := a.uix.WriteTextRawDisplay(
		ui.NOTIFICATIONS_VEHI_TEXT,
		messages,
	); err != nil {
		return err
	}
	return nil
}
