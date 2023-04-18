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
	//"reflect"
	"regexp"

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

// is there any case where packages are to be looked for in external registries?
// return a list of packages (id, name, version) found in a database matching the versioning semantics of the given version
func getPackagesExact(db *sql.DB, query models.PackageQuery) ([]models.PackageMetadata, error) {
	// find versions in db
	rows, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry WHERE NAME = ` + string(query.Name) + ` AND VERSION = ` + string(query.Version) + `;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var metadata models.PackageMetadata
	for rows.Next() {
		if err := rows.Scan(&metadata.ID, &metadata.Name, &metadata.Version); err != nil {
			return nil, fmt.Errorf("version of package not found. rows.Scan: %v", err)
		}
	}

	var ret []models.PackageMetadata
	ret = append(ret, metadata)
	return ret, nil
}

// get a list of packages from db given a bounded range of versions
func getPackagesRange(db *sql.DB, query models.PackageQuery) ([]models.PackageMetadata, error) {
	boundedRangeRegex := regexp.MustCompile(`(\d*\.\d*\.\d*)-(\d*\.\d*\.\d*)`)
	boundedRangeMatch := boundedRangeRegex.FindStringSubmatch(query.Version)

	begin := boundedRangeMatch[1]
	end := boundedRangeMatch[2]

	// find versions in db
	rows, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry WHERE NAME = ` + query.Name + ` AND VERSION >= ` + begin + `AND VERSION <= ` + end + `;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var metadataList []models.PackageMetadata
	for rows.Next() {
		var metadata models.PackageMetadata
		if err := rows.Scan(&metadata.ID, &metadata.Name, &metadata.Version); err != nil {
			return nil, fmt.Errorf("version of package not found. rows.Scan: %v", err)
		}
		metadataList = append(metadataList, metadata)
	}

	return metadataList, nil
}

// only major version must match. Any minor or patch version greater than or equal to the minimum is valid.
// are versions getting correctly checked? since comparisons are with strings
func getPackagesCarat(db *sql.DB, query models.PackageQuery) ([]models.PackageMetadata, error) {
	caratRegex := regexp.MustCompile(`\^(\d*\.\d*\.\d*)`)
	caratMatch := caratRegex.FindStringSubmatch(query.Version)

	begin := caratMatch[1]
	// get major version and increment it
	end := string(int(begin[0])+1) + ".0.0"

	// find versions in db
	rows, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry WHERE NAME = ` + query.Name + ` AND VERSION >= ` + begin + ` AND VERSION < ` + end + `;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var metadataList []models.PackageMetadata
	for rows.Next() {
		var metadata models.PackageMetadata
		if err := rows.Scan(&metadata.ID, &metadata.Name, &metadata.Version); err != nil {
			return nil, fmt.Errorf("version of package not found. rows.Scan: %v", err)
		}
		metadataList = append(metadataList, metadata)
	}

	return metadataList, nil
}

// For tilde ranges, major and minor versions must match those specified, but any patch version greater than or equal to the one specified is valid.
func getPackagesTilde(db *sql.DB, query models.PackageQuery) ([]models.PackageMetadata, error) {
	tildeRegex := regexp.MustCompile(`\~(\d*\.\d*\.\d*)`)
	tildeMatch := tildeRegex.FindStringSubmatch(query.Version)

	begin := tildeMatch[1]

	// get minor version and increment it
	end := string(begin[0]) + "." + string(int(begin[2])+1) + ".0"

	// find versions in db
	rows, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry WHERE NAME = ` + query.Name + ` AND VERSION >= ` + begin + ` AND VERSION < ` + end + `;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var metadataList []models.PackageMetadata
	for rows.Next() {
		var metadata models.PackageMetadata
		if err := rows.Scan(&metadata.ID, &metadata.Name, &metadata.Version); err != nil {
			return nil, fmt.Errorf("version of package not found. rows.Scan: %v", err)
		}
		metadataList = append(metadataList, metadata)
	}

	return metadataList, nil
}

// Get the packages from the registry
/*
Missing:
- Authentication
- Pagination
- Handling the case where "*" is passed as PackageQuery (enumerate all packages)
- Handling the case where a package is not found
*/
func handle_packages(w http.ResponseWriter, r *http.Request) {
	// check authentication
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}

	//logger.Info(headers)
	//db, err := connect()
	db, err := connect_test_db()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("\nError reading body of request\n")
		return_404_packet(w, r)
	}
	for _, q := range body.Items {
		response_arr = append(response_arr, models.PackageQuery{q.Name, q.Version})
	}

	// https://semver.org/
	//validVersion := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	exactRegex := regexp.MustCompile(`(\d*\.\d*\.\d*)`)
	boundedRangeRegex := regexp.MustCompile(`(\d*\.\d*\.\d*)-(\d*\.\d*\.\d*)`)
	caratRegex := regexp.MustCompile(`\^(\d*\.\d*\.\d*)`)
	tildeRegex := regexp.MustCompile(`\~(\d*\.\d*\.\d*)`)

	// get requested versions of packages
	// TODO: ADD PAGINATION STUFF
	for _, response := range response_arr {
		//if validVersion.FindStringSubmatch(response.Version) != nil {
		if boundedRangeRegex.FindStringSubmatch(response.Version) != nil {
			list, err := getPackagesRange(db, response)
			if err != nil {
				log.Fatal(err)
				return_404_packet(w, r)
			}
			ret = append(ret, list...)
		} else if caratRegex.FindStringSubmatch(response.Version) != nil {
			list, err := getPackagesCarat(db, response)
			if err != nil {
				log.Fatal(err)
				return_404_packet(w, r)
			}
			ret = append(ret, list...)
		} else if tildeRegex.FindStringSubmatch(response.Version) != nil {
			list, err := getPackagesTilde(db, response)
			if err != nil {
				log.Fatal(err)
				return_404_packet(w, r)
			}
			ret = append(ret, list...)
			// checking exact version last because this regex is a subset of the other regexes
		} else if exactRegex.FindStringSubmatch(response.Version) != nil {
			list, err := getPackagesExact(db, response)
			if err != nil {
				log.Fatal(err)
				return_404_packet(w, r)
			}
			ret = append(ret, list...)
		} else {
			print("goes here")
			log.Fatal("Invalid version")
			return_400_packet(w, r)
		}
		// } else {
		// 	log.Fatal("Invalid version")
		// 	return_400_packet(w, r)
		// }
	}

	print("ret: ", ret, "\n")

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
		log.Fatal(err)
		w.WriteHeader(400)
	}
	ratingsJson, err := json.Marshal(ratings)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
	}
	w.Write(ratingsJson)
	w.WriteHeader(200)
}

// Return the history of this package (all versions).
func handle_package_byname(w http.ResponseWriter, r *http.Request) {
	// db, err := connect_test_db()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// vars := mux.Vars(r)
	// id := vars["name"]

}

func handle_package_byregex(w http.ResponseWriter, r *http.Request) {

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/packages", handle_packages)
	router.HandleFunc("/package/{id}", handle_package_id)
	router.HandleFunc("/package/{id}/rate", handle_package_rate)
	router.HandleFunc("/package/byName/{name}", handle_package_byname)
	router.HandleFunc("/package/byRegEx", handle_package_byregex)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}
