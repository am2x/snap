package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Database string

const (
	Postgres Database = "postgresql"
	Mysql    Database = "mysql"
	SQLite   Database = "sqlite"
)

// var application Application

var selectDbCmd = &cobra.Command{
	Use:   "db [database]",
	Short: "select database",
	Long: `Choose between the currently supported databases snap supports which are (postgresql, mysql, sqlite).
	- You can enter db-support to check the supported databases`,
	Example:   "snap db postgresql",
	ValidArgs: []string{"postgresql", "mysql", "sqlite"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database := args[0]
		fmt.Println("Database is", database)
		if database == string(Postgres) {

		} else if database == string(Mysql) {

		} else if database == string(SQLite) {

		} else {
			fmt.Println("Invalid database type !!")
		}
	},
}

var connectDb1Cmd = &cobra.Command{
	Use:     "connect [host, port, name, username, password]",
	Short:   "enter db connection info",
	Long:    "",
	Example: "snap -host 17.198.0.52 -port 5432 -name testdb -username admin -password 12345 ",
	// ValidArgs: []string{"postgresql", "mysql", "sqlite"},
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Args are ", strings.Join(args, ""))
	},
}

var backupDbCmd = &cobra.Command{
	Use:     "backup",
	Short:   "confirm to begin backup",
	Long:    "",
	Example: "snap backup testdb",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Args are ", strings.Join(args, ""))
	},
}

var connectDb2Cmd = &cobra.Command{
	Use:     "connect [host, port, name, username, password]",
	Short:   "enter db connection info",
	Long:    "",
	Example: "snap -host 17.198.0.52 -port 5432 -name testdb -username admin -password 12345 ",
	// ValidArgs: []string{"postgresql", "mysql", "sqlite"},
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Args are ", strings.Join(args, ""))
	},
}

var restoreDbCmd = &cobra.Command{
	Use:     "restore",
	Short:   "restore destination database",
	Long:    "",
	Example: "snap restore testdb",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Args are ", strings.Join(args, ""))
	},
}
