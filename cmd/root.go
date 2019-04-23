package cmd

import (
	"os"

	"authapp/logging"

	log "github.com/sirupsen/logrus"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "chat-app",
	Short: "Simple chat platform",
	Long:  "Chat platform for people communication",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.json)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.SetDefault("config", "")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("db.url", "postgres://postgres:postgres@localhost:5432/chatapp?sslmode=disable")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
	}

	viper.SetDefault("config.port", 8080)

	viper.SetEnvPrefix("authapp")
	//viper.BindEnv("port")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(logrus.Fields{"err": err}).Error("Error reading config file")
	} else {
		log.WithFields(logrus.Fields{"file": viper.ConfigFileUsed()}).Info("Config file loaded")
	}

	logging.InitDefaultLogger()
}
