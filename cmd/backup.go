package cmd

import (
	"fmt"

	"github.com/AllanM007/snap/internal/config"
	"github.com/AllanM007/snap/internal/engine"
	"github.com/spf13/cobra"
)

/*
CREATE TABLE users (username VARCHAR(255), password VARCHAR(255));
*/
var backupDbCmd = &cobra.Command{
	Use:     "backup",
	Short:   "confirm to begin backup",
	Long:    "",
	Example: "snap backup -name testdb -user admin -password 12345 --verbose",
	Args:    cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {

		prof, err := config.LoadConfig(".")
		if err != nil {
			return err
		}

		err = engine.Backup(prof)
		if err != nil {
			return err
		}

		fmt.Println("Backup successful")
		return nil
	},
}
