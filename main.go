package main

import (
	"os"
	"log"
	"strings"
	"net/http"
	"encoding/json"
    "mili.photos/service"
    "mili.photos/dbService"
)


type urlShort struct {
	LongURL string `json:"longURL"`
	Alias  string `json:"alias"`
	Expiration  string `json:"expiration"`
}

func main() {
        
    file, err := os.OpenFile("url-shortener-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)

	dbService.InitDB()

	urlExpandHandler := http.HandlerFunc(urlExpandAndRedirectOperation)
	http.Handle("/u/", urlExpandHandler)

	urlShortenHandler := http.HandlerFunc(createURLOperation)
	http.Handle("/url", urlShortenHandler)

	urlShortenAPIHandler := http.HandlerFunc(urlShortenAPIOperations)
	http.Handle("/v1/url", urlShortenAPIHandler)

	fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}

// Handler for getting long URL associated with a short one and redirecting
func urlExpandAndRedirectOperation(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[2:][1:]
	longURL, err := service.GetLongURL(path)
	 if err != nil {
		setErrorResponse(w,"internal-error", err.Error(), http.StatusInternalServerError)
		return;
	}

	if longURL != "" { 
        http.Redirect(w, r, longURL, http.StatusFound)
    } else {
		setErrorResponse(w,"not-found", "Requested URL does not exist or is expired", http.StatusOK)
	}
}

// Handler for creating a new short URL given a long URL
func createURLOperation(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if strings.TrimSpace(r.FormValue("url_long")) == "" {
        setErrorResponse(w, "missing-long-url", "Long URL is required to create a short URL", http.StatusBadRequest)
        return;
    }

 	shortURL, errService, status := service.CreateAndSaveShortURL(r.FormValue("url_long"), r.FormValue("url_alias"), r.FormValue("url_exp"))
    if errService != nil {
		setErrorResponse(w, "internal-error", errService.Error(), status)
		http.Redirect(w, r, "/static/index.html?error="+"unableToCreate" , http.StatusFound)
		return;
	}

	http.Redirect(w, r, "/static/index.html?created="+shortURL , http.StatusFound)

}

// Handler for REST APIs
// GET /v1/url
// POST /v1/url
func urlShortenAPIOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": 
		query := r.URL.Query()
		param := query["shortURL"]

		if len(param[0]) < 1 {
			log.Println("Url param 'shortURL' is missing")
			return
		}
           
		shortURL := param[0];
		longURL, err := service.GetLongURL(shortURL)
        if err != nil {
			setErrorResponse(w,"internal-error", err.Error(), http.StatusInternalServerError)
			return;
		}
        
        response := make(map[string]string)
		response["shortURL"] = shortURL
		response["longURL"] = longURL

		setResponse(w, response, 200)

        
	case "POST":
		log.Println("Creating short URL")
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			log.Println("urlShortenAPIOperations: content-type is not application/json")
	        setErrorResponse(w, "invalid-content-type", "Content type is not application/json", http.StatusUnsupportedMediaType)
		return;
		}

		var urlInfo urlShort

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&urlInfo)

		if err != nil {
		   log.Println("urlShortenAPIOperations: malformed json in request body: ", err)	
		   setErrorResponse(w, "malformed-json", "Error parsing request." + err.Error(), http.StatusBadRequest)
		   return;
		}

        log.Println("creating shortURL")

        // Check if a long URL is provided, this field is required
        if strings.TrimSpace(urlInfo.LongURL) == "" {
        	setErrorResponse(w, "missing-long-url", "Long URL is required to create a short URL", http.StatusBadRequest)
        	return;
        }

                
        shortURL, errService, status := service.CreateAndSaveShortURL(urlInfo.LongURL, urlInfo.Alias, urlInfo.Expiration)
        if errService != nil {
        	log.Println("urlShortenAPIOperations: Failed to create short URL for url: ", urlInfo.LongURL)		
			setErrorResponse(w, errService.Error(), "Failed to create short url", status)
			return;
	    }


		response := make(map[string]string)
		response["shortURL"] = shortURL
		response["longURL"] = urlInfo.LongURL
		
		setResponse(w, response, 201)
	}
}

// Sets the response JSON
func setResponse(w http.ResponseWriter, response map[string]string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

// Sets the response JSON for errors
func setErrorResponse(w http.ResponseWriter, errorCode string, errorMessage string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	response := make(map[string]string)
	response["errorCode"] = errorCode
	response["errorMessage"] = errorMessage
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}



