package file

import (
	"errors"
	"os"
)

func HomeDir(home string) (string, error) {
	if home != "" {
		return home, nil
	}

	var err error = nil

	if home, err = os.UserHomeDir(); err == nil {
		return home, nil
	}

	home, err = os.Getwd()
	if err != nil {
		return "", errors.New("home Directory not detected")
	}
	return home, nil
}

func ExitDir(abs string) bool {
	_, err := os.Stat(abs)
	return !os.IsNotExist(err)
}
