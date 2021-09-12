package main

import (
	"os"
	"log"
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
	http.Handle("/", urlExpandHandler)

	urlShortenHandler := http.HandlerFunc(urlOperations)
	http.Handle("/url", urlShortenHandler)

	fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}

func urlExpandAndRedirectOperation(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Println("Path:", path)
	longURL, err := service.GetLongURL(path)
	 if(err != nil) {
		setErrorResponse(w,"internal-error", err.Error(), http.StatusInternalServerError)
	}

	 if(longURL != "") { 
        http.Redirect(w, r, longURL, http.StatusMovedPermanently)
    } else {
		setErrorResponse(w,"not-found", "Requested URL does not exist or is expired", http.StatusOK)
	}

}

func urlOperations(w http.ResponseWriter, r *http.Request) {
	if(r.URL.Path != "/url") {
		http.NotFound(w,r)
		return
	}

	switch r.Method {
	case "GET": 
		query := r.URL.Query()
		param := query["shortURL"]

		if len(param[0]) < 1 {
			log.Println("Url param 'shortURL' is missing")
			return
		}
           
		shortURL := param[0];
		log.Println("shortURL to get", shortURL)
		// pass this shortURL to db to get the original URL and redirect to it
		longURL, err := service.GetLongURL(shortURL)
        if(err != nil) {
		setErrorResponse(w,"internal-error", err.Error(), http.StatusInternalServerError)
		}
        
        log.Println("longURL from get", longURL)

        if(longURL != "") { 
        http.Redirect(w, r, longURL, http.StatusMovedPermanently)
        } else {
		setErrorResponse(w,"not-found", "Requested URL does not exist or is expired", http.StatusOK)
		}
        
	case "POST":
 		r.ParseForm()

 		 url, errService := service.CreateAndSaveShortURL(r.FormValue("url_long"), r.FormValue("url_alias"), r.FormValue("url_exp"))
        if(errService != nil) {
			setErrorResponse(w, "internal-error", errService.Error(), http.StatusInternalServerError)
	    }

	    w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		response := make(map[string]string)
		response["shortURL"] = url
		response["longURL"] = r.FormValue("url_long")
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)

		/*log.Println("creating shortURL")
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
	        setErrorResponse(w, "invalid-content-type", "Content type is not application/json", http.StatusUnsupportedMediaType)
		return;
		}
		log.Println("checked content type")

		var urlInfo urlShort

		decoder := json.NewDecoder(r.Body)
		//decoder.DisallowUnknownFields()
		err := decoder.Decode(&urlInfo)

		log.Println("decoded url info")
		log.Println(urlInfo.LongURL)

		if err != nil {
		   setErrorResponse(w, "malformed-json", "Error parsing request." + err.Error(), http.StatusBadRequest)
		}

        log.Println("creating shortURL")
                
        // to do: url here needs to be sent as response
        url, errService := service.CreateAndSaveShortURL(urlInfo.LongURL, urlInfo.Alias, urlInfo.Expiration)
            if(errService != nil) {
			setErrorResponse(w, "internal-error", err.Error(), http.StatusInternalServerError)
	    }

	    if errService != nil {
	    	log.Println(errService)
	    }

		log.Println(url)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		response := make(map[string]string)
		response["shortURL"] = url
		response["longURL"] = urlInfo.LongURL
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)*/
				
				//setResponse(w, url, urlInfo.LongURL, http.StatusCreated)

            	// fmt.Fprintln(w,"%s: %s\n", "LongURL", urlInfo.LongURL)
            	// fmt.Fprintln(w,"%s: %s\n", "Alias", urlInfo.Alias)
            	// fmt.Fprintln(w,"%s: %s\n", "Expiration", urlInfo.Expiration

	}
}

func setErrorResponse(w http.ResponseWriter, errorCode string, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	response := make(map[string]string)
	response["errorCode"] = errorCode
	response["message"] = message
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

