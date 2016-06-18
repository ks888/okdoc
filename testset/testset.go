package testset

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ks888/okdoc/parser"
	"github.com/ks888/okdoc/runner"
)

type Test struct {
	content *parser.CodeBlock
	name    string
	result  *runner.RunResult
}

type TestFile struct {
	path string
	list []*Test
}

type TestSet struct {
	testFiles []*TestFile
}

func (ts *TestSet) AddTestFile(path string) error {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	md_content := &parser.Markdown{Content: string(fileContent)}
	md_content.Parse()

	testFile := &TestFile{path: path}
	filename := filepath.Base(path)

	for i, codeBlock := range md_content.CodeBlocks {
		testName := fmt.Sprintf("%s#%d", filename, i)
		test := &Test{content: codeBlock, name: testName}
		testFile.list = append(testFile.list, test)
	}
	ts.testFiles = append(ts.testFiles, testFile)

	return nil
}

func (ts *TestSet) RunAllTests() error {
	for _, testFile := range ts.testFiles {
		for _, test := range testFile.list {
			runnerInst := runner.FindRunner(test.content.Lang)

			if runnerInst == nil {
				test.result = &runner.RunResult{true, false, "No test runner"}
			} else {
				test.result = runnerInst.Run(test.content.Block)
			}
		}
	}
	return nil
}

func (ts *TestSet) PrintTestStats() {
	numTests := 0
	numNoRunnerTests := 0
	failedTests := make([]*Test, 0)
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

	fmt.Printf("%d of %d tests are successful! (Plus, %d tests have no test runner)\n", numTests-len(failedTests), numTests, numNoRunnerTests)
	if len(failedTests) != 0 {
		fmt.Printf("\nThe list of failed tests: \n")

		for _, test := range failedTests {
			fmt.Printf("====== %s ======\n", test.name)
			fmt.Printf("%s\n", test.result.Message)
			fmt.Printf("===============\n")
		}
	}

}
