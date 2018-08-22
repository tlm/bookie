package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bookie",
	Short: "bookie is a dns update addon for kubernetes",
	Long: `bookie is a kubernetes controller that watches ingress resources
			and updates DNS records`,
	RunE: execute,
}

func execute(cmd *cobra.Command, args []string) error {
	fmt.Println("here")
	return nil
}
