/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name     string
	server   string
	port     string
	user     string
	password string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "Snap is a database agnostic database backup cli tool",
	Long: `
	Snap is a CLI tool built using Go that allows users to do 
	quick terminal backups and restores to their database whether 
	local, hosted or containerized. `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snap.yaml)")
	connectDbCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "database name")
	viper.BindPFlag("name", connectDbCmd.PersistentFlags().Lookup("name"))
	connectDbCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "database host")
	viper.BindPFlag("server", connectDbCmd.PersistentFlags().Lookup("server"))
	connectDbCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "database port")
	viper.BindPFlag("port", connectDbCmd.PersistentFlags().Lookup("port"))
	connectDbCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "database username")
	viper.BindPFlag("user", connectDbCmd.PersistentFlags().Lookup("user"))
	connectDbCmd.PersistentFlags().StringVarP(&password, "password", "w", "", "database password")
	viper.BindPFlag("password", connectDbCmd.PersistentFlags().Lookup("password"))

	restoreDbCmd.Flags().StringVarP(&name, "name", "n", "", "database name")

	sshConnect.PersistentFlags().StringVarP(&server, "host", "H", "", "ssh host")
	sshConnect.MarkFlagRequired("host")
	sshConnect.PersistentFlags().StringVarP(&user, "user", "u", "", "user")
	sshConnect.MarkFlagRequired("user")
	sshConnect.PersistentFlags().StringVarP(&password, "key", "k", "", "path to ssh key file")
	sshConnect.MarkFlagRequired("key")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(versionCmd, connectDbCmd, backupDbCmd, restoreDbCmd, sshConnect)
	// selectDbCmd.AddCommand(connectDbCmd, backupDbCmd, restoreDbCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
