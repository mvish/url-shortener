package service

import (
	  "fmt"
	  "log"
	  "strings"
	  "net/url"
	  "github.com/google/uuid"
      "mili.photos/dbService"
)

func CreateAndSaveShortURL(longURL string, alias string, expiration string) (string, error) {
	var shortURL string
	validURL := validate(url.QueryEscape(longURL))
	log.Println("validURL", validURL)
	if validURL {
		if alias == "" {
		shortURL = generateNewRandomCode();
		log.Println("shortURL random generated", shortURL)
        
        shortURLExist := shortURLExistsQ(shortURL)

        log.Println("shortURLExist", shortURLExist)

        for shortURLExist == true {
        	shortURL = generateNewRandomCode();
        	shortURLExist = shortURLExistsQ(shortURL)
        }

	   } else {
          shortURL = alias
	   }

	}
   log.Println("shortURL", shortURL)
   res, err := dbService.SaveShortURL(shortURL, longURL, expiration)
   if err != nil || res == 0 {
   	 return "", fmt.Errorf("failed-url-creation %v")
   }
   	 
   return shortURL, nil
}

func GetLongURL(shortURL string) (string, error) {
    url, err := dbService.GetLongURL(shortURL)
    if err != nil {
    	return url, fmt.Errorf("failed-get-long-url %v")
    }

    return url, nil
}


func validate(longURL string) bool {
	_, err := url.Parse(longURL)
	return err == nil
}

func generateNewRandomCode() string {
	uuidNew := uuid.New()
	uuidNoHyphen := strings.Replace(uuidNew.String(), "-", "", -1)
	return uuidNoHyphen[0:len(uuidNoHyphen)-17]
}

func shortURLExistsQ(key string) bool {
	keyExistsQ, _ := dbService.RandomKeyExistsQ(key)
	//if err != nil {
	//	return false, fmt.Errorf("failed-getting-key %v")
	//}

	return keyExistsQ
}


