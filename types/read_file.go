package types

import (
	"bufio"
	"io"
)

func readFile(handle io.Reader) error {
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		// Do something with line
		_ = scanner.Text()
	}
	return nil
}
