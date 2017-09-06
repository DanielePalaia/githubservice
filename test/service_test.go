package testing

import (
	"newrelic/engine"
	"newrelic/response"
	"newrelic/utility"
	"testing"
)

// Receive a login user name, a set of pull request and returns the contributions for the user
type structTestingFindRepositoryContributors struct {
	// The first two used just by the test
	inputUser string
	inputFile string
	result    response.Result
	err       error
}

// Receive a login user name, a set of pull request and returns the contributions for the user
type structTestingOrderByRepositories struct {
	// The first two used just by the test
	inputListContribution           response.ResultsOrderedByContribution
	outputListOrderedByRepositories response.ResultsOrderedByRepositories
	limitCount                      int
}

func TestExerciseindRepositoryContributors(t *testing.T) {
	tests := map[string]structTestingFindRepositoryContributors{
		// This user has done several prs but just 5 contributions in repos (4 in total) it doesn't own
		"first_test": {
			inputUser: "amegawac",
			inputFile: "pull_request_data.json",
			result:    response.Result{Login: "amegawac", Total_Contribution: 5, Number_Repositories: 4},
			err:       nil,
		},
		// This is not the user we are looking for
		"second_test": {
			inputUser: "notexistent",
			inputFile: "pull_request_data.json",
			result:    response.Result{Login: "notexistent", Total_Contribution: 0, Number_Repositories: 0},
			err:       nil,
		},
		// This is not the user we are looking for
		"third_test": {
			inputUser: "afc163",
			inputFile: "pull_request_data.json",
			result:    response.Result{Login: "notexistent", Total_Contribution: 10, Number_Repositories: 9},
			err:       nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pullRequests, err := utility.LoadReposResponseFromFile(tc.inputFile)
			if err != nil {
				t.Fatal("test failed system error during execution of the test")
			}
			res, err := service.FindRepositoryContributors(pullRequests, tc.inputUser, true)
			if err != nil {
				t.Fatal("error not expected")
			}

			if res.Total_Contribution != tc.result.Total_Contribution {
				t.Fatalf("test failed expected Total_Contribution %d got %d", tc.result.Total_Contribution, res.Total_Contribution)
			}
			if res.Number_Repositories != tc.result.Number_Repositories {
				t.Fatalf("test failed expected Number Repositories %d got %d", tc.result.Number_Repositories, res.Number_Repositories)
			}

		})
	}
}

func loadGenericInputMap() response.ResultsOrderedByContribution {
	var inputListContribution response.ResultsOrderedByContribution = make([]response.Result, 0)
	r := response.Result{Login: "amegawac", Total_Contribution: 20, Number_Repositories: 5}
	inputListContribution = append(inputListContribution, r)
	r2 := response.Result{Login: "pluto", Total_Contribution: 9, Number_Repositories: 1}
	inputListContribution = append(inputListContribution, r2)
	r3 := response.Result{Login: "daniele", Total_Contribution: 8, Number_Repositories: 8}
	inputListContribution = append(inputListContribution, r3)
	r4 := response.Result{Login: "Susan", Total_Contribution: 2, Number_Repositories: 2}
	inputListContribution = append(inputListContribution, r4)
	return inputListContribution
}

// Load generic users infos (this could be generated randomly)
func loadGenericInputMaps() (response.ResultsOrderedByContribution, response.ResultsOrderedByRepositories) {
	inputListContribution := loadGenericInputMap()
	var outputListOrderedByRepositories response.ResultsOrderedByRepositories = make([]response.Result, 0)

	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[2])
	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[0])
	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[3])
	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[1])
	return inputListContribution, outputListOrderedByRepositories
}

// Load generic users infos (this could be generated randomly)
func loadGenericInputMaps2() (response.ResultsOrderedByContribution, response.ResultsOrderedByRepositories) {
	inputListContribution := loadGenericInputMap()
	var outputListOrderedByRepositories response.ResultsOrderedByRepositories = make([]response.Result, 0)

	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[2])
	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[0])
	outputListOrderedByRepositories = append(outputListOrderedByRepositories, inputListContribution[1])
	return inputListContribution, outputListOrderedByRepositories
}

func TestExerciseOrderByRepositories(t *testing.T) {
	inputListContribution, outputListOrderedByRepositories := loadGenericInputMaps()
	inputListContribution2, outputListOrderedByRepositories2 := loadGenericInputMaps2()
	tests := map[string]structTestingOrderByRepositories{
		// This user has done several prs but just 5 contributions in repos (4 in total) it doesn't own
		"first_test": {
			inputListContribution:           inputListContribution,
			outputListOrderedByRepositories: outputListOrderedByRepositories,
			limitCount:                      4,
		},
		"second_test": {
			inputListContribution:           inputListContribution2,
			outputListOrderedByRepositories: outputListOrderedByRepositories2,
			limitCount:                      3,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			outputListOrderedByRepositories := service.OrderByRepositories(tc.inputListContribution, tc.limitCount)
			for i, user := range outputListOrderedByRepositories {
				if user.Login != tc.outputListOrderedByRepositories[i].Login {
					t.Fatalf("test failed expected %s got %s", tc.outputListOrderedByRepositories[i].Login, user.Login)
				}

			}

		})
	}
}
