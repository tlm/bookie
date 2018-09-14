package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/tlmiller/bookie/pkg/api"
	"github.com/tlmiller/bookie/pkg/engine"
	exconfig "github.com/tlmiller/bookie/pkg/executor/config"
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
	_ = api.NewServer()
	_ = engine.NewEngine()

	_, err := exconfig.ExecutorsForConfig()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "",
		"kubeconfig path")
	RootCmd.PersistentFlags().StringVar(&kubeConfig, "master", "",
		"kubernetes master")
}

func preRun(cmd *cobra.Command, args []string) {
	setupLogging()
}
