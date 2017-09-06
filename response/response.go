/* This file contains mainly the structures needed to xml unmarshal and ordered support */
package response

import (
	"net/http"
	"newrelic/servicelog"
	"time"
)

// Store the result
type Result struct {
	Login               string
	Total_Contribution  int
	Number_Repositories int
}

type ResultsOrderedByContribution []Result
type ResultsOrderedByRepositories []Result

// Implements the sort interface for result
func (r ResultsOrderedByContribution) Len() int {
	return len(r)
}
func (r ResultsOrderedByContribution) Less(i, j int) bool {
	return r[i].Total_Contribution > r[j].Total_Contribution
}
func (r ResultsOrderedByContribution) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Implements the sort interface for result
func (r ResultsOrderedByRepositories) Len() int {
	return len(r)
}
func (r ResultsOrderedByRepositories) Less(i, j int) bool {
	return r[i].Number_Repositories > r[j].Number_Repositories
}
func (r ResultsOrderedByRepositories) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// First call to know numbers of users in a city.
type UsersResponse struct {
	Total_count int     `json:"total_count"`
	Items       []Users `json:"items"`
}
type Users struct {
	Login string `json:"login"`
}

// Second call, for every users check the last completed PR and take the repo url
type ReposResponse struct {
	Total_count int     `json:"total_count"`
	Items       []Repos `json:"items"`
}
type Repos struct {
	Repository_url string `json:"repository_url"`
}

// Third call check contributor names in the repo.
type ContributorsResponse struct {
	Items []Users `json:"items"`
}

func CreateUsersResponse() *UsersResponse {
	gitResponse := new(UsersResponse)
	gitResponse.Total_count = 0
	gitResponse.Items = make([]Users, 0)

	return gitResponse
}

func CreateReposResponse() *ReposResponse {
	gitResponse := new(ReposResponse)
	gitResponse.Total_count = 0
	gitResponse.Items = make([]Repos, 0)

	return gitResponse
}

func CreateContributorsResponse() *ContributorsResponse {
	gitResponse := new(ContributorsResponse)
	gitResponse.Items = make([]Users, 0)

	return gitResponse
}

// Do requests go github using a url
func DoRequestAndReceiveResponse(url string) (*http.Response, error) {
	logger := servicelog.GetInstance()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Println(time.Now().UTC(), "Http new request error")
		return nil, err

	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(time.Now().UTC(), "bad error in client Do")
		return nil, err
	}
	return resp, nil
}
