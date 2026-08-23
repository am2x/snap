package engine

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/AllanM007/snap/internal/config"
	"github.com/AllanM007/snap/internal/db"
	"github.com/sirupsen/logrus"
)

func Ping(prof config.Profile) error {
	// Define the command without the redirection
	osCommand := exec.Command("psql", "-U", prof.User, "-d", prof.Name, "-h", prof.Server, "-p", prof.Port, "-c", "\\q")
	osCommand.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", prof.Password))

	// Execute the command
	err := osCommand.Run()
	if err != nil {
		log.Fatalf("Command execution failed: %s", err)
	}
	return nil

}

func Backup(prof config.Profile) error {
	// Define the command without the redirection
	osCommand := exec.Command("pg_dump", "-U", prof.User, "-d", prof.Name)
	osCommand.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", prof.Password))

	filePath := fmt.Sprintf("backups/%s.sql", prof.Name)

	// Create or open the file for writing the output
	outputFile, err := os.Create(filePath)
	if err != nil {
		logrus.Fatalf("Failed to create backup file: %s", err)
	}
	defer outputFile.Close()

	// Set the command's stdout to the file
	osCommand.Stdout = outputFile

	logrus.Info(osCommand)
	err = osCommand.Run()
	if err != nil {
		logrus.Fatalf("Command execution failed: %s", err)
	}
	return nil
}

func Restore(dbName string) error {
	filePath := fmt.Sprintf("backups/%s.sql", dbName)

	// Create the PostgreSQL restore command
	// osCommand := exec.Command("pg_restore", "-h", dbServer, "-U", dbUser, "-d", dbName, "-1", fileName)
	restoreCommand := exec.Command("psql", "-d", dbName, "-f", filePath)

	// Run the command
	logrus.Info(restoreCommand)
	output, err := restoreCommand.CombinedOutput()
	if err != nil {
		logrus.Error(string(output))

		if strings.Contains(string(output), "does not exist") {

			//create os pipe for running subsequent commands
			r, w, err := os.Pipe()
			if err != nil {
				logrus.Errorf("error opening os pipe : %v", err)
				return err
			}
			defer r.Close()

			dialect := db.GetDialect("")
			safeName := dialect.QuoteIdentifier(dbName)

			createCommand := exec.Command("psql", "-U", "postgres", "-c", fmt.Sprintf("CREATE DATABASE %s;", safeName))
			// log.Info(createCommand)
			createCommand.Stdout = w

			err = createCommand.Start()
			if err != nil {
				logrus.Errorf("error running create command: %v", err)
				return err
			}
			defer createCommand.Wait()
			w.Close()

			restoreCommand.Stdin = r
			restoreCommand.Stdout = os.Stdout
			err = restoreCommand.Run()
			if err != nil {
				logrus.Errorf("error restoring database : %v", err)
				return err
			}
		} else {
			logrus.Errorf("failed to run psql restore: %s", string(output))
			return err
		}
	}
	return nil
}
