package whirly

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type ModeType uint8

const (
	Kitt ModeType = iota
	Empty
)

type ProgressIndicator struct {
	didStart      bool
	mode          ModeType
	configuration progressConfiguration
	stopChannel   chan struct{}
}

type progressConfiguration struct {
	sequence     []string
	isReversable bool
}

var configurations = map[ModeType]progressConfiguration{
	Kitt: {
		sequence:     []string{"█▒▒▒▒▒", "▒█▒▒▒▒", "▒▒█▒▒▒", "▒▒▒█▒▒", "▒▒▒▒█▒", "▒▒▒▒▒█"},
		isReversable: true,
	},
	Empty: {
		sequence:     []string{},
		isReversable: false,
	},
}

func (config *progressConfiguration) fullSequence() []string {
	if !config.isReversable {
		return config.sequence
	}
	forward := append([]string(nil), config.sequence...)
	reversed := append([]string(nil), config.sequence...)
	slices.Reverse(reversed[1:])
	return append(forward, reversed[2:]...)
}

func New(mode ModeType) ProgressIndicator {
	return ProgressIndicator{
		mode:          mode,
		configuration: configurations[mode],
		stopChannel:   make(chan struct{}),
	}
}

func (pi *ProgressIndicator) Start() {
	if pi.mode == Empty || pi.didStart {
		return
	}
	pi.didStart = true
	go func() {
		pi.hideCursor()
		indicatorSequnce := pi.configuration.fullSequence()
		for {
			for _, seq := range indicatorSequnce {
				select {
				case <-pi.stopChannel:
					pi.showCursor()
					return
				default:
					fmt.Print(seq)
					time.Sleep(250 * time.Millisecond)
					pi.eraseLine()
				}
			}
		}
	}()
}

func (pi *ProgressIndicator) Stop() {
	if pi.mode == Empty || !pi.didStart {
		return
	}
	pi.stopChannel <- struct{}{}
	pi.didStart = false
	pi.eraseLine()
}

func (pi *ProgressIndicator) eraseLine() {
	eraseCodeString := strings.Builder{}
	eraseCodeString.WriteString("\r\033[K")
	fmt.Print(eraseCodeString.String())
}

func (pi *ProgressIndicator) hideCursor() {
	fmt.Print("\033[?25l")
}

func (pi *ProgressIndicator) showCursor() {
	fmt.Print("\033[?25h")
}
