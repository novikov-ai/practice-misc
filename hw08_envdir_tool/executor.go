package main

import (
	"log"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	err := configEnvironment(env)
	if err != nil {
		log.Fatal(err)
	}

	command := exec.Cmd{
		Path:   cmd[0],
		Env:    os.Environ(),
		Args:   cmd,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = command.Run()
	if err != nil {
		log.Fatal(err)
	}

	return 0
}

func configEnvironment(env Environment) error {
	for envName, envValue := range env {
		err := os.Setenv(envName, envValue.Value)
		if err != nil {
			return err
		}

		if !envValue.NeedRemove {
			continue
		}

		err = os.Unsetenv(envName)
		if err != nil {
			return err
		}
	}
	return nil
}
