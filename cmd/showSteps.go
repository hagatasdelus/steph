package cmd

import (
	"fmt"
	"time"

	"https://github.com/hagatasdelus/steph/db"
	"github.com/spf13/cobra"
)

var showStepsCmd = &cobra.Command{
	Use:   "show-steps [YYYY/MM/DD]",
	Short: "Show steps for a specific date",
	Long:  `Display the number of steps recorded for a specific date in YYYY/MM/DD format. If no date is provided, shows steps for the past week.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showWeeklySteps()
			return
		}

		dateStr := args[0]
		parsedDate, err := time.Parse("2006/01/02", dateStr)
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY/MM/DD format.")
			return
		}
		
		date := time.Date(
			parsedDate.Year(),
			parsedDate.Month(),
			parsedDate.Day(),
			0, 0, 0, 0,
			time.Local,
		)

		steps, err := db.GetStepsForDate(date)
		if err != nil {
			fmt.Println("Failed to get steps:", err)
			return
		}

		if steps == 0 {
			fmt.Printf("No steps recorded for %s.\n", dateStr)
		} else {
			fmt.Printf("Steps for %s: %d\n", dateStr, steps)
		}
	},
}

func showWeeklySteps() {
	today := time.Now()
	startDate := today.AddDate(0, 0, -6)

	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	histories, err := db.GetStepHistoriesByDateRange(startDate, today)
	if err != nil {
			fmt.Println("Failed to retrieve step histories:", err)
			return
	}

	fmt.Println("Steps for the past week:")
	fmt.Println("------------------------\n")

	stepsMap := make(map[string]uint)
	for _, h := range histories {
		dateKey := h.Date.Format("2006/01/02")
		stepsMap[dateKey] = h.Steps
	}

	foundAny := false
	totalSteps := uint(0)
	for i := 0; i <= 6; i++ {
		date := startDate.AddDate(0, 0, i)
		dateKey := date.Format("2006/01/02")
		steps, exists := stepsMap[dateKey]
		totalSteps += steps

		if exists {
			fmt.Printf("%s: %d steps\n", dateKey, steps)
			foundAny = true
		} else {
			fmt.Printf("%s: No steps recorded\n", dateKey)
		}
	}
	fmt.Println("\n------------------------")
	fmt.Printf("Total steps: %d\nAverage: %d\n", totalSteps, totalSteps/7)
	
	if !foundAny {
		fmt.Println("\nNo step data recorded for the past week.")
	}
}

func init() {
	rootCmd.AddCommand(showStepsCmd)
}
