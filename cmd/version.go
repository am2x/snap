package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current project version",
	Long:  "A longer version is returned reflecting the current project release version available",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Snap version 0.0.1 -- HEAD")
	},
}

var versionConnect = &cobra.Command{
	Use:     "ssh-connect",
	Short:   "secure connect to remote server",
	Long:    "",
	Example: "snap ssh-connect --user centos --host 176.48.03.79 --key /Users/mac/key --verbose",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	}}
