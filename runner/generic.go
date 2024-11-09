package runner

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/AidynClarke/golang-live-reload/buffer"
)

type genericRunner struct {
	command string
	args []string
	clearConsole bool

	process *exec.Cmd
	cmdBuffer *buffer.CmdBuffer
}

func newGenericRunner(command string, args []string, clearConsole bool, onClose func(), onError func(error)) *genericRunner {
	runner := &genericRunner{
		command: command, 
		args: args, 
		clearConsole: clearConsole,
	}

	runner.cmdBuffer = buffer.NewCmdBuffer(time.Second, runner.restart)

	return runner
}

func (r *genericRunner) Start() error {
	var err error

	if r.clearConsole {
		r.clear()
	}

	r.process = exec.Command(r.command, r.args...)

	r.process.Stdout = os.Stdout
	r.process.Stderr = os.Stderr

	err = r.process.Start()
	if err != nil {
		return err
	}

	return nil
}

func (r *genericRunner) Stop() (pid int, err error) {
	pid = r.process.Process.Pid

	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))

	return pid, kill.Run()
}

func (r *genericRunner) Restart() {
	r.cmdBuffer.Call()
}

func (r *genericRunner) restart() {
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

func (r *genericRunner) clear() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}