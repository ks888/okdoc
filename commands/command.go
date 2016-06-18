package commands

import (
	"errors"
	"os"

	"github.com/ks888/okdoc/testset"
	"github.com/spf13/cobra"
)

var OkdocCmd = &cobra.Command{
	Use:   "okdoc: ",
	Short: "okdoc tests your documents",
	Long:  "okdoc tests your documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("The number of args must be 1")
		}

		filepath := args[0]
		if _, err := os.Stat(filepath); err != nil {
			return errors.New("File does not exist")
		}

		testSet := new(testset.TestSet)
		testSet.AddTestFile(filepath)

		testSet.RunAllTests()
		testSet.PrintTestStats()

		return nil
	},
}
