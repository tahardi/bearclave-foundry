package foundry

import "fmt"

func wrapError(baseErr error, msg string, err error) error {
	switch {
	case msg == "" && err == nil:
		return baseErr
	case msg != "" && err != nil:
		return fmt.Errorf("%w: %s: %w", baseErr, msg, err)
	case msg != "":
		return fmt.Errorf("%w: %s", baseErr, msg)
	default:
		return fmt.Errorf("%w: %w", baseErr, err)
	}
}

func anvilError(msg string, err error) error {
	return wrapError(ErrAnvil, msg, err)
}

func forgeError(msg string, err error) error {
	return wrapError(ErrForge, msg, err)
}

func foundryError(msg string, err error) error {
	return wrapError(ErrFoundry, msg, err)
}
