package dbService

import (
	"fmt"
	"database/sql"
	"time"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// Struct type replicating url table
type url struct {
	shortURL string
	longURL string
	created string
	expiration string
	totalCalls int
}

var db *sql.DB

// Opens a connection to database
// Pings the database to check if it is in usable state
// Creates tables url
func InitDB() {
	var err error
    db, err = sql.Open("sqlite3", "database/URLShortener.db")
    if err != nil {
    	log.Println("Failed to connect to database/URLShortener.db")
        log.Fatal(err)
    }
    
   // defer db.Close()

    pingErr := db.Ping()
    if pingErr != nil {
    	log.Println("Connected to database but Ping returned no response")
        log.Fatal(pingErr)
    }
    
    _, err = db.Exec("create table if not exists url (shortURL text, longURL text, created text, expiration string, total_calls int)")

    if err != nil {
    	log.Println("Failed to initialize tables")
    }

    log.Println("Connected to database/URLShortener.db and tables initialized successfully")
    fmt.Println("Connected to database successful")
}

// Gets the long URL associated with a short URL
// Returns the long URL associated with a short URL
func GetLongURL(shortURL string) (string, error) {
	var longURL, expiration string

    if err := db.QueryRow("select longURL, expiration from url where shortURL = ?", shortURL).Scan(&longURL, &expiration); err != nil {
		longURL = ""
		return longURL, fmt.Errorf("empty row %v")
	}

    log.Println("longURL from db", longURL)
    log.Println("expiration from db:", expiration)

    // If the URL has an expiration date, check if it is expired
    // If the URL is expired delete the entry for that URL
    if(expiration != "") {   
		if(URLExpired(shortURL, expiration)) {
			DeleteShortURL(shortURL)
			longURL = ""
	 	}
    }

    log.Println("longURL to send:", longURL)

    // If the URL is not expired, update the number of total calls
    if longURL != "" {
	callsUpdated, err := UpdateTotalCalls(shortURL)
	if(err != nil) {
		log.Println("Failed to update total calls for url:", shortURL)
	}
    
    log.Println("Total calls for shortURL" ,shortURL, "updated to ", callsUpdated)

    }
    

    return longURL, nil    
}

// Check if a short URL is expired
func URLExpired(shortURL string, expiration string) bool {
	const dateform = "2006-01-02"
    exp, _ := time.Parse(dateform, expiration)
    return time.Now().After(exp)
}

// Saves the short URL, long URL, date created, expiration and total calls in database                                                  
func SaveShortURL(shortURL string, longURL string, expiration string) (int64, error) {
	log.Println(db)
	log.Println(shortURL)
	log.Println(longURL)
	log.Println(expiration)
	result, err := db.Exec("insert into url (shortURL, longURL, created, expiration, total_calls) values (?, ?, ?, ?, ?)", shortURL, longURL, time.Now(), expiration, 0)
	if err != nil {
	log.Println("insert-failed", err)
	 return 0, fmt.Errorf("insert-failed %v")
	}
	
	res, _ := result.RowsAffected()

	return res, nil
}

// Deletes a short URL once it is expired
func DeleteShortURL(shortURL string) (int64, error) {
	result, err := db.Exec("delete from url where shortURL = ?", shortURL)
	if err != nil {
	 return 0, fmt.Errorf("update-failed %v")
	}
	
	res, _ := result.RowsAffected()

	return res, nil
}

// Gets the total number of times a URL has been called since its creation
func getURLTotalCalls(shortURL string) (int, error) {
	 var totalCalls int
	if err := db.QueryRow("select total_calls from url where shortURL = ?", shortURL).Scan(&totalCalls); err != nil {
		return 0, fmt.Errorf("failed-row-read %v")
	}

    return totalCalls, nil
}

// Updates the total number of times a URL has been called
func UpdateTotalCalls(shortURL string) (int64, error) {
	calls, _ := getURLTotalCalls(shortURL)
	result, err := db.Exec("update url set total_calls = ? where shortURL = ?", calls+1, shortURL)

	if err != nil {
	 return 0, fmt.Errorf("update-failed %v")
	}
	
	res, _ := result.RowsAffected()

	return res, nil
}

// Checks if a short URL exists in the database
func RandomKeyExistsQ(shortURL string) (bool, error) {
	log.Println(db)
	log.Println(shortURL)
	var exists bool
	if err := db.QueryRow("select exists(select shortURL from url where shortURL =" + shortURL +" )").Scan(&exists); err != nil {
 		return false, fmt.Errorf("failed-read %v")
	}

    return exists, nil
}


