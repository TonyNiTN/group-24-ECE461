package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"group-24-ECE461/internal/models"
	"sort"
	"time"

	"github.com/google/go-github/github"
)

func Base64Encode(str string) string { // function to encode a string to base64
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string { // function to decode a base64 encoded string
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

func GetTopFiveContributions(contr []*github.Contributor) int { // function to get the top 5 contributors of a list of contributors on a github repository and returns the sum of contributions from all 5 contributors.
	sort.Slice(contr, func(i, j int) bool {
		return *contr[i].Contributions > *contr[j].Contributions
	})

	topFive := contr[:5]
	var sum int
	for _, c := range topFive {
		sum += *c.Contributions
	}
	return sum
}

func CountCommits(commits []*github.RepositoryCommit) int { // function to count the number of commits in a single page of response. Returns the count as an integer.
	var sum int
	for i := range commits {
		sum = i + 1
	}

	return sum
}

func GetLastWeek() string { // function to get the day 1 week before today. Returns a string in the format YYYY-MM-DD.
	date := time.Now().Add(time.Duration(-24*7) * time.Hour)
	YYYYMMDD := "2006-01-02"
	return date.Format(YYYYMMDD)
}

func PrintRepo(repo *models.Repository) { // function to print the fields of a repository in json format
	b, err := json.MarshalIndent(repo, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Print(string(b))
	}
}
