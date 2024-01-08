package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Tracker is an app for tracking prices",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("inside of a main command")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
