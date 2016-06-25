package testrunner

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type TestResult struct {
	Success   bool
	HasRunner bool
	Message   string
}

func Run(cmd, testCode string) *TestResult {
	testFilePath, err := saveTestCode(testCode)
	if err != nil {
		return &TestResult{false, true, err.Error()}
	}
	defer os.RemoveAll(path.Dir(testFilePath))

	if running, msg := dockerRunning(); !running {
		return &TestResult{false, true, msg}
	}

	success, msg := runCmdByDocker(cmd, testFilePath)
	return &TestResult{success, true, msg}
}

func dockerRunning() (bool, string) {
	out, err := exec.Command("docker", "info").CombinedOutput()
	if err != nil {
		return false, string(out)
	}
	return true, ""
}

func runCmdByDocker(cmd, testFilePath string) (bool, string) {
	dockerTestDir := "/tmp/okdoc/"
	dockerTestFile := dockerTestDir + path.Base(testFilePath)
	volumeMapping := path.Dir(testFilePath) + ":" + dockerTestDir

	args := []string{"run", "-v", volumeMapping, viper.GetString("docker-image"), cmd, dockerTestFile}

	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		message := string(out)
		message += "\n"
		message += err.Error()
		return false, message
	}

	return true, string(out)
}

func saveTestCode(script string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	useDockerMachine := false
	switch runtime.GOOS {
	case "darwin":
		useDockerMachine = true
	default:
		return "", errors.New("Only Mac OS is supported")
	}

	if useDockerMachine && !strings.HasPrefix(wd, "/Users") {
		return "", errors.New("Current directory has to be under /Users due to the restriction of docker-machine")
	}

	testScriptDir, err := ioutil.TempDir(wd, "temp_")
	if err != nil {
		return "", err
	}

	testScriptFile, err := ioutil.TempFile(testScriptDir, "script_")
	if err != nil {
		return "", err
	}

	_, err = testScriptFile.WriteString(script)
	if err != nil {
		return "", err
	}

	return testScriptFile.Name(), nil
}
