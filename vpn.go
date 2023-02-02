package vpn

import "fmt"

const (
	ErrNotFound = fmt.Errorf(`error: does not appear to be through a known vpn`)
)

func Check(addr string) (string, error) {
	return "", ErrNotFound
}
