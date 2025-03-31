package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	Version = "0.1.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of StepHist`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("StepH v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	
	rootCmd.Flags().BoolP("version", "V", false, "Print version information and exit")
	
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("StepH v%s\n", Version)
			os.Exit(0)
		}
	}
}
