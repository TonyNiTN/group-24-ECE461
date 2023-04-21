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
	"cloud.google.com/go/storage"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/packit461/packit23/sql/models"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

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
		if request.Header["X-Authorization"] != nil {
			token, err := jwt.Parse(request.Header["X-Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return_400_packet(writer, request)
					log.Print("Error validating JWT")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				return_400_packet(writer, request)
				log.Print("Error in verifyJWT")
			}
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				return_400_packet(writer, request)
				log.Print("Error in verifyJWT")
			}
		}
	})
}

func getBucketObject(bucketName string, objectName string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
	}
	defer rc.Close()
	bucketObjet, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return bucketObjet, nil
}

// return all versions of package name in db
func getMetadataFromName(db *sql.DB, query models.PackageQuery) ([]models.PackageMetadata, error) {
	var metadataList []models.PackageMetadata
	rows, err := db.Query("SELECT ID, NAME, VERSION FROM Registry WHERE NAME = ?;", query.Name)
	if err != nil {
		log.Print(err)
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
		log.Print(err)
		return_500_packet(w, r)
	}
	defer db.Close()
	// parse query for offset (pagination). if empty, return the first page of results
	query := r.URL.Query()
	offset := query.Get("offset")
	if offset == "" {
		offset = "1"
	}

	// parse body for versions to find
	var response_arr []models.PackageQuery
	var body models.PackagesBody
	var packages_metadata []models.PackageMetadata
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return_404_packet(w, r)
		log.Print("\nError reading body of request\n")
	}
	for _, q := range body.Items {
		response_arr = append(response_arr, models.PackageQuery{q.Name, q.Version})
	}

	// make version field a range to look for
	for _, response := range response_arr {
		// check if response version field contains '-' character without surrounding whitespace, if it doesn't add it
		char_idx := strings.Index(response.Version, "-")
		if strings.Contains(response.Version, "-") && response.Version[char_idx-1] != ' ' && response.Version[char_idx+1] != ' ' {
			response.Version = strings.Replace(response.Version, "-", " - ", -1)
		}
		c, err := semver.NewConstraint(response.Version)
		if err != nil {
			return_400_packet(w, r)
			log.Print(err)
		}

		// query all versions of a package if found in db
		metadataList, err := getMetadataFromName(db, response)
		if err != nil {
			return_400_packet(w, r)
			log.Print(err)
		}
		// check which version is in range
		for _, md := range metadataList {
			v, err := semver.NewVersion(md.Version)
			if err != nil {
				return_500_packet(w, r)
				log.Print(err)
			}
			if c.Check(v) {
				packages_metadata = append(packages_metadata, md)
			}
		}
	}

	if packages_metadata == nil {
		return_400_packet(w, r)
	}
	json.NewEncoder(w).Encode(packages_metadata)
	w.WriteHeader(200)
}

// Return this package (ID)
// Get from the bucket
func handle_package_id(w http.ResponseWriter, r *http.Request) {
	//db, err := connect()
	db, err := connect_test_db()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return_404_packet(w, r)
		log.Print("Emppty {id} in path")
	}
	var meta models.PackageMetadata
	rows, err := db.Query("SELECT ID, NAME, VERSION FROM Registry WHERE ID = " + id + ";")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&meta.ID, &meta.Name, &meta.Version)
		if err != nil {
			log.Print(err)
		}
	}
	// res, err := db.Query(`BEGIN SELECT
	// 						B.BINARY_FILE,
	// 						A.URL
	// 						B.JS_PROGRAM
	// 						FROM Registry AS A
	// 						WHERE A.ID == ?
	// 						INNER JOIN Binaries AS B
	// 							ON A.BINARY_PK == B.ID
	// 						END;`, id)
	// bucket name, object name
	getBucketObject()

	if err != nil {
		log.Print(err)
	}
	defer res.Close()
	var packData models.PackageData
	// Need to append NULL for JSProgram
	err = res.Scan(&packData.Content, &packData.URL, &packData.JSProgram)
	if err != nil {
		log.Print(err)
	}
	totalPack := models.PackageModel{Metadata: &meta, Data: &packData}
	packJson, err := json.Marshal(totalPack)
	if err != nil {
		log.Print(err)
	}
	w.Write(packJson)
	w.WriteHeader(200)
}

