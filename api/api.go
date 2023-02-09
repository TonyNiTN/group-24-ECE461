package api

import (
	"context"
	"encoding/json"
	"fmt"
	"group-24-ECE461/internal/error"
	"group-24-ECE461/internal/helper"
	"group-24-ECE461/internal/models"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

//TODO:
//identify all the endpoints needed:
//
//Ramp-up Time: parse or use regex on readme to get installation, quickstart, docs and url ;; using REST API
//Correctness: number of stargazers, found in Repository ;; using GraphQL API
//Bus Factor: Total commits, and top 5 contributors ;; using REST API
//Responsiveness: pull requests in the last week to open issues in the last week ;; using REST API
//License Compatibility: 1 If license, 0 otherwise (will use regex if need to search for a specific license) ;; using GraphQL API

func SendRequests(client *github.Client, graphqlClient *githubv4.Client, ctx context.Context, graphqlCtx context.Context, repo *models.Repository, logger *zap.Logger) {
	GetStars(graphqlClient, ctx, repo, logger)
	GetReadme(client, ctx, repo, logger)
	GetLicense(graphqlClient, ctx, repo, logger)
	GetPullRequests(client, ctx, repo, logger)
	GetIssues(client, ctx, repo, logger)
	GetContributors(client, ctx, repo, logger)
	GetCommits(client, ctx, repo, logger)
}

func GetRepoOwnerFromNPM(pack string) string {
	url := "https://registry.npmjs.org/" + pack
	res, err := http.Get(url)
	if err != nil {

		fmt.Println(err)
		return ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		//log err
		return ""
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		//log error
		return ""
	}

	repo := data["repository"]
	value, ok := repo.(map[string]interface{})
	var repo_url string
	if !ok {
		fmt.Println("error reading values from response")
	}
	repo_url = value["url"].(string)
	parts := strings.Split(repo_url, ".com/")
	parts = strings.Split(parts[1], "/")
	return parts[0]
}

func CreateRESTClient() (*github.Client, context.Context) { // function to create github REST api client
	ctx := context.Background()                                                           // create empty context
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")}) // configure auth header for the client
	tc := oauth2.NewClient(ctx, ts)                                                       // create new http client
	client := github.NewClient(tc)                                                        // create new github rest api client from the http client template
	return client, ctx                                                                    // returns the github rest api client and the empty context
}

func CreateGQLClient() (*githubv4.Client, context.Context) { // function to creategithub GraphQL api client
	ctx := context.Background()                                                            // create empty context
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")}) // configure auth header for the client
	tc := oauth2.NewClient(ctx, ts)                                                        // create new http client
	graphqlClient := githubv4.NewClient(tc)                                                // create new github graphql api client from the http client template
	return graphqlClient, ctx                                                              // returns the github graphql api client and the empty context
}

func GetPullRequests(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get request for pull requests
	since := helper.GetLastWeek()                                                              // get date to filter results by
	s := fmt.Sprintf("org:%s repo:%s created:>%s is:pr is:open", repo.Owner, repo.Name, since) // create query string
	prs, _, err := client.Search.Issues(ctx, s, &github.SearchOptions{})                       // make the request
	if err != nil {
		fmt.Println(err)
		logger.Debug("Error getting pull requests from the server")
		return
	} else {
		logger.Info("Success -- Pull requests")
		repo.OpenPRs = *prs.Total // populate repository field
	}
}

func GetIssues(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get requests for issues
	since := helper.GetLastWeek()                                                              // get date to filter results by
	s := fmt.Sprintf("org:%s repo:%s created:>%s is:pr is:open", repo.Owner, repo.Name, since) // create query string
	issues, _, err := client.Search.Issues(ctx, s, &github.SearchOptions{})                    // make the request
	if err != nil {
		fmt.Println(err)
		logger.Debug("Error getting issues from the server")
		return
	} else {
		logger.Info("Success -- Issues")
		repo.OpenIssues = *issues.Total // populate repository field
	}
}

func GetCommits(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get requests for commits
	opt := &github.CommitsListOptions{ // define options structure to indicate results per page
		ListOptions: github.ListOptions{PerPage: 30}, // here, it's 30 results per page, as default.
	}

	_, res, err := client.Repositories.ListCommits(ctx, repo.Owner, repo.Name, opt) // make the request
	if err != nil {
		fmt.Println(err)
		logger.Debug("Error getting commits from the server")
		return
	} else {
		logger.Info("Success -- Commits")
		repo.Commits = res.LastPage * 30 // populate repository field
	}
}

func GetContributors(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get requests for contributors
	contr, _, err := client.Repositories.ListContributors(ctx, repo.Owner, repo.Name, nil) // make the requests

	if err != nil {
		fmt.Println(err)
		logger.Debug("Error getting contributors from the server")
		return
	} else {
		logger.Info("Success -- Contributors")
		repo.TopContributions = helper.GetTopFiveContributions(contr) // populate repository field with the total contributions of top 5 contributors
	}
}

func GetLicense(client *githubv4.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get requests for license
	variables := map[string]interface{}{ // variables to dynamically populate the graphql query structure
		"owner": githubv4.String(repo.Owner),
		"name":  githubv4.String(repo.Name),
	}
	err := client.Query(ctx, &models.License, variables) // make the graphql request
	if err != nil {
		fmt.Println(err)
		logger.Debug("Error getting License from the server")
		return
	} else {
		logger.Info("Success -- License")
		repo.License = models.License.Repository.LicenseInfo.Body // populate repository field
	}
}

func GetStars(client *githubv4.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get requests for stargazers
	variables := map[string]interface{}{ // variables to dynamically populate the graphql query structure
		"owner": githubv4.String(repo.Owner),
		"name":  githubv4.String(repo.Name),
	}
	err := client.Query(ctx, &models.Stars, variables) // make the graphql request
	if err != nil {
		newError := error.NewGraphQLError("query failed", "query.getStars")
		logger.Debug(newError.Error())
		return
	}
	fmt.Println("success")
	logger.Info("Success -- Stargazers")
	repo.StarsCount = models.Stars.Repository.StargazerCount // populate repository field

}

func GetReadme(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger) { // function to make get requests for readme
	readme, _, err := client.Repositories.GetReadme(ctx, repo.Owner, repo.Name, &github.RepositoryContentGetOptions{}) // make the requests
	if err != nil {
		newError := error.NewGraphQLError("query failed", "query.getStars")
		logger.Error(newError.Error())
		return
	} else {
		logger.Info("Success -- Readme")
		repo.Readme = helper.Base64Decode(*readme.Content) // populate repository field
	}
}

func ParseUrl(url string) (owner string, name string) {
	res, _ := regexp.MatchString(`(?i)github\b`, url)
	if res {
		owner, name = helper.GetOwnerAndName(url)
	}
	res, _ = regexp.MatchString(`(?i)npmjs\b`, url)
	if res {
		name = helper.GetPackageName(url)
		owner = GetRepoOwnerFromNPM(name)
	}

	return owner, name
}
