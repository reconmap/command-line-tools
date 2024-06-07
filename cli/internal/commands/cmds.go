package commands

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/reconmap/shared-lib/pkg/logging"
	"github.com/reconmap/shared-lib/pkg/models"

	"github.com/reconmap/cli/internal/terminal"
	"github.com/reconmap/shared-lib/pkg/io"
)

func RunCommand(command *models.CommandUsage, vars []string) error {
	logger := logging.GetLoggerInstance()

	var err error
	argsRendered := terminal.ReplaceArgs(command, vars)
	log.Println("Command to run: " + command.ExecutablePath + " " + argsRendered)

	cmd := exec.Command(command.ExecutablePath, strings.Fields(argsRendered)...) // #nosec G204
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = io.CopyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = io.CopyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout), string(stderr)

	outputFilename := filepath.Clean(strconv.Itoa(command.ID) + ".out")
	f, err := os.Create(outputFilename)

	defer func() {
		if err := f.Close(); err != nil {
			logger.Warn("Error closing file: %s", err)
		}
	}()
	_, err = f.WriteString(outStr)
	if err != nil {
		logger.Error(err)
	}
	command.OutputFilename = outputFilename

	if len(errStr) > 0 {
		log.Println(errStr)
	}

	return err
}
