package github_apis

import (
	"fmt"
	"github.com/packit461/packit23/package_rater/internal/logger"
	"github.com/packit461/packit23/package_rater/internal/models"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

var OWNER string = "nodejs"
var NAME string = "node"

func TestCreateRESTClient(t *testing.T) { // test CreateRESTClient function in the api package
	client, ctx := CreateRESTClient()

	if client == nil || ctx == nil {
		t.Error("Error creating REST client!")
	}
}

func TestCreateGQLClient(t *testing.T) { // test CreateGQLClient function in the api package
	client, ctx := CreateGQLClient()

	if client == nil || ctx == nil {
		t.Error("Error creating GraphQL client!")
	}
}

// func TestCreateBuildDependenciesQuery(t *testing.T) {
// 	c := cache.New(5*time.Minute, 10*time.Minute)

// 	client, ctx := CreateGQLClient()
// 	repo := models.NewRepository()
// 	repo.Owner = "nodejs" /// testing
// 	repo.Name = "node"
// 	logger, _ := logger.InitLogger()
// 	query := buildDependenciesQuery(client, ctx, repo, logger, c)

// 	if ok := len(query) > 0; !ok {
// 		fmt.Println("the length is 0")
// 		t.Error("Error making dependencies query!")
// 	} else {
// 		// bs, _ := json.Marshal(query)
// 		fmt.Println(query)
// 	}

// }

func TestGetPullRequests(t *testing.T) { // test GetPullRequests function in the api package
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateRESTClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetPullRequests(client, ctx, repo, logger, c)
	val := interface{}(repo.OpenPRs)
	if _, ok := val.(int); !ok {
		t.Error("Error getting pull requests!")
	}
}

func TestGetIssues(t *testing.T) { // test GetIssues function in the api package
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateRESTClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetIssues(client, ctx, repo, logger, c)
	val := interface{}(repo.OpenIssues)
	if _, ok := val.(int); !ok {
		t.Error("Error getting issues!")
	}
}

func TestGetCommits(t *testing.T) { // test GetCommits function in the api package
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateRESTClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetCommits(client, ctx, repo, logger, c)
	val := interface{}(repo.Commits)
	if _, ok := val.(int); !ok {
		t.Error("Error getting commits!")
	}
}

func TestGetContributors(t *testing.T) { // test GetContributors function in the api package
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateRESTClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetContributors(client, ctx, repo, logger, c)
	val := interface{}(repo.TopContributions)
	if _, ok := val.(int); !ok {
		t.Error("Error getting top 5 contributions!")
	}
}

func TestGetStars(t *testing.T) { // test GetStars function in the api package
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateGQLClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetStars(client, ctx, repo, logger, c)
	val := interface{}(repo.StarsCount)
	if _, ok := val.(int); !ok {
		t.Error("Error getting stargazers count!")
		fmt.Println()
	}
	// fmt.Println(repo.StarsCount)
}

func TestGetDependencyQuery(t *testing.T) {
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateGQLClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetDependencyQuery(client, ctx, repo, logger, c)
	val := interface{}(repo.DependencyCount)
	if _, ok := val.(int); !ok {
		t.Error("Error getting dependency count!")
		fmt.Println()
	}
	fmt.Printf("%T", val)
}

func TestGetReadme(t *testing.T) { // test GetReadme function in the api package
	c := cache.New(5*time.Minute, 10*time.Minute)

	client, ctx := CreateRESTClient()
	repo := models.NewRepository()
	var OWNER string = "nodejs"
	var NAME string = "node"
	repo.Owner = OWNER
	repo.Name = NAME
	logger, _ := logger.InitLogger()
	GetReadme(client, ctx, repo, logger, c)
	val := interface{}(repo.Readme)
	if _, ok := val.(string); !ok {
		t.Error("Error getting readme")
	}
}
