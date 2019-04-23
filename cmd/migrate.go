package cmd

import (
	"authapp/db"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// migrateCmd migrates DB
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates DB",
	Long:  `Migrates DB to required version`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("DB migration")
		if _, err := db.OpenDB(); err != nil {
			log.Error("Failed to open DB")
			return
		}
		err := db.Create()
		if err != nil {
			log.Error("Failed to migrate DB")
		} else {
			log.Info("DB migration finished successfully")
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
