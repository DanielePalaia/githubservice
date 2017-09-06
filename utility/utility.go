package utility

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"newrelic/response"
	"newrelic/servicelog"
	"os"
	"strings"
	"time"
)

var Credentials string = ""
var Threading bool = false

func LoadAuth() {
	file, err := os.Open("conf")
	var Threads string
	if err != nil {
		Credentials = ""
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	Username, _ := scanner.ReadString(':')
	Username, _ = scanner.ReadString('\n')

	Username = Username[:len(Username)-1]
	Passwd, _ := scanner.ReadString(':')
	Passwd, _ = scanner.ReadString('\n')
	Passwd = Passwd[:len(Passwd)-1]

	if strings.Compare(Username, "") == 0 || strings.Compare(Passwd, "") == 0 {
		Credentials = ""
	} else {
		Credentials = fmt.Sprintf("%s:%s@", Username, Passwd)
	}
	Threads, _ = scanner.ReadString(':')
	Threads, _ = scanner.ReadString('\n')
	Threads = Threads[:len(Threads)-1]
	if strings.Compare(Threads, "yes") == 0 {
		Threading = true
	}
}

// Load pull request repository from file
func LoadContributorResponseFromFile(path string) (*response.ContributorsResponse, error) {
	outputRepos := response.CreateContributorsResponse()
	file := path[strings.LastIndex(path, "/")+1:]
	file = file + ".json"
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("error loading file")
		return nil, err
	}

	err = json.Unmarshal(raw, &outputRepos.Items)

	if err != nil {
		fmt.Printf("error unmarshalling file")
		return nil, err
	}

	return outputRepos, nil
}

func LoadContributorResponseFromNetwork(url string) (*response.ContributorsResponse, error) {
	var resp3 *http.Response
	var err error
	var record3 *response.ContributorsResponse
	// Removing http:://
	url = url[8:]
	url = fmt.Sprintf("https://%s%s", Credentials, url)
	// Load from network
	resp3, err = response.DoRequestAndReceiveResponse(url)
	if err != nil {
		return nil, err
	}

	record3 = response.CreateContributorsResponse()
	defer resp3.Body.Close()

	if err := json.NewDecoder(resp3.Body).Decode(&record3.Items); err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "error decoding")
	}

	return record3, nil
}

// Load pull request repository from file
func LoadReposResponseFromFile(file string) (*response.ReposResponse, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print("error loading")
		return nil, err
	}

	outputRepos := response.CreateReposResponse()
	err = json.Unmarshal(raw, outputRepos)
	if err != nil {
		fmt.Print("error unmarshalling")
		return nil, err
	}

	return outputRepos, nil
}
