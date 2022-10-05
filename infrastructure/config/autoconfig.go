package config

import (
	"bufio"
	"dag-cli/pkg/color"
	"fmt"
	"io"
)

type Autoconfig struct {
	scanner *bufio.Scanner
}

func New(r io.Reader) *Autoconfig {
	return &Autoconfig{
		scanner: bufio.NewScanner(r),
	}
}

func (ac *Autoconfig) Ask(question string, _default string) string {
	if _default == "" {
		fmt.Printf("%s: ", color.Ize(color.Yellow, question))
	} else if _default == "no" {
		fmt.Printf("%s [%s]: ", color.Ize(color.Yellow, question), color.Ize(color.Red, _default))
	} else {
		fmt.Printf("%s [%s]: ", color.Ize(color.Yellow, question), color.Ize(color.Green, _default))
	}
	ac.scanner.Scan()

	text := ac.scanner.Text()
	if text == "" {
		return _default
	} else {
		return text
	}
}