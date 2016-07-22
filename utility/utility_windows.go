package utility

import "os"

func getExecRoot() (path string, err error) {
	return os.Getwd()
}
