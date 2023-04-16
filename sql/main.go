package main

/*
Provide GET functionality for the following endpoints:
/packages
/package/(id)
/package/(id)/rate
/package/byName/(name)
*/

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

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

func return_error_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal error"))
}

// Get the packages from the registry
func handle_packages(w http.ResponseWriter, r *http.Request) {
	//db, err := connect()
	db, err := connect_test_db()
	if err != nil {
		log.Fatal(err)
	}

	//logger.Info(fmt.Sprintf("Received %s request", r.Method))
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}
	//logger.Info(headers)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		//logger.Info("\nError reading body of request\n")
		return_error_packet(w, r)
		return
	}
	print(body)
	//logger.Info(fmt.Sprintf("Body:\n%s\n", body))
	// TODO: ADD PAGINATION STUFF
	res, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry;`)
	if err != nil {
		log.Fatal(err)
	}
	print("res is: ", res, "\n")
	defer res.Close()

	w.WriteHeader(200)
	// --------- DEBUGGING/EXPERIMENTAL CODE TO VIEW RETURN ---------
	for res.Next() {
		//check for correct indexing
		var pack models.PackageMetadata // I have no idea what to put here
		err := res.Scan(&pack.ID, &pack.Version, &pack.Name)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", pack)
	}
	// --------------------------------------------------------------
}

// Return this package (ID)
func handle_packages_id(w http.ResponseWriter, r *http.Request) {
	//db, err := connect()
	db, err := connect_test_db()
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	// logger.Info(fmt.Sprintf("Received %s request", r.Method))
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}
	// logger.Info(headers)

	id := r.Header.Get("Id")
	if id == "" {
		// logger.Info("\nNo Matching Value to key Id\n")
		return_error_packet(w, r)
		return
	}

	res, err := db.Query("SELECT ID, NAME, VERSION FROM Registry WHERE ID ==" + id + ";")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()
	var meta models.PackageMetadata
	err = res.Scan(&meta.ID, &meta.Name, &meta.Version)
	if err != nil {
		log.Fatal(err)
	}
	res, err = db.Query(`BEGIN SELECT 
							B.BINARY_FILE, 
							A.URL 
							B.JS_PROGRAM
							FROM Registry AS A
							WHERE A.ID == ` + id + `
							INNER JOIN Binaries AS B
								ON A.BINARY_PIK == B.ID
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
func handle_packages_rate(w http.ResponseWriter, r *http.Request) {
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
	id := r.Header.Get("Id")
	if id == "" {
		// logger.Info("\nNo Matching Value to key Id\n")
		return_error_packet(w, r)

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
func handle_packages_byname(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/packages", handle_packages)
	http.HandleFunc("/packages/id", handle_packages_id)
	http.HandleFunc("/packages/id/rate", handle_packages_rate)
	http.HandleFunc("/packages/byName/name", handle_packages_byname)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
