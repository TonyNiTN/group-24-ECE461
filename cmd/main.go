package main

import (
	"fmt"
	"group-24-ECE461/api"
	"group-24-ECE461/internal/error"
	"group-24-ECE461/internal/helper"
	"group-24-ECE461/internal/logger"
	"group-24-ECE461/internal/models"
	"group-24-ECE461/internal/scorer"
)

func main() {
	logger, err := logger.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	newError := error.NewGraphQLError("query failed", "query.getUser")

	logger.Info(newError.Error())

	logger.Info("Starting Application")
	owner := api.GetRepoOwnerFromNPM("express")
	fmt.Println(owner)
	o, n := helper.GetOwnerAndName("https://github.com/expressjs/express")
	fmt.Println(o, n)
	var repos []*models.Repository
	client, ctx := api.CreateRESTClient()
	graphqlClient, graphqlCtx := api.CreateGQLClient()
	for owner, name := range models.Repos {
		repo := models.NewRepository()
		repo.Name = name
		repo.Owner = owner
		api.SendRequests(client, graphqlClient, ctx, graphqlCtx, repo, logger)
		scorer.CalculatePackageScore(repo)
		repos = append(repos, repo)
	}
	//helper.DisplayResults(repos)
}
