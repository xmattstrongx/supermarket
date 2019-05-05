package cmd

import (
	"fmt"

	"github.com/xmattstrongx/supermarket/api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func daemonRun(cmd *cobra.Command, args []string) {
	fmt.Println("GOo")

	log.SetFormatter(&log.JSONFormatter{})
	api.Serve()
}
