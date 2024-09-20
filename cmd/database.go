package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Database string

const (
	Postgres Database = "postgresql"
	Mysql    Database = "mysql"
	SQLite   Database = "sqlite"
)

var (
	name     string
	server   string
	port     int
	user     string
	password string
)

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
			log.Println("db is postgres")
		} else if database == string(Mysql) {
			log.Println("db is mysql")
		} else if database == string(SQLite) {
			log.Println("db is sqlite")
		} else {
			fmt.Println("Invalid database type !!")
		}
	},
}

var connectDbCmd = &cobra.Command{
	Use:     "connect",
	Short:   "enter db connection info",
	Long:    "",
	Example: "snap connect -name testdb -server 8.8.8.8 -password 5432 -user admin -password 12345 ",
	// ValidArgs: []string{"postgresql", "mysql", "sqlite"},
	// Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbName, _ := cmd.Flags().GetString("name")
		dbServer, _ := cmd.Flags().GetString("server")
		dbPort, _ := cmd.Flags().GetInt("port")
		dbUser, _ := cmd.Flags().GetString("user")
		dbPassword, _ := cmd.Flags().GetString("password")

		osCommand := exec.Command("export PGPASSWORD=" + dbPassword + "; psql -h " + dbServer + " -p " + string(dbPort) + " -U " + dbUser + " -d " + dbName)

		if err := osCommand.Run(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Database connected and live")

		//SSL terminal access
		/*
					psql "sslmode=require host=<db-address> dbname=<db-name> user=<username>"

			        //example
			        psql "sslmode=require host=my-psql-db.cloud.neon.tech dbname=tutorials_db user=admin"
		*/

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

		dbName, _ := cmd.Flags().GetString("name")
		// dbServer, _ := cmd.Flags().GetString("server")
		// dbPort, _ := cmd.Flags().GetInt("port")
		dbUser, _ := cmd.Flags().GetString("user")
		dbPassword, _ := cmd.Flags().GetString("password")

		fileName := fmt.Sprintf(dbName, ".sql")
		// Create the MySQL dump command
		osCommand := exec.Command("psql", "-u"+dbUser, "-p"+dbPassword, dbName)

		// Redirect the output to a file
		outfile, err := exec.Command("bash", "-c", fmt.Sprintf(" > %s", fileName)).Output()
		if err != nil {
			fmt.Errorf("failed to create backup file: %w", err)
			return
		}

		// Run the command
		if err := osCommand.Run(); err != nil {
			fmt.Errorf("failed to run mysqldump: %w", err)
			return
		}
		fmt.Println("Database", fileName, "backup created:", outfile)
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
