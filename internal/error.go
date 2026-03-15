package internal

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

func broadcastError(msg string, err error) error {
	return wrapError(ErrBroadcast, msg, err)
}

func logError(msg string, err error) error {
	return wrapError(ErrLog, msg, err)
}

func parseError(msg string, err error) error {
	return wrapError(ErrParse, msg, err)
}

func receiptError(msg string, err error) error {
	return wrapError(ErrReceipt, msg, err)
}

func transactionError(msg string, err error) error {
	return wrapError(ErrTransaction, msg, err)
}
