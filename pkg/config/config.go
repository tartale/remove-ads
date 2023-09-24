package config

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/stringz"
	"github.com/tartale/go/pkg/structs"
)

var (
	Values values
)

type values struct {
	LogLevel string `mapstructure:"RMADS_LOG_LEVEL" default:"INFO"`
	DryRun   bool   `mapstructure:"RMADS_DRY_RUN" default:"false"`

	SkipFilePath   string `default:"" validate:"required,file"`
	InputFilePath  string `default:""`
	OutputFilePath string `default:""`
}

func (v *values) SetDefaults() {
	defaults.SetDefaults(v)
}

func (v *values) ResolveVariables() error {

	err := structs.Walk(&Values, func(sf reflect.StructField, sv reflect.Value) error {

		val := sv.Interface()
		err := stringz.Envsubst(&val)
		if err != nil && errors.Is(err, errorz.ErrInvalidType) {
			return nil
		}
		if err != nil {
			return err
		}
		tsv := sv.Type()
		vval := reflect.ValueOf(val)
		sv.Set(vval.Convert(tsv))

		return nil
	})

	return err
}

func (v *values) Validate() error {

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(v)
	if err != nil {
		return err
	}

	return nil
}
