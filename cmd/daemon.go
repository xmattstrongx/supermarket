package cmd

import (
	"github.com/xmattstrongx/supermarket/api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "daemon runs the supermarket daemon as a docker container",
	Run:   daemonRun,
}

func daemonRun(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.JSONFormatter{})
	server := api.NewServer()
	server.Serve()
}
