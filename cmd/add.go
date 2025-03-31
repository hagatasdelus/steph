package cmd

import (
	"fmt"
	"strconv"
	"time"

	"https://github.com/hagatasdelus/steph/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [steps]",
	Short: "Add steps to a specific day",
	Long:  `Add steps to the current count for today or yesterday. Use -y or --yesterday to add to yesterday's steps.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		steps, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			fmt.Println("Invalid steps value:", err)
			return
		}

		date := time.Now()
		if isYesterday {
			date = date.AddDate(0, 0, -1)
		}

		if err := db.AddSteps(date, uint(steps)); err != nil {
			fmt.Println("Failed to add steps:", err)
			return
		}

		currentSteps, err := db.GetStepsForDate(date)
		if err != nil {
			fmt.Println("Failed to retrieve updated steps:", err)
			return
		}

		fmt.Printf("Added %d steps. Total for %s: %d steps\n", steps, date.Format("2006/01/02"), currentSteps)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&isYesterday, "yesterday", "y", false, "Add steps for yesterday")
}
