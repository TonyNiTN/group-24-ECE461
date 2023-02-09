package cli

import (
	"fmt"
	"group-24-ECE461/api"
	"group-24-ECE461/internal/models"
	"group-24-ECE461/internal/scorer"
	"os/exec"

	"go.uber.org/zap"
)

func Install(logger *zap.Logger) {
	fmt.Println("Installing....")
	cmd := exec.Command("go", "install")
	err := cmd.Run()
	if err != nil {
		logger.Debug("Error installing program!")
		fmt.Println("Error compiling program:", err)
		return
	}
	fmt.Println("Installation successful")
}

func Score(urlLinks []string, logger *zap.Logger) {
	fmt.Println("Scoring.....")
	var repos []*models.Repository
	client, ctx := api.CreateRESTClient()
	graphqlClient, graphqlCtx := api.CreateGQLClient()
	for _, url := range urlLinks {
		owner, name := api.ParseUrl(url)
		repo := models.NewRepository()
		repo.Name = name
		repo.Owner = owner
		api.SendRequests(client, graphqlClient, ctx, graphqlCtx, repo, logger)
		scorer.CalculatePackageScore(repo)
		repos = append(repos, repo)
	}
	models.DisplayResults(repos)
}

func Build(logger *zap.Logger) {
	fmt.Println("Building......")
	cmd := exec.Command("go", "build")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		logger.Debug("Build Error!")
		return
	} else {
		logger.Info("Build successful")
		fmt.Println("Build successful")
	}
}

func Test(logger *zap.Logger) {
	fmt.Println("Testing.....")
	cmd := exec.Command("go", "test", "./...", "-cover")
	err := cmd.Run()
	if err != nil {
		logger.Debug("Error running tests!")
		fmt.Println("Error running tests:", err)
		return
	}
}
