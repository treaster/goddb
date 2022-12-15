package lambdatk

import (
	"errors"
	"fmt"
)

func newErrf(format string, a ...any) error {
	return errors.New(fmt.Sprintf(format, a...))
}
