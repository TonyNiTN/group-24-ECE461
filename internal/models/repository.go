package models

import (
	"fmt"
	"sort"
	"encoding/json"
)

type Repository struct {
	Name                      string  `json:"name"`
	Owner                     string  `json:"owner"`
	Url string `json:"Url"`
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


type Package struct {
	Url string `json:"Url"`
	NetScore                  float64 `json:"netScore"`
	NetPercentage             float64 `json:"netPercentage"`
	RampUpTimeScore           float64 `json:"rampUpTimeScore"`
	CorrectnessScore          float64 `json:"correctnessScore"`
	BusFactorScore            float64 `json:"busFactorScore"`
	ResponsivenessScore       float64 `json:"responsivenessScore"`
	LicenseCompatibilityScore float64 `json:"licenseCompatibilityScore"`

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

func SortRepositories(repos []*Repository) []*Repository {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].NetScore > repos[j].NetScore
	})

	return repos
}

func ShowResults(repos []*Repository) {
	var data []Package

	for _, repo := range repos {
		m := Package{
			Url: repo.Url,
			NetScore : float64(int(repo.NetScore*100)) / 100,
			NetPercentage: float64(int(repo.NetPercentage*100)) / 100,
			RampUpTimeScore : float64(int(repo.RampUpTimeScore*100)) / 100,
			CorrectnessScore : float64(int(repo.CorrectnessScore*100)) / 100,
			BusFactorScore : float64(int(repo.BusFactorScore*100)) / 100,
			ResponsivenessScore : float64(int(repo.ResponsivenessScore*100)) / 100,
			LicenseCompatibilityScore : float64(int(repo.LicenseCompatibilityScore*100)) / 100,
		}

		data = append(data, m)
	}

	for _, obj := range data {
		b, err := json.Marshal(obj)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(string(b))
	}
}
