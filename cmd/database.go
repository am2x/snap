package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
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

type DbConfig struct {
	Name     string `mapstructure:"NAME"`
	Server   string `mapstructure:"SERVER"`
	Port     string `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
}

func LoadConfig(path string) (config DbConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// var selectDbCmd = &cobra.Command{
// 	Use:   "db [database]",
// 	Short: "select database",
// 	Long: `Choose between the currently supported databases snap supports which are (postgresql, mysql, sqlite).
// 	- You can enter db-support to check the supported databases`,
// 	Example:   "snap db postgresql",
// 	ValidArgs: []string{"postgresql", "mysql", "sqlite"},
// 	Args:      cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		database := args[0]
// 		fmt.Println("Database is", database)
// 		if database == string(Postgres) {
// 			log.Println("db is postgres")
// 		} else if database == string(Mysql) {
// 			log.Println("db is mysql")
// 		} else if database == string(SQLite) {
// 			log.Println("db is sqlite")
// 		} else {
// 			fmt.Println("Invalid database type !!")
// 		}
// 	},
// }

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

		viper.WriteConfigAs("config.env")

		// Define the command without the redirection
		osCommand := exec.Command("psql", "-U", dbUser, "-d", dbName, "-h", dbServer, "-p", dbPort, "-c", "\\q")
		osCommand.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

		//sample ssh command
		// osCommand := exec.Command("ssh", "-i", "Documents/file.pem", "ubuntu@176.48.03.79")

		// Execute the command
		log.Info(osCommand)
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

		config, _ := LoadConfig(".")

		// Define the command without the redirection
		osCommand := exec.Command("pg_dump", "-U", config.User, "-d", config.Name)
		osCommand.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", config.Password))

		filePath := fmt.Sprintf("backups/%s.sql", config.Name)

		// Create or open the file for writing the output
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create backup file: %s", err)
		}
		defer outputFile.Close()

		// Set the command's stdout to the file
		osCommand.Stdout = outputFile

		log.Info(osCommand)
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
	// Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dbName, _ := cmd.Flags().GetString("name")

		filePath := fmt.Sprintf("backups/%s.sql", dbName)

		// Create the PostgreSQL restore command
		// osCommand := exec.Command("pg_restore", "-h", dbServer, "-U", dbUser, "-d", dbName, "-1", fileName)
		restoreCommand := exec.Command("psql", "-d", dbName, "-f", filePath)

		// Run the command
		log.Info(restoreCommand)
		output, err := restoreCommand.CombinedOutput()
		if err != nil {
			log.Error(string(output))

			if strings.Contains(string(output), "does not exist") {

				//create os pipe for running subsequent commands
				r, w, err := os.Pipe()
				if err != nil {
					log.Errorf("error opening os pipe : %v", err)
					return
				}
				defer r.Close()

				createCommand := exec.Command("psql", "-U", "postgres", "-c", fmt.Sprintf("CREATE DATABASE %s;", dbName))
				// log.Info(createCommand)
				createCommand.Stdout = w

				err = createCommand.Start()
				if err != nil {
					log.Errorf("error running create command: %v", err)
					return
				}
				defer createCommand.Wait()
				w.Close()

				restoreCommand.Stdin = r
				restoreCommand.Stdout = os.Stdout
				err = restoreCommand.Run()
				if err != nil {
					log.Errorf("error restoring database : %v", err)
					return
				}
			} else {
				log.Errorf("failed to run psql restore: %s", string(output))
				return
			}
		}
		log.Info("Database ", dbName, " backup restored")
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		viper.Reset()
	},
}
