package cmd

import "fmt"

const (
	scYes           = "yes"
	scNo            = "no"
	scVerbose       = "verbose"
	DefaultScOption = scNo
)

func scorecardOptions() []string {
	return []string{scNo, scYes, scVerbose}
}

func ValidateScorecardOption(opt string) error {
	for _, o := range scorecardOptions() {
		if o == opt {
			return nil
		}
	}
	return fmt.Errorf("invalid scorecard option: %s", opt)
}

func IsScorecardEnabled(when string) bool {
	return when == scYes || when == scVerbose
}

func IsScorecardVerbose(when string) bool {
	return when == scVerbose
}
