package dbService

import (
	"fmt"
	"database/sql"
	"time"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type url struct {
	shortURL string
	longURL string
	created string
	expiration string
	totalCalls int
}

var db *sql.DB

func InitDB() {
	var err error
    db, err = sql.Open("sqlite3", "database/URLShortener.db")
    if err != nil {
        log.Fatal(err)
    }
    
   // defer db.Close()

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    
    db.Exec("create table if not exists url (shortURL text, longURL text, created text, expiration string, total_calls int)")

    fmt.Println(db);
    fmt.Println("Connected to database successful")
}

//	fmt.Println(db.Ping());
//	_, err := db.Exec("create table if not exists url (shortURL text, longURL text, created text, expiration string, total_calls int)")
//	if (err != nil) {
//		return 1, fmt.Errorf("init-table-failed %v")
//	}

//	return 0, nil
//}

func GetLongURL(shortURL string) (string, error) {
	var longURL, expiration string

    if err := db.QueryRow("select longURL, expiration from url where shortURL = ?", shortURL).Scan(&longURL, &expiration); err != nil {
		longURL = ""
		return longURL, fmt.Errorf("empty row %v")
	}

    log.Println("longURL from db", longURL)
    log.Println("expiration from db:", expiration)

    if(expiration != "") {
    const dateform = "2006-01-02"
    exp, _ := time.Parse(dateform, expiration)
    
    log.Println("expiration date:", exp)

	if(time.Now().After(exp)) {
		DeleteShortURL(shortURL)
		longURL = ""
	 }
    }
    log.Println("longURL to send:", longURL)

	UpdateTotalCalls(shortURL)

    return longURL, nil    
}                                                  

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
	
	res, err := result.RowsAffected()

	if(err != nil) {
		return 0, fmt.Errorf("no-rows-affected %v")
	}

	return res, nil
}

func DeleteShortURL(shortURL string) (int64, error) {
	result, err := db.Exec("delete from url where shortURL = ?", shortURL)
	if err != nil {
	 return 0, fmt.Errorf("update-failed %v")
	}
	
	res, err := result.RowsAffected()

	if(err != nil) {
		return 0, fmt.Errorf("no-rows-affected %v")
	}

	return res, nil
}

func getURLTotalCalls(shortURL string) (int, error) {
	 var totalCalls int
	if err := db.QueryRow("select total_calls from url where shortURL = ?", shortURL).Scan(&totalCalls); err != nil {
		return 0, fmt.Errorf("failed-row-read %v")
	}

    return totalCalls, nil
}

func UpdateTotalCalls(shortURL string) (int64, error) {
	_, err := getURLTotalCalls(shortURL)
	if err != nil {
		return 0, fmt.Errorf("failed-calls-update %v")
	}
	result, err := db.Exec("update url set total_calls = ? where shortURL = ?", shortURL)

	if err != nil {
	 return 0, fmt.Errorf("update-failed %v")
	}
	
	res, err := result.RowsAffected()

	if(err != nil) {
		return 0, fmt.Errorf("no-rows-updated %v")
	}

	return res, nil
}

func RandomKeyExistsQ(shortURL string) (bool, error) {
	log.Println(db)
	log.Println(shortURL)
	var exists bool
	if err := db.QueryRow("select exists(select shortURL from url where shortURL =" + shortURL +" )").Scan(&exists); err != nil {
 		return false, fmt.Errorf("failed-read %v")
	}

    return exists, nil
}


