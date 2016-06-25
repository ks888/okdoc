package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/ks888/okdoc/testrunner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	docExt = "md"
)

var OkdocCmd = &cobra.Command{
	Use:   "okdoc [flags] [target dir/file]",
	Short: "okdoc tests your documents",
	Long:  "okdoc tests your documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("target dir/file is not specified")
		}
		path := args[0]

		docs, err := collectDocs(path, docExt)
		if err != nil {
			return err
		}

		testSet := new(testrunner.TestSet)
		for _, docpath := range docs {
			testSet.AddTestFile(docpath)
		}

		err = testSet.RunAllTests()
		if err != nil {
			return err
		}
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

func init() {
	OkdocCmd.PersistentFlags().String("docker-image", "golang:1.6.2", "Docker image for isolated test env")
	viper.BindPFlag("docker-image", OkdocCmd.PersistentFlags().Lookup("docker-image"))
}
