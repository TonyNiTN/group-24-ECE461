package cli

import (
	"fmt"
	"group-24-ECE461/api"
	"group-24-ECE461/internal/models"
	"group-24-ECE461/internal/scorer"
	"strings"
	"os/exec"

	"go.uber.org/zap"
)

func Install(logger *zap.Logger) {
	out, err := exec.Command("go", "list", "-m", "all").Output()
	if err != nil {
		fmt.Println("Error executing 'go list -m all':", err)
		return
	}
	dependencies := strings.Split(string(out), "\n")
	fmt.Printf("Number of dependencies: %d\n", len(dependencies)-1)
}

func Score(urlLinks []string, logger *zap.Logger) {
	if len(urlLinks) == 0 {
		logger.Info("Input file is empty, No Url to read!")
		fmt.Println("Input file is empty, No Url to read!")
		return
	}
	fmt.Println("Scoring.....")
	var repos []*models.Repository
	client, ctx := api.CreateRESTClient()
	graphqlClient, graphqlCtx := api.CreateGQLClient()

	for _, url := range urlLinks {
		owner, name := api.ParseUrl(url)
		if owner == "" || name == "" {
			fmt.Println("Error parsing Url string (Invalid Url type)")
		} else {
			repo := models.NewRepository()
			repo.Name = name
			repo.Owner = owner
			repo.Url = url
			flag := api.SendRequests(client, graphqlClient, ctx, graphqlCtx, repo, logger)
			if flag != 0 {
				continue
			}
			
			scorer.CalculatePackageScore(repo)
			repos = append(repos, repo)
			
		}
	}
	models.ShowResults(repos)
}

// func Test(logger *zap.Logger) {
// 	fmt.Println("Testing.....")
// 	cmd := exec.Command("go", "test", "./...", "-cover")
// 	fmt.Println(cmd)
// 	err := cmd.Run()
// 	fmt.Println(err)
// 	if err != nil {
// 		logger.Debug("Error running tests!")
// 		fmt.Println("Error running tests:", err)
// 		return
// 	}
// }
