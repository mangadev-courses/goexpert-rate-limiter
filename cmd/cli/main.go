package main

import (
	"fmt"
	"os"

	"github.com/mangadev-courses/goexpert-rate-limiter/internal/cli/cmd"
	"github.com/mangadev-courses/goexpert-rate-limiter/pkg/goten"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stressTest",
	Short: "Stress Test",
	Long:  "This is a CLI tool to stress test the application",
}

func main() {
	goten := goten.New()

	rootCmd.AddCommand(cmd.LoadCmd(goten))

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
