package cmd

import (
	"github.com/xmattstrongx/supermarket/api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func daemonRun(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.JSONFormatter{})
	server := api.NewServer()
	server.Serve()
}
