package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "video-splitter",
	Short: "A tool to split videos based on a plan file",
	Long:  `video-splitter is a CLI tool that splits a video into multiple segments based on a provided plan file.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
