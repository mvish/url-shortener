package dbService

import (
	"errors"
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

type URLCount struct {
	Url string
	TotalCalls int
}

type TopUrls []URLCount

var db *sql.DB

// Opens a connection to database
// Pings the database to check if it is in usable state
// Creates tables url
func InitDB() (int64, error) {
	var errCreateURLTable error
	var errCreateURLHourly error
	var err error
    db, err = sql.Open("sqlite3", "database/URLShortener.db")
    if err != nil {
    	log.Println("dbService: Failed to connect to database/URLShortener.db", err)
        return 1, errors.New("failure:db-connection")
    }
    
   // defer db.Close()

    pingErr := db.Ping()
    if pingErr != nil {
    	log.Println("dbService: Connected to database but Ping returned no response", pingErr)
        return 1, errors.New("failure:db-ping-no-response")
    }
    
    _, errCreateURLTable = db.Exec("create table if not exists url (short_url text NOT NULL, long_url text, created text, expiration text, total_calls, PRIMARY KEY (short_url))")
    _, errCreateURLHourly = db.Exec("create table if not exists url_count_hourly (start text NOT NULL, short_url text NOT NULL, count int, PRIMARY KEY (start, short_url))")

    if errCreateURLTable != nil || errCreateURLHourly != nil{
    	log.Println("dbService: Failed to initialize tables", errCreateURLTable, errCreateURLHourly)
    	return 1, errors.New("failure:tables-init")

    }

    log.Println("dbService: Connected to database/URLShortener.db and tables initialized successfully")
   
    return 0, nil
}

// Gets the long URL associated with a short URL
// Returns the long URL associated with a short URL
func GetLongURL(shortURL string) (string, bool, error) {
	var longURL, expiration string
	deleted := false

    if err := db.QueryRow("select long_url, expiration from url where short_url = ?", shortURL).Scan(&longURL, &expiration); err != nil {
		log.Println("dbService: Failed to get long URL for ", shortURL, " assigning longURL to empty string")
		return shortURL, false, errors.New("failure:get-long-url")
	}


    // If the URL has an expiration date, check if it is expired
    if(expiration != "") {   
		if(URLExpired(shortURL, expiration)) {
			// If the URL is expired delete the entry for that URL
			DeleteShortURL(shortURL)
			log.Println("dbService: Requested short URL: ", shortURL, " expired. Deleted URL, setting long URL to empty string")
			deleted = true
	 	}
    }
    
    return longURL, deleted, nil    
}

// Check if a short URL is expired
func URLExpired(shortURL string, expiration string) bool {
	const dateform = "2006-01-02"
    exp, _ := time.Parse(dateform, expiration)
    return time.Now().After(exp)
}

// Saves the short URL, long URL, date created, expiration and total calls in database                                                  
func SaveShortURL(shortURL string, longURL string, expiration string) (int64, error) {
	result, err := db.Exec("insert into url (short_url, long_url, created, expiration, total_calls) values (?, ?, ?, ?, ?)", shortURL, longURL, time.Now(), expiration, 0)
	if err != nil {
	 log.Println("dbService: Failed to insert row in table url: ", err)
	 return 0, errors.New("failure:db-insert-shortURL")
	}
	
	res, _ := result.RowsAffected()

	return res, nil
}

// Deletes a short URL once it is expired
func DeleteShortURL(shortURL string) (int64, error) {
	result, err := db.Exec("delete from url where short_url = ?", shortURL)
	if err != nil {
		log.Println("dbService: Failed to delete short URL: ", shortURL, err)
		return 0, errors.New("failure:db-delete-shortURL")
	}
	
	res, _ := result.RowsAffected()

	return res, nil
}

// Gets the total number of times a URL has been called since its creation
func getURLTotalCalls(shortURL string) (int, error) {
	var totalCalls int
	if err := db.QueryRow("select total_calls from url where short_url = ?", shortURL).Scan(&totalCalls); err != nil {
		log.Println("dbService: Failed to get total calls for URL: ", shortURL, err)
		return 0, errors.New("failure:get-total-calls")
	}

    return totalCalls, nil
}

// Updates the total number of times a URL has been called
func UpdateTotalCalls(shortURL string) (int64, error) {
	calls, _ := getURLTotalCalls(shortURL)
	result, err := db.Exec("update url set total_calls = ? where short_url = ?", calls+1, shortURL)

	if err != nil {
		log.Println("dbService: Failed to update total calls for URL: ", shortURL, err)
		return 0, errors.New("failure:update-total-calls")
	}
	
	res, _ := result.RowsAffected()

	return res, nil
}

// Checks if a short URL exists in the database
func URLExistsQ(shortURL string) (bool, error) {
	var exists int
	if err := db.QueryRow("select exists(select 1 from url where short_url = ?)", shortURL).Scan(&exists); err != nil {
		log.Println("dbService: Failed to check if URL exists: ", shortURL, err)
		return false, errors.New("failure:check-short-url-exists")
	}

    return exists == 1, nil
}

// Analytics handlers

// Gets the top n most called URLs
func Top(limit string) ([]URLCount, error) {
	//urlToCountMap := make(map[string]int)
	var urlCount TopUrls

	rows, err := db.Query("select short_url, total_calls from url order by total_calls desc limit ?", limit)
	if(err != nil) {
		log.Println("dbService: Failed to get top 5 URLs")
		return nil, errors.New("failure:get-top-five-urls")
	}

	defer rows.Close()

	for rows.Next() {
		var urlCountRow URLCount
		if err := rows.Scan(&urlCountRow.Url, &urlCountRow.TotalCalls); err != nil {
			log.Println("dbService: Failed to scan a row from top 5 URLs")
			return nil, errors.New("failure:get-top-five-row")
		}

		urlCount = append(urlCount, urlCountRow)
	}

	return urlCount, nil
}

// Updates the hourly total calls for a URL
func UpdateHourlyCalls(startTime string, shortURL string) (bool, error) {
	result, err := db.Exec("insert into url_count_hourly (start, short_url, count) values (?, ?, ?) on conflict (start, short_url) do update set count = count+1", startTime, shortURL, 1)
	if err != nil {
		log.Println("dbService: Updating hourly count failed for short url: ", shortURL, err)
		return false, errors.New("failure:db-update-hourly")
	}

	res, _ := result.RowsAffected()

	return res == 1, nil
}

// Gets the total calls for a URL in the past n hours
func GetURLCountPastnHours(shortURL string, n string) (int, error) {
	var sumCount int
	if err := db.QueryRow("select sum(count) from url_count_hourly where short_url = ? and datetime('now') >= datetime(start, ?)", shortURL, "-" + n + " Hour").Scan(&sumCount); err != nil {
		log.Println("dbService: Getting count for past ", n, "hours failed for short url:", shortURL, sumCount, err)
		return 0, errors.New("failure:url-count-by-hours")
	}

	return sumCount, nil
}

// Gets the total calls for a URL in the past n days
func GetURLCountPastnDays(shortURL string, n string) (int, error) {
	var sumCount int
	if err := db.QueryRow("select sum(count) from url_count_hourly where short_url = ? and date() >= date(start, ?)", shortURL, "-" + n + " Day").Scan(&sumCount); err != nil {
		log.Println("dbService: Getting count for past ", n, "days failed for short url: ", shortURL, err)
		return 0, errors.New("failure:url-count-by-days")
	}

	return sumCount, nil
}


