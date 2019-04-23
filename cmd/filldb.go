package cmd

import (
	"authapp/db"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// fillDBCmd fills DB with test data
var fillDBCmd = &cobra.Command{
	Use:   "filldb",
	Short: "Fills DB",
	Long:  `Fills DB with test data`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("DB filling")
		if _, err := db.OpenDB(); err != nil {
			log.Error("Failed to open DB")
			return
		}
		err := db.FillTestData()
		if err != nil {
			log.Error("Failed to fill DB")
		} else {
			log.Info("DB filled with test data successfully")
		}
	},
}

func init() {
	RootCmd.AddCommand(fillDBCmd)
}
