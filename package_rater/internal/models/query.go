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
			TotalCount int
			Nodes      struct {
				FileName string
			}
			Edges struct {
				Node struct {
					BlobPath     string
					Dependencies struct {
						TotalCount int `json:"dependencyCount"`
						Nodes      struct {
							PackageName     string
							Requirements    string
							HasDependencies bool
							PackageManager  string
						}
					}
				}
			}
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}
