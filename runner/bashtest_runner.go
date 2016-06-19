package runner

import (
	"io/ioutil"
	"os/exec"
)

type bashtestRunner struct {
	baseRunner
	command string
}

func init() {
	RegisterRunner(bashtestRunner{baseRunner: baseRunner{[]string{"bashtest"}}, command: "bash"})
}

func (runner bashtestRunner) Run(testCode string) *RunResult {
	testScriptFile, err := ioutil.TempFile("", "okdoc_")
	if err != nil {
		return &RunResult{false, true, "Failed to create temp file"}
	}

	_, err = testScriptFile.WriteString("set -e\n")

	_, err = testScriptFile.WriteString(testCode)
	if err != nil {
		return &RunResult{false, true, "Failed to write test script to temp file"}
	}

	// 'out' includes stdout AND stderr
	out, err := exec.Command(runner.command, testScriptFile.Name()).CombinedOutput()
	if err != nil {
		message := string(out)
		message += "\n"
		message += err.Error()
		return &RunResult{false, true, message}
	} else {
		return &RunResult{true, true, string(out)}
	}
}
