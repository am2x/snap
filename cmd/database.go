package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	port     string
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
	Example: "snap connect -name testdb -server 8.8.8.8 -port 5432 -user admin -password 12345 ",
	// ValidArgs: []string{"postgresql", "mysql", "sqlite"},
	// Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dbName, _ := cmd.Flags().GetString("name")
		dbServer, _ := cmd.Flags().GetString("server")
		dbPort, _ := cmd.Flags().GetString("port")
		dbUser, _ := cmd.Flags().GetString("user")
		dbPassword, _ := cmd.Flags().GetString("password")

		logrus.Warn(viper.Get("name"))

		// Define the command without the redirection
		osCommand := exec.Command("psql", "-U", dbUser, "-d", dbName, "-h", dbServer, "-p", dbPort, "-c", "\\q")
		osCommand.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

		//sample ssh command
		// osCommand := exec.Command("ssh", "-i", "Documents/file.pem", "ubuntu@176.48.03.79")

		logrus.Info(osCommand)
		// Execute the command
		err := osCommand.Run()
		if err != nil {
			log.Fatalf("Command execution failed: %s", err)
		}

		fmt.Println("Database ping successful!!")

		//SSL terminal access
		/*
					psql "sslmode=require host=<db-address> dbname=<db-name> user=<username>"

			        //example
			        psql "sslmode=require host=my-psql-db.cloud.neon.tech dbname=tutorials_db user=admin"
		*/

	},
}

/*
CREATE TABLE users (username VARCHAR(255), password VARCHAR(255));
*/
var backupDbCmd = &cobra.Command{
	Use:     "backup",
	Short:   "confirm to begin backup",
	Long:    "",
	Example: "snap backup -name testdb -user admin -password 12345 ",
	// Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		dbName := viper.Get("name")
		dbUser := viper.Get("user")
		dbPassword := viper.Get("password")

		logrus.Warn(viper.Get("name"))

		// Define the command without the redirection
		osCommand := exec.Command("pg_dump", "-U", dbUser.(string), "-d", dbName.(string))
		osCommand.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

		filePath := fmt.Sprintf("backups/%s.sql", dbName)

		// Create or open the file for writing the output
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create backup file: %s", err)
		}
		defer outputFile.Close()

		// Set the command's stdout to the file
		osCommand.Stdout = outputFile

		logrus.Info(osCommand)
		// Execute the command
		err = osCommand.Run()
		if err != nil {
			log.Fatalf("Command execution failed: %s", err)
		}

		fmt.Println("Backup successful")
	},
}

var restoreDbCmd = &cobra.Command{
	Use:     "restore",
	Short:   "restore destination database",
	Long:    "",
	Example: "snap restore -n testdb ",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dbName, _ := cmd.Flags().GetString("name")

		fileName := fmt.Sprintf("backups/%s.sql", dbName)

		// Create the PostgreSQL restore command
		// osCommand := exec.Command("pg_restore", "-h", dbServer, "-U", dbUser, "-d", dbName, "-1", fileName)
		osCommand := exec.Command("psql", "-d", dbName, "-f", fileName)

		// Run the command
		if err := osCommand.Run(); err != nil {
			logrus.Error("failed to run psql restore: %w", err)
			return
		}
		fmt.Println("Database", fileName, "backup restored:", dbName)
	},
}
