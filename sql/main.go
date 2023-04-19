package main

/*
Provide GET functionality for the following endpoints:
/package/{id}
/package/{id}/rate
/package/byName/{name}
*/

/*
Provide POST functionality for the following endpoints:
/packages
/package/byRegEx
*/

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	//"reflect"
	//"regexp"
	"github.com/Masterminds/semver"
	"github.com/golang-jwt/jwt"
	"github.com/packit461/packit23/sql/models"

	//"github.com/packit461/packit23/package_rater/internal/logger"

	"cloud.google.com/go/cloudsqlconn"
	// "cloud.google.com/go/cloudsqlconn/mysql/mysql"
	sql_driver "github.com/go-sql-driver/mysql"
)

func connect_test_db() (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		"db_user:oldpassword!!!@tcp(127.0.0.1:3306)/test_db",
	)
	if err != nil {
		log.Fatal(err)
		print("FATAL")
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return db, nil
}

func connect() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.", k)
		}
		return v
	}

	var (
		dbUser   = mustGetenv("DB_USER")
		dbPwd    = mustGetenv("DB_PASSWORD")
		dbName   = mustGetenv("DB_NAME")
		project  = mustGetenv("PROJECT_ID")
		region   = mustGetenv("REGION")
		instance = mustGetenv("INSTANCE_NAME")
	)

	var instanceConnectionName = project + region + instance

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}
	var opts []cloudsqlconn.DialOption
	sql_driver.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instanceConnectionName, opts...)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return dbPool, nil
}

func return_500_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal error"))
}

func return_404_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not found"))
}

func return_400_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 - There is missing field(s) in the PackageQuery/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
}

func return_413_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	w.Write([]byte("413 - Request entity too large"))
}

func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["X-Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return_400_packet(writer, request)
				}
				return "", nil
			})
			if err != nil {
				return_400_packet(writer, request)
			}
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				return_400_packet(writer, request)
			}
		}
	})
}

// return all versions of package name in db
func getMetadataFromName(db *sql.DB, query models.PackageQuery) ([]models.PackageMetadata, error) {
	var metadataList []models.PackageMetadata
	rows, err := db.Query("SELECT ID, NAME, VERSION FROM Registry WHERE NAME = ?;", query.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var md models.PackageMetadata
		if err := rows.Scan(&md.ID, &md.Name, &md.Version); err != nil {
			return nil, fmt.Errorf("version of package not found. rows.Scan: %v", err)
		}
		metadataList = append(metadataList, md)
	}
	return metadataList, nil
}

/*
// Get the packages from the registry
Missing:
- Pagination
*/
func handle_packages(w http.ResponseWriter, r *http.Request) {
	//db, err := connect()
	db, err := connect_test_db()
	if err != nil {
		log.Fatal(err)
		return_500_packet(w, r)
	}

	//logger.Info(fmt.Sprintf("Received %s request", r.Method))
	// parse query for offset (pagination). if empty, return the first page of results
	query := r.URL.Query()
	offset := query.Get("offset")
	if offset == "" {
		offset = "1"
	}

	// parse body for versions to find
	var response_arr []models.PackageQuery
	var body models.PackagesBody
	var ret []models.PackageMetadata
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return_404_packet(w, r)
		log.Fatal("\nError reading body of request\n")
	}
	for _, q := range body.Items {
		response_arr = append(response_arr, models.PackageQuery{q.Name, q.Version})
	}

	for _, response := range response_arr {
		// make version field a range to look for
		c, err := semver.NewConstraint(response.Version)
		if err != nil {
			return_400_packet(w, r)
			log.Fatal(err)
		}

		// query all versions of a package if found in db
		metadataList, err := getMetadataFromName(db, response)
		if err != nil {
			return_400_packet(w, r)
			log.Fatal(err)
		}
		// check which version is in range
		for _, md := range metadataList {
			v, err := semver.NewVersion(md.Version)
			if err != nil {
				return_500_packet(w, r)
				log.Fatal(err)
			}
			if c.Check(v) {
				ret = append(ret, md)
			}
		}
	}

	//print("ret: ", ret, "\n")

	if ret == nil {
		return_400_packet(w, r)
	}
	json.NewEncoder(w).Encode(ret)
	w.WriteHeader(200)
}

// Return this package (ID)
func handle_package_id(w http.ResponseWriter, r *http.Request) {
	//db, err := connect()
	db, err := connect_test_db()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// logger.Info(fmt.Sprintf("Received %s request", r.Method))
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}
	// logger.Info(headers)
	vars := mux.Vars(r)
	id := vars["id"]

	//id := r.Header.Get("Id")
	if id == "" {
		// logger.Info("\nNo Matching Value to key Id\n")
		return_404_packet(w, r)
		return
	}
	var meta models.PackageMetadata
	rows, err := db.Query("SELECT ID, NAME, VERSION FROM Registry WHERE ID = " + id + ";")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&meta.ID, &meta.Name, &meta.Version)
		if err != nil {
			log.Fatal(err)
		}
	}
	print("\n", meta.Version, "\n")
	res, err := db.Query(`BEGIN SELECT 
							B.BINARY_FILE, 
							A.URL 
							B.JS_PROGRAM
							FROM Registry AS A
							WHERE A.ID == ` + id + `
							INNER JOIN Binaries AS B
								ON A.BINARY_PK == B.ID
							END;`)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()
	var packData models.PackageData
	// Need to append NULL for JSProgram
	err = res.Scan(&packData.Content, &packData.URL, &packData.JSProgram)
	if err != nil {
		log.Fatal(err)
	}
	totalPack := models.PackageModel{Metadata: &meta, Data: &packData}
	packJson, err := json.Marshal(totalPack)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(packJson)
	w.WriteHeader(200)
}

