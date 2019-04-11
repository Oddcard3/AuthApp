package config

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Init initialize config
func Init() {
	// default values
	viper.SetDefault("port", 8080)
	viper.SetDefault("config", "")
	viper.SetDefault("logger.level", "info")

	//aliases
	viper.RegisterAlias("port", "config.port")

	// flags
	flag.Int("port", 8080, "port for binding")
	flag.String("config", "", "path to config file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// env
	viper.SetEnvPrefix("authapp")
	viper.BindEnv("port")

	// config file
	if cfgFilePath := viper.GetString("config"); cfgFilePath != "" {
		viper.SetConfigFile(cfgFilePath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/authapp/")
		viper.AddConfigPath("$HOME/.authapp")
		viper.AddConfigPath("./")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading file %s", err)
	}
}
