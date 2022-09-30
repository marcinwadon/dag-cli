package layer

import (
	"errors"
	"fmt"
)

type Layer string

const (
	L0 Layer = "l0"
	L1 Layer = "l1"
)

var Layers = [...]Layer{L0, L1}

func ParseString(in string) (*Layer, error) {
	for _, l := range Layers {
		if in == fmt.Sprint(l) {
			return &l, nil
		}
	}

	return nil, errors.New("unknown layer")
}