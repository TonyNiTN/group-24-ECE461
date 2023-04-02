package models

var License struct { // graphql query structure to get License body of a repository
	Repository struct {
		LicenseInfo struct {
			Body string `json:"body"`
		}
	} `graphql:"repository(owner: $owner, name: $name)"` // passing in owner and name for the repository to look for
}

var Stars struct { // graphql query structure to get stargazers count of a repository
	Repository struct {
		StargazerCount int `json:"stars"`
	} `graphql:"repository(owner: $owner, name: $name)"` // passing in owner and name for the repository to look for
}

var Dependency struct {
	Repository struct {
		DependencyGraphManifests struct {
			TotalCount int `json:"totalCount"`
			Nodes      []struct {
				Filename string `json:"filename"`
			}
			Edges []struct {
				Node struct {
					BlobPath     string `json:"blobPath"`
					Dependencies struct {
						TotalCount int `json:"dependencyCount"`
						Nodes      []struct {
							PackageName     string `json:"packageName"`
							Requirements    string `json:"requirements"`
							HasDependencies bool   `json:"hasDependencies"`
							PackageManager  string `json:"packageManager"`
						}
					}
				}
			}
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}
