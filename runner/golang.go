package runner

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"test/buffer"
	"time"
)

type golangRunner struct {
	command string
	args []string
	clearConsole bool

	process *exec.Cmd
	cmdBuffer *buffer.CmdBuffer
}

func newGolangRunner(command string, args []string, clearConsole bool, onClose func(), onError func(error)) *golangRunner {
	runner := &golangRunner{
		command: command, 
		args: args, 
		clearConsole: clearConsole,
	}

	runner.cmdBuffer = buffer.NewCmdBuffer(time.Second, runner.restart)

	return runner
}

func (r *golangRunner) Start() error {
	var err error

	if r.clearConsole {
		r.clear()
	}

	exe := "tmp/" + strings.Split(r.command, ".")[0] + ".exe"
	os.MkdirAll("tmp", os.ModePerm)

	// Build
	r.process = exec.Command("go", []string{"build", "-o", exe, r.command}...)
	r.process.Stdout = os.Stdout
	r.process.Stderr = os.Stderr
	err = r.process.Run()
	if err != nil {
		return err
	}

	// Execute
	r.process = exec.Command(exe, r.args...)
	r.process.Stdout = os.Stdout
	r.process.Stderr = os.Stderr
	err = r.process.Start()
	if err != nil {
		return err
	}

	return nil
}

func (r *golangRunner) Stop() (int, error) {
	pid := r.process.Process.Pid

	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))

	killErr := kill.Run()

	rmErr := os.RemoveAll("tmp")

	if rmErr != nil {
		log.Println("Error removing tmp directory:", rmErr)
	}

	return pid, killErr
}

func (r *golangRunner) Restart() {
	r.cmdBuffer.Call()
}

func (r *golangRunner) restart() {
	_, err := r.Stop()

	if err != nil {
		log.Println("Error stopping process:", err)
		return
	}

	err = r.Start()
	if err != nil {
		log.Println("Error starting process:", err)
	}
}

func (r *golangRunner) clear() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}