// Return the rating. Only use this if each metric was computed successfully.
func handle_package_rate(w http.ResponseWriter, r *http.Request) {
	db, err := connect_test_db()
	if err != nil {
		log.Fatal(err)
	}

	// Get the Requested package's ID
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}
	// logger.Info(headers)
	//id := r.Header.Get("Id")
	vars := mux.Vars(r)
	id := vars["id"]
	//print("key is", key)
	if id == "" {
		// logger.Info("\nNo Matching Value to key Id\n")
		return_404_packet(w, r)

		w.WriteHeader(404)
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Query(`BEGIN SELECT 
						A.BUS_FACTOR,
						A.CORRECTNESS,
						A.RAMP_UP,
						A.RESPONSIVENESS,
						A.LICENSE_SCORE,
						A.PINNING_PRACTICE,
						A.PULL_REQUEST,
						A.NET_SCORE
						FROM Rating AS A
						WHERE B.ID == ` + id + `
						INNER JOIN Registry AS B
							ON A.ID == B.RATING_PK
						END;`)
	var ratings models.PackageRating
	err = res.Scan(&ratings.BusFactor, &ratings.Correctness, &ratings.RampUp, &ratings.ResponsiveMaintainer, &ratings.LicenseScore, &ratings.GoodPinningPractice, &ratings.PullRequest, &ratings.NetScore)
	if err != nil {
		return_400_packet(w, r)
		log.Fatal(err)

	}
	ratingsJson, err := json.Marshal(ratings)
	if err != nil {
		return_500_packet(w, r)
		log.Fatal(err)
	}
	w.Write(ratingsJson)
	w.WriteHeader(200)
}

// Return the history of this package (all versions).
/*

type User struct {
	Name string `json:"name"`
	// Is this user an admin?
	IsAdmin bool `json:"isAdmin"`
}

type PackageMetadata struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	// Package version
	Version string `json:"Version"`
}

type PackageHistoryEntry struct {
	User *User `json:"User"`
	// Date of activity using ISO-8601 Datetime standard in UTC format.
	Date time.Time `json:"Date"`

	PackageMetadata *PackageMetadata `json:"PackageMetadata"`

	Action string `json:"Action"`
}
*/

// registry (ID, NAME, RATING_PK, AUTHOR_PK, URL, BINARY_PK, VERSION, UPLOADED, IS_EXTERNAL)
// User - go into Users
// Date - go into Registry(UPLOADED)
// PackageMetadata - go into Registry(ID,NAME,VERSION)
// Action - not ready yet
// db needs to link package name with the User who uploaded and add Action (somewhere)

// return the package history with package name from path
func handle_package_byname(w http.ResponseWriter, r *http.Request) {
	db, err := connect_test_db()
	if err != nil {
		return_500_packet(w, r)
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		return_404_packet(w, r)
	}
	var ret []models.PackageHistoryEntry
	var metadataList []models.PackageMetadata
	var times []time.Time
	// get registry entry from name
	rows, err := db.Query("SELECT ID, NAME, VERSION, UPLOADED FROM Registry WHERE NAME = ?", name)
	if err != nil {
		return_500_packet(w, r)
		log.Fatal(err)
	}
	defer rows.Close()
	// get all versions of named package
	for rows.Next() {
		var time time.Time
		var md models.PackageMetadata
		if err := rows.Scan(&md.ID, &md.Name, &md.Version, &time); err != nil {
			// package with name not found
			return_500_packet(w, r)
			log.Fatal(err)
		}
		metadataList = append(metadataList, md)
		times = append(times, time)
	}

	// iterate through versions of package and get rest of history
	for i, md := range metadataList {
		print("\n")
		var history models.PackageHistoryEntry
		history.User = &models.User{Name: "test", IsAdmin: false}
		history.Date = times[i]
		history.PackageMetadata = &md
		history.Action = "Uploaded"
		print(history.PackageMetadata.Name, history.PackageMetadata.Version, history.PackageMetadata.ID, history.User.Name, history.User.IsAdmin, history.Action)
		print("\n")
	}

	retJson, err := json.Marshal(ret)
	if err != nil {
		return_500_packet(w, r)
		log.Fatal(err)
	}
	w.Write(retJson)
}

func handle_package_byregex(w http.ResponseWriter, r *http.Request) {

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/packages", verifyJWT(handle_packages))
	router.HandleFunc("/package/{id}", verifyJWT(handle_package_id))
	router.HandleFunc("/package/{id}/rate", verifyJWT(handle_package_rate))
	router.HandleFunc("/package/byName/{name}", (handle_package_byname))
	router.HandleFunc("/package/byRegEx", verifyJWT(handle_package_byregex))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}
