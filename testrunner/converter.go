package testrunner

type NoConverterError struct{}

func (e NoConverterError) Error() string {
	return "No matched converter"
}

type converter struct {
	orgCmd       string
	convertedCmd string
	convert      func(orgCode string) (string, error)
}

var (
	converters = make([]converter, 0)
)

func Convert(cmd, testCode string) (string, string, error) {
	for _, cvt := range converters {
		if cmd == cvt.orgCmd {
			convertedTestCode, err := cvt.convert(testCode)
			if err != nil {
				return "", "", err
			}
			return cvt.convertedCmd, convertedTestCode, nil
		}
	}
	return "", "", NoConverterError{}
}

func bashConvert(testCode string) (string, error) {
	return "set -ex\n" + testCode, nil
}

func init() {
	bashConverter := converter{orgCmd: "bashtest", convertedCmd: "bash", convert: bashConvert}
	converters = append(converters, bashConverter)
}
