package cmd

import (
	"github.com/spf13/cobra"
)

var (
	kubeConfig string
	kubeMaster string
)

var RootCmd = &cobra.Command{
	Use:   "bookie",
	Short: "bookie is a dns update addon for kubernetes",
	Long: `bookie is a kubernetes controller that watches ingress resources
			and updates DNS records`,
	PersistentPreRun: preRun,
	Run:              execute,
}

func execute(cmd *cobra.Command, args []string) {
}

func init() {
	RootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "",
		"kubeconfig path")
	RootCmd.PersistentFlags().StringVar(&kubeConfig, "master", "",
		"kubernetes master")
}

func preRun(cmd *cobra.Command, args []string) {
	setupLogging()
}
