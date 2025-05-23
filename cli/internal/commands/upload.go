package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/reconmap/cli/internal/terminal"
	"github.com/reconmap/shared-lib/pkg/api"
	"github.com/reconmap/shared-lib/pkg/configuration"
	"github.com/reconmap/shared-lib/pkg/models"
)

func UploadResults(projectId int, usage *models.CommandUsage) error {
	return UploadCommandOutputUsingFileName(projectId, usage)
}

func UploadCommandOutputUsingFileName(projectId int, usage *models.CommandUsage) error {
	if len(strings.TrimSpace(usage.OutputFilename)) == 0 {
		return errors.New("The command has not defined an output filename. Nothing has been uploaded to the server.")
	}

	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}
	var remoteURL string = config.ApiUrl + "/commands/outputs"

	var client *http.Client = &http.Client{}
	err = Upload(client, remoteURL, usage.OutputFilename, usage.ID, projectId)
	return err
}

func Upload(client *http.Client, url string, outputFileName string, usageId int, projectId int) (err error) {

	if _, err := os.Stat(outputFileName); os.IsNotExist(err) {
		return fmt.Errorf("Output file '%s' could not be found", outputFileName)
	}

	file, err := os.Open(filepath.Clean(outputFileName))
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("resultFile", filepath.Base(outputFileName))
	_, err = io.Copy(part, file)

	if err = writer.WriteField("commandUsageId", strconv.Itoa(usageId)); err != nil {
		return
	}
	if projectId != 0 {
		if err = writer.WriteField("projectId", strconv.Itoa(projectId)); err != nil {
			return
		}
	}

	if err = writer.Close(); err != nil {
		return
	}

	req, err := api.NewRmapRequest("POST", url, body)
	if err != nil {
		return
	}

	err = api.AddBearerToken(req)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	terminal.PrintYellowDot()
	fmt.Printf(" Uploading command output '%s' to the server.\n", outputFileName)
	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode == http.StatusUnauthorized {
		err = fmt.Errorf("your session has expired. Please login again")
	}
	terminal.PrintGreenTick()
	fmt.Printf(" Done\n")

	return
}
