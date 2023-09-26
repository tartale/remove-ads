package config

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/stringz"
	"github.com/tartale/go/pkg/structs"
)

var (
	Values values
)

type values struct {
	LogLevel            string        `mapstructure:"RMADS_LOG_LEVEL" default:"INFO"`
	TempDir             string        `mapstructure:"RMADS_TEMP_DIR" default:"${PWD}/tmp"`
	StillFramesInterval time.Duration `mapstructure:"RMADS_THUMBNAIL_INTERVAL" default:"5s"`

	FFmpegFilePath  string `default:"" validate:"required,file"`
	FFprobeFilePath string `default:"" validate:"required,file"`
	SkipFilePath    string `default:"" validate:"required,file"`
	InputFilePath   string `default:"" validate:"required,file"`
	OutputFilePath  string `default:""`
}

func (v *values) SetDefaults() {

	defaults.SetDefaults(v)
}

func (v *values) ResolveVariables() error {

	var err error
	Values.FFmpegFilePath, err = exec.LookPath("ffmpeg")
	if err != nil || !filez.Exists(Values.FFmpegFilePath) {
		return fmt.Errorf("%s ffmpeg must be installed and in the PATH", errorz.ErrFatal)
	}
	Values.FFprobeFilePath, err = exec.LookPath("ffprobe")
	if err != nil || !filez.Exists(Values.FFmpegFilePath) {
		return fmt.Errorf("%s ffprobe must be installed and in the PATH", errorz.ErrFatal)
	}

	err = structs.Walk(&Values, func(sf reflect.StructField, sv reflect.Value) error {

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
