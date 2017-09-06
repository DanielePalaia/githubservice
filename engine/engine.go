package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"newrelic/response"
	"newrelic/servicelog"
	"newrelic/utility"
	"sort"
	"strings"
	"sync"
	"time"
)

func CallGitApi(w http.ResponseWriter, location string, limitcount int) error {
	// It will contain the users ordered by contributions
	var outputByContribution response.ResultsOrderedByContribution = make([]response.Result, 0)
	// GitHub returns maximum 150 results per page so we have to loop through pages
	i := 1
	for {
		fmt.Printf("iteration: %d\n", i)
		// Find all users in a given location (100 results)
		record, err := findLocationUsers(location, i)
		if err != nil {
			return err
		}

		// If we reached the end of the pages we stop
		if record.Total_count == 0 {
			break
		}

		// For every user find the last 100 PR closed that are in other repositories and returns the result
		outputByContribution, err = findRepositoriesContributed(outputByContribution, record)
		if err != nil {
			return err
		}
		i++
	}

	//Order them for number of repositories contributed
	outputByRepositories := OrderByRepositories(outputByContribution, limitcount)
	if err := encodeToJson(w, outputByRepositories); err != nil {
		return err
	}
	return nil
}

// Find the users in a given location
func findLocationUsers(location string, i int) (*response.UsersResponse, error) {
	urlLocation := fmt.Sprintf("https://%sapi.github.com/search/users?q=+location:%s&sort=repositories&page=%d&per_page=150", utility.Credentials, location, i)
	resp, err := response.DoRequestAndReceiveResponse(urlLocation)
	if err != nil {
		return nil, err
	}
	usersRecord := response.CreateUsersResponse()

	if err := json.NewDecoder(resp.Body).Decode(&usersRecord); err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "Json Decoding error")
		return nil, err
	}

	resp.Body.Close()
	return usersRecord, nil
}

func routineFindUserRepositoriesContributed(user string, mutex *sync.Mutex, output *response.ResultsOrderedByContribution, wg *sync.WaitGroup) {
	var r response.Result
	if wg != nil {
		defer wg.Done()
	}
	fmt.Printf("processing user: %s\n", user)
	// Search for that last 100 closed pull requests the user submitted and take the repository
	url := fmt.Sprintf("https://%sapi.github.com/search/issues?q=type:pr+state:closed+closed:>2017-01-01+author:%s&per_page=100&page=1", utility.Credentials, user)
	resp, err := response.DoRequestAndReceiveResponse(url)
	if err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "error in DoRequestAndReceiveResponse")
		return
	}
	pullRequestsrecord := response.CreateReposResponse()
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&pullRequestsrecord); err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "Json Decoding error")
		return
	}

	if r, err = FindRepositoryContributors(pullRequestsrecord, user, false); err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "error in FindRepositoryContributors")
		return
	}

	// Output is accessed concurrently by all the goroutines an it needs to be protected.
	if mutex != nil {
		mutex.Lock()
	}
	*output = append(*output, r)
	if mutex != nil {
		mutex.Unlock()
	}

}

// For every user search for PRS closed and accepted from the beginning of the year
func findRepositoriesContributed(output response.ResultsOrderedByContribution, users *response.UsersResponse) (response.ResultsOrderedByContribution, error) {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	if utility.Threading == true {
		wg.Add(len(users.Items))
	}
	for _, user := range users.Items {
		// We are running in multithreading
		if utility.Threading == true {
			// Run a goroutine for every user to fine its contributions
			go routineFindUserRepositoriesContributed(user.Login, &mutex, &output, &wg)
		} else {
			routineFindUserRepositoriesContributed(user.Login, nil, &output, nil)
		}

	}

	if utility.Threading == true {
		wg.Wait()
	}

	return output, nil
}

// For every repository then look that the user is a contributor of that repository
func FindRepositoryContributors(pullRequests *response.ReposResponse, user string, test bool) (response.Result, error) {
	var contributions int
	var repos int
	duplicateRepos := make(map[string]int)

	for _, repo := range pullRequests.Items {
		var dup int
		var ok bool
		if dup, ok = duplicateRepos[repo.Repository_url]; ok != true {
			var record3 *response.ContributorsResponse
			var err error
			if !test {
				url := fmt.Sprintf("%s/contributors", repo.Repository_url)
				record3, err = utility.LoadContributorResponseFromNetwork(url)
				if err != nil {
					return response.Result{}, err
				}
			} else {
				file := fmt.Sprintf("%s", repo.Repository_url)
				record3, err = utility.LoadContributorResponseFromFile(file)
			}

			for _, contributor := range record3.Items {
				if strings.Compare(contributor.Login, user) == 0 {
					contributions++
					repos++
					duplicateRepos[repo.Repository_url]++
				}
			}
			if contributions == 0 {
				duplicateRepos[repo.Repository_url] = 0
			}

		}

		if dup > 0 {
			contributions++
		}

	}
	r := response.Result{Login: user, Number_Repositories: repos, Total_Contribution: contributions}
	return r, nil
}

func encodeToJson(w http.ResponseWriter, response response.ResultsOrderedByRepositories) error {
	// Encode in a json format
	js, err := json.Marshal(response)

	if err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "Decoding error")
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func OrderByRepositories(outputByContribution response.ResultsOrderedByContribution, limitcount int) response.ResultsOrderedByRepositories {
	var outputByRepositories response.ResultsOrderedByRepositories = make([]response.Result, 0)
	// Order by contributions
	sort.Sort(outputByContribution)
	//Take the first limitcount contributor
	for i, user := range outputByContribution {
		outputByRepositories = append(outputByRepositories, user)
		if i >= limitcount-1 {
			break
		}
	}

	sort.Sort(outputByRepositories)
	return outputByRepositories
}
