package runner

import "strings"


type Runner interface {
	Start() error
	Stop() (pid int, err error)
	Restart()
}

func NewRunner(command string, args []string, clearConsole bool, onClose func(), onError func(error)) Runner {
	if strings.HasSuffix(command, ".go") {
		return newGolangRunner(command, args, clearConsole, onClose, onError)
	}

	return newGenericRunner(command, args, clearConsole, onClose, onError)
}
