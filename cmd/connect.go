package cmd

import (
	"fmt"

	"github.com/AllanM007/snap/internal/config"
	"github.com/AllanM007/snap/internal/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: set logrus logging level using --verbose tag; default to warn-level, human-readable output
var connectDbCmd = &cobra.Command{
	Use:     "connect",
	Short:   "enter db connection info",
	Long:    "",
	Example: "snap connect -name testdb -server 8.8.8.8 -port 5432 -user admin -password 12345",
	// ValidArgs: []string{"postgresql", "mysql", "sqlite"},
	// Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		dBName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		dbServer, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}
		dbPort, err := cmd.Flags().GetString("port")
		if err != nil {
			return err
		}
		dbUser, err := cmd.Flags().GetString("user")
		if err != nil {
			return err
		}
		dbPassword, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		//TODO: ping before config settings save
		viper.WriteConfigAs("config.env")

		newProf := config.Profile{
			Name:     dBName,
			Server:   dbServer,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
		}

		err = engine.Ping(newProf)
		if err != nil {
			return err
		}

		fmt.Println("Database ping successful!!")
		return nil
	},
}
