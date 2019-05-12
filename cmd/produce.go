package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	produceClientCmd = &cobra.Command{
		Use:     "produce",
		Short:   "produce client CLI to interact with the supermarket API",
		Aliases: []string{"produce", "pr"},
	}

	produceClientCmdEndpoint string
	produceClientCmdTimeout  string

	produceClientListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list produce in the inventory",
		Run:     produceClientList,
	}

	produceClientDeleteCmd = &cobra.Command{
		Use:     "delete [id]",
		Aliases: []string{"d"},
		Short:   "delete a produce item from the inventory",
		Run:     produceClientDelete,
	}

	produceClientListCmdParamSortBy string
	produceClientListCmdParamOrder  string
	produceClientListCmdParamLimit  string
	produceClientListCmdParamOffset string

	produceClientDeleteCmdParamProduceCode string
)

func init() {
	produceClientCmd.PersistentFlags().StringVarP(&produceClientCmdEndpoint, "endpoint", "e", "http://localhost:8000", "endpoint for the produce client to use")
	produceClientCmd.PersistentFlags().StringVarP(&produceClientCmdTimeout, "timeout", "t", "10s", "timeout for the produce client to set")

	produceClientCmd.AddCommand(produceClientListCmd)
	produceClientListCmd.Flags().StringVar(&produceClientListCmdParamSortBy, "sort_by", "", "optional value to choose how list response is sorted. Available case insensitive values are name, producecode, unitprice.")
	produceClientListCmd.Flags().StringVar(&produceClientListCmdParamOrder, "order", "", "optional value to choose how the response is order. Available case insensitive values are desc or descending. If no value or an unaccepted value is passed response will default to ascending.")
	produceClientListCmd.Flags().StringVar(&produceClientListCmdParamLimit, "limit", "", "optional value to choose the limit of the response after any optional offset. If an invalid value is passed this parameter will be ignored.")
	produceClientListCmd.Flags().StringVar(&produceClientListCmdParamOffset, "offset", "", "optional value to choose many items to start offset the response values. If an invalid value is passed this parameter will be ignored.")

	produceClientCmd.AddCommand(produceClientDeleteCmd)
	produceClientDeleteCmd.Flags().StringVar(&produceClientDeleteCmdParamProduceCode, "sort_by", "", "optional value to choose how list response is sorted. Available case insensitive values are name, producecode, unitprice.")
}

func produceClientList(cmd *cobra.Command, args []string) {
	client, err := newClient(
		withEndpoint(produceClientCmdEndpoint),
		withTimeout(produceClientCmdTimeout),
	)
	if err != nil {
		log.Fatalf("failed to create produce client: %s", err)
	}

	resp, statusCode, err := client.listProduce(
		produceClientListCmdParamSortBy,
		produceClientListCmdParamOrder,
		produceClientListCmdParamLimit,
		produceClientListCmdParamOffset,
	)
	if err != nil {
		log.Fatalf("failed to list produce: %s", err)
	}

	printResponse(resp, statusCode)
}

func produceClientDelete(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("must provide a produce code to delete")
	}

	client, err := newClient(
		withEndpoint(produceClientCmdEndpoint),
		withTimeout(produceClientCmdTimeout),
	)
	if err != nil {
		log.Fatalf("failed to create produce client: %s", err)
	}

	produceCode := args[0]

	resp, statusCode, err := client.deleteProduce(produceCode)
	if err != nil {
		log.Fatalf("failed to delete produce: %s", err)
	}

	printResponse(resp, statusCode)

}

func printResponse(obj interface{}, statusCode int) {
	js, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("failed to marshal output: %s", err)
		return
	}

	fmt.Fprintf(os.Stdout, fmt.Sprintf("Status Code: %d\n", statusCode))
	fmt.Fprintf(os.Stdout, string(js))
}
