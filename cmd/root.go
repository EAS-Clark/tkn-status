package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tkn-status/helpers"
	"os"
)

var (
	namespace        string
	pipelineRunName  string
	taskNames        []string
	pollingFrequency string
)

var rootCmd = &cobra.Command{
	Use:   "your-app",
	Short: "Check the status of Tekton PipelineRun tasks",
	Run: func(cmd *cobra.Command, args []string) {
		err := helpers.RunTaskStatusChecker(namespace, pipelineRunName, taskNames, pollingFrequency)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the PipelineRun")
	rootCmd.Flags().StringVarP(&pipelineRunName, "pipelinerun", "p", "", "Name of the PipelineRun")
	rootCmd.Flags().StringSliceVarP(&taskNames, "tasks", "t", []string{}, "Names of the tasks to check (comma-separated)")
	rootCmd.Flags().StringVar(&pollingFrequency, "polling-frequency", "10s", "Frequency at which to poll task status")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
