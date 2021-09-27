package config

import (
	"flag"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configEnvPrefix = "mtz"

// Config contiene la configuración del servicio
type Config struct {
	*viper.Viper
}

// NewConfig genera la configuración del servicio a partir de los argumentos
// de línea de comandos y variables de entorno
func NewConfig(flags *flag.FlagSet) *Config {
	vp := viper.New()

	pflag.CommandLine.AddGoFlagSet(flags)
	pflag.Parse()

	err := vp.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}

	vp.SetEnvPrefix(configEnvPrefix)

	replacer := strings.NewReplacer(".", "_")
	vp.SetEnvKeyReplacer(replacer)
	vp.AutomaticEnv()

	return &Config{vp}
}
