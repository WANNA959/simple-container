package container

import (
	"log"
	"os"
	"os/exec"
	"simple-container/pkg/cgroups"
	"simple-container/pkg/cgroups/subsystems"
	"strings"
	"syscall"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) error {
	parent, writePipe := NewParentProcess(tty)
	if parent == nil {
		log.Fatalln("New parent process error")
	}
	if err := parent.Start(); err != nil {
		log.Fatalln(err)
	}

	//netns.Set()

	cgroupManager := cgroups.NewCgroupManager("simple-container")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(0)
	return nil
}

// call init command here
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Printf("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Fatalln("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

// 生成匿名管道
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
