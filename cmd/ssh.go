package cmd

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sshConnect = &cobra.Command{
	Use:     "ssh-connect",
	Short:   "secure connect to remote server",
	Long:    "",
	Example: "snap ssh-connect --user centos --host 176.48.03.79 --key /Users/mac/key --verbose",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
		user, err := cmd.Flags().GetString("user")
		if err != nil {
			return err
		}
		keyFilePath, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}

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
			return err
		}
		return nil
	},
}
