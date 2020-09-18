package custom

import "github.com/codepuree/tilo-railway-company/pkg/traincontrol"

func Track(tc *traincontrol.TrainControl) {
	tc.SetSwitch("a", "0")
	tc.SetSwitch("b", "0")
	tc.SetSwitch("e", "0")
	tc.SetSwitch("f", "0")
}
