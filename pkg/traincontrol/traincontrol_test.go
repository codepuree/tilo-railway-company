package traincontrol

import "testing"

func TestTrainControl_interpret(t *testing.T) {
	tc := &TrainControl{
		Sensors:  make(map[int]Sensor),
		Switches: make(map[string]Switch),
		Signals:  make(map[string]Signal),
		Blocks:   make(map[string]Block),
	}

	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"Incomplete message", "01l", true},
		{"empty message", "", true},
		{"Sensor", "01lz", false},
		{"Switch", "ya0z", false},
		{"Signal", "xa1z", false},
		{"Block direction", "adfz", false},
		{"Block speed", "a50z", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tc.interpret(tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("TrainControl.interpret(\"%s\") error = %v, wantErr %v", tt.msg, err, tt.wantErr)
			}
		})
	}
}
