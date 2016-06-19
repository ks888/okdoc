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
	testDir, err := ioutil.TempDir("", "okdoc_")
	if err != nil {
		return &RunResult{false, true, "Failed to create test directory"}
	}

	testScriptFile, err := ioutil.TempFile(testDir, "script_")
	if err != nil {
		return &RunResult{false, true, "Failed to create temp file"}
	}

	_, err = testScriptFile.WriteString("set -ex\n" + testCode)
	if err != nil {
		return &RunResult{false, true, "Failed to write test script to temp file"}
	}

	// 'out' includes stdout AND stderr
	cmd := exec.Command(runner.command, testScriptFile.Name())
	cmd.Dir = testDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		message := string(out)
		message += "\n"
		message += err.Error()
		return &RunResult{false, true, message}
	} else {
		return &RunResult{true, true, string(out)}
	}
}
