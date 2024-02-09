package cache

import (
	"fmt"
)

type ExecutableNotFoundError struct {
	Ref string
}

func (e ExecutableNotFoundError) Error() string {
	return fmt.Sprintf("unable to find executable with reference %s", e.Ref)
}

func (e ExecutableNotFoundError) Unwrap() error {
	return fmt.Errorf("executable not found")
}

func NewExecutableNotFoundError(ref string) ExecutableNotFoundError {
	return ExecutableNotFoundError{Ref: ref}
}
