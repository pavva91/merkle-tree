package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func DirPathIsValid(dirPath string) bool {
	_, err := os.ReadDir(dirPath)
	return err == nil
}

func PathValidation(dirPath string) (string, error) {
	home, _ := os.UserHomeDir()
	fmt.Println("HOME:", home)

	firstChar := dirPath[:1]
	if firstChar != "." {
		containsHome := strings.Contains(dirPath, home)
		if !containsHome {
			err := errors.New("folder with absolute path must be inside home")
			return "", err
		}
	}
	if firstChar == "~" {
		dirPath = home + dirPath[1:]
	}
	if !DirPathIsValid(dirPath) {
		err := fmt.Errorf("folder %s does not exist", dirPath)
		return "", err
	}
	return dirPath, nil
}
