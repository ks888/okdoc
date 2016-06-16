package testset

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/ks888/okdoc/parser"
)

type Test struct {
	content *parser.CodeBlock
	name    string
	result  string
	output  string
}

type TestFile struct {
	path string
	list []*Test
}

type TestSet struct {
	testFiles []*TestFile
}

func (ts *TestSet) AddTestFile(filepath string) error {
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	md_content := &parser.Markdown{Content: string(fileContent)}
	md_content.Parse()

	testFile := &TestFile{path: filepath}
	for i, codeBlock := range md_content.CodeBlocks {
		test := &Test{content: codeBlock, name: string(i)}
		testFile.list = append(testFile.list, test)
	}
	ts.testFiles = append(ts.testFiles, testFile)

	return nil
}

func (ts *TestSet) RunAllTests() error {
	for _, testFile := range ts.testFiles {
		fmt.Printf("TestFile: %s\n", testFile.path)

		for _, test := range testFile.list {
			testScriptFile, err := ioutil.TempFile("", "okdoc_")
			if err != nil {
				return errors.New("Failed to create temp file")
			}

			_, err = testScriptFile.WriteString(test.content.Block)
			if err != nil {
				return errors.New("Failed to write test script to temp file")
			}

			// 'out' includes stdout AND stderr
			out, err := exec.Command(test.content.Lang, testScriptFile.Name()).CombinedOutput()
			test.output = string(out)
			if err != nil {
				test.result = err.Error()
			}
		}
	}
	return nil
}

func (ts *TestSet) PrintTestStats() {
	numTests := 0
	failedTests := make([]*Test, 0)
	for _, testFile := range ts.testFiles {
		for _, test := range testFile.list {
			if test.result != "" {
				failedTests = append(failedTests, test)
			}
			numTests++
		}
	}

	fmt.Printf("%d of %d tests are successful!\n", numTests-len(failedTests), numTests)
	if len(failedTests) != 0 {
		fmt.Printf("The list of failed tests: \n")
		for _, test := range failedTests {
			fmt.Printf("* %s: %s\n", test.name, test.result)
		}
	}

}
