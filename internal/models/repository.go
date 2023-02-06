package models

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
