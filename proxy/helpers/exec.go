package helpers

import (
	"../models"
	"../shared"
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteShellCommand(name string, params []string) (string, string, error) {
	cmd := exec.Command(name, params...)
	cmdOutput := &bytes.Buffer{}
	cmdError := &bytes.Buffer{}

	cmd.Stdout = cmdOutput
	cmd.Stderr = cmdError

	err := cmd.Run()

	return string(cmdOutput.Bytes()), string(cmdError.Bytes()), err
}

func ExecuteShellCommandForTask(name string, params []string, task models.Task, sharedContext *shared.SharedContext) (string, string, error) {
	cmdOut, cmdErr, err := ExecuteShellCommand(name, params)

	task.Log(fmt.Sprintf("executing command %s %s", name, params), "info", sharedContext)

	if err != nil {
		task.Log(fmt.Sprintf("failed with: %s", err), "info", sharedContext)
	}

	if cmdOut != "" {
		task.Log(fmt.Sprintf("out: %s", cmdOut), "info", sharedContext)
	}

	if cmdErr != "" {
		task.Log(fmt.Sprintf("err: %s", cmdErr), "info", sharedContext)
	}

	return cmdOut, cmdErr, err
}
