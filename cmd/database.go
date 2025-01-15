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

		// Execute the command
		err := osCommand.Run()
		if err != nil {
			log.Fatalf("Command execution failed: %s", err)
		}

		fmt.Println("Database ping successful!!")
	},
}

var sshConnect = &cobra.Command{
	Use:     "ssh-connect",
	Short:   "secure connect to remote server",
	Long:    "",
	Example: "snap ssh-connect --user centos --host 176.48.03.79 --key /Users/mac/key",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		keyFilePath, _ := cmd.Flags().GetString("key")

		server := fmt.Sprintf("%s@%s", user, host)
		command := "docker ls"

		//ssh -i /Users/mac/key ubuntu@176.48.03.79
		sshCmd := exec.Command("ssh", "-i", keyFilePath, server, command)
		log.Info(sshCmd)

		// Set the output to the standard output and error
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Run(); err != nil {
			log.Errorf("error executing ssh command : %v\n", err)
			return
		}
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
