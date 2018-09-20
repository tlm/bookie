package cmd

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	//"github.com/tlmiller/bookie/pkg/api"
	"github.com/tlmiller/bookie/pkg/engine"
	exconfig "github.com/tlmiller/bookie/pkg/executor/config"
	k8client "github.com/tlmiller/bookie/pkg/k8/client"
	k8config "github.com/tlmiller/bookie/pkg/k8/config"
	//"github.com/tlmiller/bookie/pkg/k8/controller/ingress"
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
	//	_ = api.NewServer()
	log.Info("creating main engine")
	eng := engine.NewEngine()

	log.Info("building configuration executors")
	executors, err := exconfig.ExecutorsForConfig()
	if err != nil {
		log.Fatalf("building configuration executors: %v", err)
	}

	log.Info("adding executors to main engine")
	err = eng.AddExecutors(executors...)
	if err != nil {
		log.Fatalf("adding executors to main engine: %v", err)
	}

	log.Info("fetching kubernetes client")
	c, err := k8client.NewForConfig(kubeMaster, kubeConfig)
	if err != nil {
		log.Fatalf("fetching kubernetes client: %v", err)
	}

	log.Info("building configuration kubernetes controllers")
	controllers, err := k8config.ControllersForConfig(c)
	if err != nil {
		log.Fatalf("building configuration kubernetes controllers: %v", err)
	}

	log.Info("adding controllers to main engine")
	eng.AddControllers(controllers...)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	log.Info("starting main engine")
	eng.Run()

	sig := <-sigCh
	log.Infof("recieved sig %s, terminating", sig)

	log.Info("waiting for main engine to stop")
	eng.Stop()

	//cl, err := client.ClientForConfig(kubeMaster, kubeConfig)
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}

	//ingress.NewController(cl)
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