// Return the package rating with ID (id from path)
func handle_package_rate(w http.ResponseWriter, r *http.Request) {
	db, err := connect_test_db()
	if err != nil {
		log.Print(err)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return_404_packet(w, r)
		w.WriteHeader(404)
		log.Print("id cannot be empty")
	}
	res, err := db.Query("SELECT A.BUS_FACTOR, A.CORRECTNESS, A.RAMP_UP, A.RESPONSIVENESS, A.LICENSE_SCORE, A.PINNING_PRACTICE, A.PULL_REQUEST, A.NET_SCORE FROM Ratings AS A INNER JOIN Registry AS B ON A.ID = B.RATING_PK WHERE B.ID = ?", id)
	if err != nil {
		return_400_packet(w, r)
		log.Print(err)
	}
	var ratings models.PackageRating
	for res.Next() {
		err = res.Scan(&ratings.BusFactor, &ratings.Correctness, &ratings.RampUp, &ratings.ResponsiveMaintainer, &ratings.LicenseScore, &ratings.GoodPinningPractice, &ratings.PullRequest, &ratings.NetScore)
		if err != nil {
			return_400_packet(w, r)
			log.Print(err)
		}
	}
	ratingsJson, err := json.Marshal(ratings)
	if err != nil {
		return_500_packet(w, r)
		log.Print(err)
	}
	w.Write(ratingsJson)
	w.WriteHeader(200)
}

// return the package history with package name from path (all versions)
// const mysqlFormat = "2006-01-02 15:04:05"
// const timeFormat = "2006-01-02T15:04:05Z"
func handle_package_byname(w http.ResponseWriter, r *http.Request) {
	db, err := connect_test_db()
	if err != nil {
		return_500_packet(w, r)
		log.Print(err)
	}
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		return_404_packet(w, r)
	}
	var ret []models.PackageHistoryEntry
	var metadataList []models.PackageMetadata
	var times []string
	// get registry entry from name
	rows, err := db.Query("SELECT ID, NAME, VERSION, UPLOADED FROM Registry WHERE NAME = ?", name)
	if err != nil {
		return_500_packet(w, r)
		log.Print(err)
	}
	defer rows.Close()
	// get all versions of named package
	for rows.Next() {
		var timevar string
		var md models.PackageMetadata
		if err := rows.Scan(&md.ID, &md.Name, &md.Version, &timevar); err != nil {
			// package with name not found
			return_500_packet(w, r)
			//log.Fatal(err)
		}
		if err != nil {
			return_500_packet(w, r)
		}
		t, err := time.Parse("2006-01-02 15:04:05", timevar)
		if err != nil {
			return_500_packet(w, r)
		}
		timevar = t.Format(time.RFC3339)
		metadataList = append(metadataList, md)
		times = append(times, timevar)
	}

	// iterate through versions of package and get rest of history
	for i, md := range metadataList {
		var history models.PackageHistoryEntry
		history.User = &models.User{Name: "test", IsAdmin: false}
		history.Date, err = time.Parse("2006-01-02T15:04:05Z", times[i])
		if err != nil {
			return_500_packet(w, r)
		}
		history.PackageMetadata = &md
		history.Action = "Uploaded"
		ret = append(ret, history)
	}

	json.NewEncoder(w).Encode(ret)
	w.WriteHeader(200)
}

func handle_package_byregex(w http.ResponseWriter, r *http.Request) {

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/packages", verifyJWT(handle_packages))
	router.HandleFunc("/package/{id}", verifyJWT(handle_package_id))
	router.HandleFunc("/package/{id}/rate", verifyJWT(handle_package_rate))
	router.HandleFunc("/package/byName/{name}", verifyJWT(handle_package_byname))
	router.HandleFunc("/package/byRegEx", verifyJWT(handle_package_byregex))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}
