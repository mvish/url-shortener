package service

import (
	  "log"
	  "errors"
	  "strings"
	  "net/url"
	  "regexp"
	  "time"
	  "github.com/google/uuid"
      "mili.photos/dbService"
)

// Processes the short URL requested (in case an alias is provided) or generate a new random URL
// Saves generated short URL to database
func CreateAndSaveShortURL(longURL string, alias string, expiration string) (string, error, int) {
	var shortURL string

	// Check if the provided URL is a valid URL
	validURL := validateLongURL(url.QueryEscape(longURL))
	log.Println("shortenURLService: Long URL is valid")

	if validURL {

		// If user provided no alias, generate a random URL of length 8
		if alias == "" {
		shortURL = generateNewRandomCode();
		log.Println("shortURL random generated", shortURL)
        
        // Check if the short URL already exist
        shortURLExists, err := shortURLExistsQ(shortURL)

        if err != nil {
        	log.Println("Error checking if short URL already exists")
        	return shortURL, errors.New("failure:short-url-existence-unknown"), 500
        }

        // If the short URL already exists keep generting until one which does not exist is found and set that as the acceptable short url
        for shortURLExists == true {
        	shortURL = generateNewRandomCode();
        	shortURLExists, _ = shortURLExistsQ(shortURL)
        	log.Println("shortURLExist2: ", shortURLExists)
        }

	    } else {
	   		// If the user provided an alias check if the alias is already taken
	   		// Also check if the alias is 8 characters or less and contains only alphanumeric characters
	   		aliasExists, _ := shortURLExistsQ(alias)
	   		aliasValid := validateAlias(alias)
	   		if !aliasExists && aliasValid {
          		shortURL = alias
            } else {
            	log.Println("Alias provided either exists or contains invalid characters")
          		return alias, errors.New("failure:invalid-alias"), 409
          	}
	    }

	}

    // Save the accepted short URL in database	
   log.Println("shortenURLService: Accepted shortURL: ", shortURL)
   res, err := dbService.SaveShortURL(shortURL, longURL, expiration)
   if err != nil || res == 0 {
   		log.Println("shortenURLService: Failed to save url created: ", shortURL)
   	 	return "", errors.New("failure:failed-saving-url-created"), 500
   }
   
   // Update url_count_hourly table to have a row for the newly created short URL
   timeStamp := getTimeStamp()
   _, err = dbService.UpdateHourlyCalls(timeStamp, shortURL)
   if err != nil {
    	log.Println("shortenURLService: Failed to update hourly calls for url: ", shortURL)
   }

   return shortURL, nil, 201
}

// Gets the long URL associated with a short URL
func GetLongURL(shortURL string) (string, bool, error) {
    url, deleted, err := dbService.GetLongURL(shortURL)
    if err != nil {
    	log.Println("shortenURLService: Failed to get long URL for: ", shortURL)
    	return url, deleted, errors.New("failure:failed-get-long-url")
    }

    // If the URL is not expired, update the number of total calls and hourly calls
    if deleted == false {
    	// Update hourly calls for the URL
    	timeStamp := getTimeStamp()
    	_, err := dbService.UpdateHourlyCalls(timeStamp, shortURL)
    	if err != nil {
    		log.Println("shortenURLService: Failed to update hourly calls for url: ", shortURL)
    	}

    	// Update total calls for the URL
    	totalCallsUpdated, err := dbService.UpdateTotalCalls(shortURL)
		if(err != nil) {
			log.Println("dbService: Failed to update total calls for url:", shortURL)
		}
    
    	log.Println("dbService: Total calls for shortURL" ,shortURL, "updated by ", totalCallsUpdated)

	}

    return url, deleted, nil
}

// Checks if a given URL is valid or not
func validateLongURL(longURL string) bool {
	_, err := url.Parse(longURL)
	return err == nil
}

// Checks if a given alias is valid or not
func validateAlias(alias string) bool {
	exp := regexp.MustCompile("^[a-zA-Z0-9]*$")
	if len(alias) <= 50 {
		return exp.MatchString(alias)
	}

	return false
}

// Generates a random string to be used as short URL
func generateNewRandomCode() string {
	uuidNew := uuid.New()
	uuidNoHyphen := strings.Replace(uuidNew.String(), "-", "", -1)
	return uuidNoHyphen[0:len(uuidNoHyphen)-24]
}

// Checks if a given short URL already exists
func shortURLExistsQ(shortURL string) (bool, error) {
	log.Println("In shortURLExistsQ", shortURL)
	keyExistsQ, err := dbService.URLExistsQ(shortURL)
	if err != nil {
		log.Println("shortenURLService: Failed to check if url ", shortURL, "exists or not")
		return false, errors.New("failure:check-short-url-exists")
	}

	return keyExistsQ, nil
}

// Gets the current time stamp in the form yyyy-mm-dd HH:00:00
func getTimeStamp() string {
	currentTime := time.Now()
    layout := "2006-01-02 15:04:05"
    roundedTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), 0, 0, 0, currentTime.Location()).Format(layout)

    return roundedTime;
}
