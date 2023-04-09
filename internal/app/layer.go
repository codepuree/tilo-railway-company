package app

import (
	"fmt"
	"log"

	"github.com/codepuree/tilo-railway-company/internal/communication"
	"github.com/codepuree/tilo-railway-company/pkg/scriptcontrol"
	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

type Layer struct {
	Arduino       *Arduino
	TrainControl  *traincontrol.TrainControl
	ScriptControl *scriptcontrol.ScriptControl //map[string]scriptcontrol.Func

	trainControlConfig *traincontrol.TrainControlConfig
}

func NewLayer(port string, baudRate int, tc *traincontrol.TrainControlConfig, scriptFolder string) (Layer, error) {
	var layer Layer

	// Arduino
	a := NewArduino(port, baudRate)
	layer.Arduino = a

	// Train Control
	layer.TrainControl = traincontrol.NewTrainControl(tc)

	// Script Control
	layer.ScriptControl = scriptcontrol.NewScriptControl(scriptFolder)

	return layer, nil
}

func (lyr *Layer) Bind(evtSys *EventSystem, ws *communication.Publisher) error {
	err := lyr.Arduino.Connect()
	if err != nil {
		log.Println(fmt.Errorf("unable to connect to Arduino: %w", err))
		// return fmt.Errorf("unable to connect to Arduino: %w", err) // TODO: enable
	}

	lyr.ScriptControl.Bind(*ws)
	err = lyr.ScriptControl.Load()
	if err != nil {
		return fmt.Errorf("unable to load scriptcontrol: %w", err)
	}

	return nil
}

func (lyr *Layer) Close() error {
	return lyr.ScriptControl.Close()
}
