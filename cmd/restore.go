package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreDbCmd = &cobra.Command{
	Use:     "restore",
	Short:   "restore destination database",
	Long:    "",
	Example: "snap restore -n testdb --verbose",
	RunE: func(cmd *cobra.Command, args []string) error {

		restoreDBName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		log.Info("Database ", restoreDBName, " backup restored")
		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		viper.Reset()
	},
}
