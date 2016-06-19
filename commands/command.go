package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/ks888/okdoc/testset"
	"github.com/spf13/cobra"
)

const (
	docExt = "md"
)

var OkdocCmd = &cobra.Command{
	Use:   "okdoc: ",
	Short: "okdoc tests your documents",
	Long:  "okdoc tests your documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("The number of args must be 1")
		}
		path := args[0]

		docs, err := collectDocs(path, docExt)
		if err != nil {
			return err
		}

		testSet := new(testset.TestSet)
		for _, docpath := range docs {
			testSet.AddTestFile(docpath)
		}

		testSet.RunAllTests()
		testSet.PrintTestStats()

		return nil
	},
}

func collectDocs(path, ext string) ([]string, error) {
	docs := make([]string, 128)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ext) {
			docs = append(docs, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return docs, nil
}
