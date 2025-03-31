package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"https://github.com/hagatasdelus/steph/db"
)

var rootCmd = &cobra.Command{
	Use:   "steph",
	Short: "Walking & Running Contribution History",
	Long: `StepH is a TUI library that can record the number of 
steps taken and display them in a mapping history format.
Steps history for those who don't want to leave CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := db.InitDB(); err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
