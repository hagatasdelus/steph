package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hagatasdelus/steph/db"
	"github.com/spf13/cobra"
)

var isYesterday bool

var setCmd = &cobra.Command{
	Use:   "set [steps]",
	Short: "Set steps for a specific day",
	Long:  `Set the number of steps for today or yesterday. Use -y or --yesterday to set yesterday's steps.`,
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

		if err := db.SetSteps(date, uint(steps)); err != nil {
			fmt.Println("Failed to set steps:", err)
			return
		}

		fmt.Printf("Steps set to %d for %s\n", steps, date.Format("2006/01/02"))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().BoolVarP(&isYesterday, "yesterday", "y", false, "Set steps for yesterday")
}
