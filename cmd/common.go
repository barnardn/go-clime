package cmd

import "github.com/barnardn/go-clime/whirly"

// command flags
var (
	isImperial    bool
	isQuiet       bool
	justTemp      bool
	justFeelsLike bool
	justZip       bool
)

func newWhirly() whirly.ProgressIndicator {
	whirlyType := whirly.Kitt
	if isQuiet {
		whirlyType = whirly.Empty
	}
	return whirly.New(whirlyType)

}
