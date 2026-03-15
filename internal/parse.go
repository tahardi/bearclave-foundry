package internal

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	HexBase         = 16
	HexStringPrefix = "0x"
	NumSize         = 64
)

var (
	ErrParse = errors.New("parse")
)

func BytesToHexString(b []byte) string {
	return HexStringPrefix + strings.ToLower(hex.EncodeToString(b))
}

func Uint64ToHexString(n uint64) string {
	return HexStringPrefix + strconv.FormatUint(n, HexBase)
}

func ParseBytesFromHexString(s string) ([]byte, error) {
	bytes, err := hex.DecodeString(
		strings.ToLower(
			strings.TrimPrefix(s, HexStringPrefix),
		),
	)
	if err != nil {
		msg := fmt.Sprintf("decoding hex string: %s", err)
		return nil, parseError(msg, err)
	}
	return bytes, nil
}

func ParseUint64FromHexString(s string) (uint64, error) {
	val, err := strconv.ParseUint(
		strings.ToLower(strings.TrimPrefix(s, HexStringPrefix)),
		HexBase,
		NumSize,
	)
	if err != nil {
		msg := fmt.Sprintf("parsing uint64 from hex string: %s", err)
		return 0, parseError(msg, err)
	}
	return val, nil
}
