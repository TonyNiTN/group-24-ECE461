package models

import (
	"encoding/json"
	"fmt"
	"group-24-ECE461/internal/logger"
	"sort"
)

type Repository struct {
	Name                      string  `json:"name"`
	Owner                     string  `json:"owner"`
	StarsCount                int     `json:"starsCount"`       // Used in the calculation of correctness score
	OpenIssues                int     `json:"openIssues"`       // Used in the calculation of responsiveness score
	OpenPRs                   int     `json:"openPullRequests"` // Used in the calculation of responsiveness score
	TopContributions          int     `json:"topContributions"` // Used in the calculation of bus factor score
	Commits                   int     `json:"commits"`          // Used in the calculation of bus factor score
	License                   string  `json:"license"`          // Used in the calculation of license compatibility score
	Readme                    string  `json:"readme"`           // Used in the calculation of ramp up time score
	RampUpTimeScore           float64 `json:"rampUpTimeScore"`
	CorrectnessScore          float64 `json:"correctnessScore"`
	BusFactorScore            float64 `json:"busFactorScore"`
	ResponsivenessScore       float64 `json:"responsivenessScore"`
	LicenseCompatibilityScore float64 `json:"licenseCompatibilityScore"`
	NetScore                  float64 `json:"netScore"`
	NetPercentage             float64 `json:"netPercentage"`
}

func NewRepository() *Repository { // Initialize empty *Repository object
	return &Repository{}
}

var Weights = map[string]float64{
	"Ramp Up Time":          2.0,
	"Correctness":           2.0,
	"Bus Factor":            3.0,
	"Responsiveness":        2.0,
	"License Compatibility": 2.0,
}

func PrintRepo(repo *Repository) { // function to print the fields of a repository in json format
	b, err := json.MarshalIndent(repo, "", " ")
	if err != nil {
		logger, _ := logger.InitLogger()
		logger.Error("Error printing repository")
		fmt.Println(err)
		return
	} else {
		fmt.Print(string(b))
	}
}

func SortRepositories(repos []*Repository) []*Repository {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].NetScore > repos[j].NetScore
	})

	return repos
}

func DisplayResults(repos []*Repository) {
	sorted_repos := SortRepositories(repos)
	fmt.Println("Name             | Place | Ramp-Up Time | Correctness | Bus Factor | Responsiveness | License Compatibility | Net Score | Net %")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i, repo := range sorted_repos {
		fmt.Printf("%-17s|%-7d|%-14.2f|%-13.2f|%-12.2f|%-16.2f|%-23.2f|%-11.2f|%-6.2f\n", repo.Name, i+1, repo.RampUpTimeScore, repo.CorrectnessScore, repo.BusFactorScore, repo.ResponsivenessScore, repo.LicenseCompatibilityScore, repo.NetScore, repo.NetPercentage)
	}

}
