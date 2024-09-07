package cmd

import "github.com/spf13/cobra"

var selectDBCmd = &cobra.Command{
	Use:   "db",
	Short: "select database",
	Long: `Choose between the currently supported databases snap supports which are (postgresql, mysql, sqlite).
	- You can enter db-support to check the supported databases`,
	Example:   "snap -db postgresql",
	ValidArgs: []string{"postgresql", "mysql", "sqlite"},
}
