package handler

import (
	"errors"
	"fmt"
)

var (
	ErrNegativeZeroBuilder = func(f string) error {
		return fmt.Errorf("%s is negative or zero", f)
	}
	ErrNotFoundBuilder = func(f string) error {
		return fmt.Errorf("%s not found", f)
	}

	ErrHeightOutOfRange     = errors.New("height must be 1 to 30")
	ErrCoordinateOutOfBound = errors.New("coordinate out of bound")
	ErrTreeExist            = errors.New("plot already has tree")
)
