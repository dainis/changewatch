package main

import (
	"log"
	"os/exec"
	"strings"
	"sync"
)

type ExecLoop struct {
	executing     bool
	haveToExecute bool
	command       string
	args          []string
	commandString string
	mu            sync.Mutex
}

func (l *ExecLoop) Exec() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.haveToExecute = true

	if l.executing {
		log.Println("\x1b[0;37mAlready executing\x1b[0m")
		return
	}

	go l.execLoop()
}

func (l *ExecLoop) execLoop() {
	l.mu.Lock()
	l.executing = true
	l.mu.Unlock()

	for l.haveToExecute {
		l.mu.Lock()
		l.haveToExecute = false
		l.mu.Unlock()
		l.executeCommand()
	}

	l.mu.Lock()
	l.executing = false
	l.mu.Unlock()
}

func (l *ExecLoop) executeCommand() {
	log.Printf("\x1b[0;32mWill execute \x1b[1;32m%s\x1b[0m", l.commandString)

	cmd := exec.Command(l.command, l.args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("\x1b[1;31mFailed to execute command %s\x1b[0m", err)
		log.Printf("\x1b[0;31mOutput : \n\x1b[0m%s\n", string(output))
	} else {
		log.Printf("\x1b[1;32mWatch command executed\x1b[0m")
		log.Printf("\x1b[0;32mOutput : \n\x1b[0m%s\n", string(output))
	}
}

func NewExecLoop(command string, args []string) *ExecLoop {
	return &ExecLoop{
		command:       command,
		args:          args,
		commandString: command + " " + strings.Join(args, " "), executing: false,
	}
}
