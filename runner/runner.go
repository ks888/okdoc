package runner

type Runner interface {
	Run(testCode string) *RunResult
	Patterns() []string
}

type RunResult struct {
	Success   bool
	HasRunner bool
	Message   string
}

var runners []Runner

func RegisterRunner(r Runner) {
	runners = append(runners, r)
}

func FindRunner(str string) Runner {
	for _, r := range runners {
		if RunnerMatch(r, str) {
			return r
		}
	}
	return nil
}

func RunnerMatch(r Runner, str string) bool {
	for _, x := range r.Patterns() {
		if str == x {
			return true
		}
	}
	return false
}

type baseRunner struct {
	patterns []string
}

func (r baseRunner) Run(testCode string) *RunResult {
	return &RunResult{true, true, "This is base runner. Runs no tests."}
}

func (r baseRunner) Patterns() []string {
	return r.patterns
}
