package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Application Description",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
