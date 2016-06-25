package testrunner

import (
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
)

type Test struct {
	content *CodeBlock
	name    string
	result  *TestResult
}

type TestFile struct {
	path string
	list []*Test
}

type TestSet struct {
	testFiles []*TestFile
}

var (
	blue = color.New(color.FgBlue).SprintFunc()
	red  = color.New(color.FgRed).SprintFunc()
)

func (ts *TestSet) AddTestFile(path string) error {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	mdContent := &Markdown{Content: string(fileContent)}
	mdContent.Parse()

	testFile := &TestFile{path: path}

	for _, codeBlock := range mdContent.CodeBlocks {
		testName := fmt.Sprintf("%s#L%d", path, codeBlock.StartLine)
		test := &Test{content: codeBlock, name: testName}
		testFile.list = append(testFile.list, test)
	}
	ts.testFiles = append(ts.testFiles, testFile)

	return nil
}

func (ts *TestSet) RunAllTests() error {
	for _, testFile := range ts.testFiles {
		for _, test := range testFile.list {
			cmd := test.content.Command
			testCode := test.content.Block

			var err error
			cmd, testCode, err = Convert(cmd, testCode)
			if err != nil {
				_, ok := err.(NoConverterError)
				if ok {
					test.result = &TestResult{true, false, err.Error()}
					continue
				} else {
					return err
				}
			}

			test.result = Run(cmd, testCode)

			if test.result.Success {
				if test.result.HasRunner {
					fmt.Printf(blue("%s...ok\n"), test.name)
				}
			} else {
				fmt.Printf(red("%s...failed\n"), test.name)
				fmt.Printf("%s\n", test.result.Message)
			}
		}
	}
	return nil
}

func (ts *TestSet) PrintTestStats() {
	numTests := 0
	numNoRunnerTests := 0
	var failedTests []*Test
	for _, testFile := range ts.testFiles {
		for _, test := range testFile.list {
			if !test.result.HasRunner {
				numNoRunnerTests++
				continue
			}

			if !test.result.Success {
				failedTests = append(failedTests, test)
			}
			numTests++
		}
	}

	col := blue
	if len(failedTests) != 0 {
		col = red
	}
	fmt.Printf(col("%d of %d tests are successful! (Plus, %d tests have no test runner)\n"), numTests-len(failedTests), numTests, numNoRunnerTests)
}
