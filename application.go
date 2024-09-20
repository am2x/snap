package main

import (
	"context"
	"fmt"
	"os/exec"
)

type ServerConfig struct {
	// DbType    Database
	Name      string
	Host      string
	Port      int
	User      string
	Password  string
	Container bool
}

type Application struct {
	context context.Context
	server  ServerConfig
}

const (
	MysqlPrefix    string = "mysql"
	PostgresPrefix string = "psql"
)

func (a *Application) SSHBackup() error {
	return nil
}

func (a *Application) SSHRestore() error {
	return nil
}

func LocalBackup(a *Application) error {
	fileName := fmt.Sprintf(a.server.Name, ".sql")
	// Create the MySQL dump command
	cmd := exec.Command("psql", "-u"+a.server.User, "-p"+a.server.Password, a.server.Name)

	// Redirect the output to a file
	outfile, err := exec.Command("bash", "-c", fmt.Sprintf(" > %s", fileName)).Output()
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run mysqldump: %w", err)
	}
	fmt.Println("Database", fileName, "backup created:", outfile)
	return nil
}

func (a *Application) LocalRestore() error {
	fileName := fmt.Sprintf(a.server.Name, ".sql")
	// Create the MySQL restore command
	cmd := exec.Command("mysql", "-u"+a.server.User, "-p"+a.server.Password, a.server.Name)

	// Redirect input from the file
	cmd.Stdin = exec.Command("bash", "-c", fmt.Sprintf(" < %s", fileName)).Stdin

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run mysql restore: %w", err)
	}
	fmt.Println("Database restored from:", fileName)
	return nil
}
