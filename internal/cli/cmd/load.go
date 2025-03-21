package cmd

import (
	"fmt"

	"github.com/mangadev-courses/goexpert-rate-limiter/internal/cli/flags"
	"github.com/spf13/cobra"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

type Loader interface {
	LoadTest(url, apiKeyHeader string, requests int, concurrency int) error
}

func LoadCmd(loader Loader) *cobra.Command {
	var url string
	var apiKeyHeader string
	var requests int
	var concurrency int

	var loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Load Test",
		Long:  `This step executes the Load Test Commands.`,
		RunE:  runLoadCmd(loader),
	}

	flags.StringVarPRequired(loadCmd, &url, "url", "u", "", "URL to be tested")
	flags.IntVarPRequired(loadCmd, &requests, "requests", "r", 0, "Number of requests to be sent")
	flags.IntVarPRequired(loadCmd, &concurrency, "concurrency", "c", 0, "Number of concurrent requests to be sent")

	loadCmd.Flags().StringVarP(&apiKeyHeader, "apiKeyHeader", "a", "", "API Key Header")

	return loadCmd
}

func runLoadCmd(loader Loader) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		apiKeyHeader, _ := cmd.Flags().GetString("apiKeyHeader")
		requests, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		fmt.Printf("URL: %s\n", url)
		fmt.Printf("API Key Header: %s\n", apiKeyHeader)
		fmt.Printf("Requests: %d\n", requests)
		fmt.Printf("Concurrency: %d\n", concurrency)

		loader.LoadTest(url, apiKeyHeader, requests, concurrency)

		return nil
	}

}
