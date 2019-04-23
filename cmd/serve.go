package cmd

import (
	"authapp/db"
	"authapp/server"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := db.OpenDB(); err != nil {
			log.Error("Failed to open DB")
			return
		}
		s, _ := server.NewServer()
		s.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().Int("port", 8080, "port for binding")
	viper.BindPFlag("config.port", serveCmd.PersistentFlags().Lookup("port"))
	viper.RegisterAlias("port", "config.port")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
