package cmd

import (
	"os"

	"github.com/ch3nnn/sql2pb/cmd/generation"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sql2pb",
	Short: "Generates a protobuf file from your database",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generation.GenCmd)

}
