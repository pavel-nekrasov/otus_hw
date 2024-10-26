package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	process := exec.Command(cmd[0], cmd[1:]...) //nolint:all
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	setEnv(process, env)
	err := process.Run()
	if err != nil {
		return -1
	}
	return process.ProcessState.ExitCode()
}

func setEnv(process *exec.Cmd, env Environment) {
	processEnv := make([]string, 0)
	enrichWithExternal(env)
	for key, value := range env {
		if value.NeedRemove {
			os.Unsetenv(key)
		} else {
			processEnv = append(processEnv, fmt.Sprintf("%v=%v", key, value.Value))
		}
	}
	process.Env = processEnv
}

func enrichWithExternal(env Environment) {
	for _, e := range os.Environ() {
		items := strings.Split(e, "=")
		_, ok := env[items[0]]
		if !ok {
			env[items[0]] = EnvValue{Value: items[1]}
		}
	}
}
