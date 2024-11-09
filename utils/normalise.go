package utils

import (
	"os"
	"strings"
)

func Normalise(path string) string {
	return strings.ReplaceAll(path, string(os.PathSeparator), "/")
}