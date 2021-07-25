// Code generated by 'yaegi extract github.com/codepuree/tilo-railway-company/pkg/traincontrol'. DO NOT EDIT.

package trclib

import (
	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"go/constant"
	"go/token"
	"reflect"
)

var Symbols map[string]map[string]reflect.Value

func init() {
    Symbols = make(map[string]map[string]reflect.Value)
	Symbols["github.com/codepuree/tilo-railway-company/pkg/traincontrol/traincontrol"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Backward":         reflect.ValueOf(constant.MakeFromLiteral("98", token.INT, 0)),
		"Bent":             reflect.ValueOf(constant.MakeFromLiteral("49", token.INT, 0)),
		"EmergencyStopped": reflect.ValueOf(constant.MakeFromLiteral("120", token.INT, 0)),
		"Forward":          reflect.ValueOf(traincontrol.Forward),
		"NewTrainControl":  reflect.ValueOf(traincontrol.NewTrainControl),
		"Stopped":          reflect.ValueOf(constant.MakeFromLiteral("115", token.INT, 0)),
		"Straight":         reflect.ValueOf(traincontrol.Straight),

		// type definitions
		"Block":              reflect.ValueOf((*traincontrol.Block)(nil)),
		"BlockDirection":     reflect.ValueOf((*traincontrol.BlockDirection)(nil)),
		"Message":            reflect.ValueOf((*traincontrol.Message)(nil)),
		"Sensor":             reflect.ValueOf((*traincontrol.Sensor)(nil)),
		"Signal":             reflect.ValueOf((*traincontrol.Signal)(nil)),
		"Switch":             reflect.ValueOf((*traincontrol.Switch)(nil)),
		"SwitchState":        reflect.ValueOf((*traincontrol.SwitchState)(nil)),
		"Track":              reflect.ValueOf((*traincontrol.Track)(nil)),
		"Train":              reflect.ValueOf((*traincontrol.Train)(nil)),
		"TrainControl":       reflect.ValueOf((*traincontrol.TrainControl)(nil)),
		"TrainControlConfig": reflect.ValueOf((*traincontrol.TrainControlConfig)(nil)),

		// interface wrapper definitions
		"_Message": reflect.ValueOf((*_github_com_codepuree_tilo_railway_company_pkg_traincontrol_Message)(nil)),
	}
}

// _github_com_codepuree_tilo_railway_company_pkg_traincontrol_Message is an interface wrapper for Message type
type _github_com_codepuree_tilo_railway_company_pkg_traincontrol_Message struct {
	IValue interface{}
}
