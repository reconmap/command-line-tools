package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/reconmap/shared-lib/pkg/configuration"
	"github.com/reconmap/shared-lib/pkg/models"
)

func GetCommandById(id int) (*models.Command, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/commands/" + strconv.Itoa(id)

	client := &http.Client{}
	req, err := NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("unable to read response from server")
	}

	var command *models.Command = &models.Command{}

	if err = json.Unmarshal([]byte(body), command); err != nil {
		return command, err
	}

	return command, nil
}
func GetCommandUsageById(id int) (*models.CommandUsage, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/commands/usage/" + strconv.Itoa(id)

	client := &http.Client{}
	req, err := NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("unable to read response from server")
	}

	var command *models.CommandUsage = &models.CommandUsage{}

	if err = json.Unmarshal([]byte(body), command); err != nil {
		return command, err
	}

	return command, nil
}

func GetCommandsByKeywords(keywords string) (*models.Commands, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/commands?keywords=" + keywords

	client := &http.Client{}
	req, err := NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("unable to read response from server")
	}

	var commands *models.Commands = &models.Commands{}

	if err = json.Unmarshal(body, commands); err != nil {
		return commands, err
	}

	return commands, nil
}

func GetTasksByKeywords(keywords string) (*models.Tasks, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/tasks?keywords=" + keywords

	client := &http.Client{}
	req, err := NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("unable to read response from server")
	}

	var tasks *models.Tasks = &models.Tasks{}

	if err = json.Unmarshal(body, tasks); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func GetVulnerabilities() (*models.Vulnerabilities, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/vulnerabilities"

	client := &http.Client{}
	req, err := NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("unable to read response from server")
	}

	var vulnerabilities *models.Vulnerabilities = &models.Vulnerabilities{}

	if err = json.Unmarshal(body, vulnerabilities); err != nil {
		return vulnerabilities, err
	}

	return vulnerabilities, nil
}
