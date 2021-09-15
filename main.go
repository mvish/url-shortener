package main

import (
    "os"
    "log"
    "strings"
    "net/http"
    "net/url"
    "encoding/json"
    "mili.photos/service"
    "mili.photos/dbService"
    lru "github.com/flyaways/golang-lru"
)


type urlShort struct {
    LongURL string `json:"longURL"`
    Alias  string `json:"alias"`
    Expiration  string `json:"expiration"`
}

var cache lru.LRUCache
const apiURL = "/api/v1/url/"

func main() {
        
    file, err := os.OpenFile("url-shortener-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)

    dbService.InitDB()

    cache, _ = lru.New(200)
 
    urlExpandHandler := http.HandlerFunc(urlExpandAndRedirectOperation)
    http.Handle("/u/", urlExpandHandler)

    urlShortenHandler := http.HandlerFunc(createURLOperation)
    http.Handle("/url", urlShortenHandler)

    urlShortenAPIHandler := http.HandlerFunc(urlShortenAPIOperations)
    http.Handle(apiURL, urlShortenAPIHandler)

    urlAnalyticsHandler := http.HandlerFunc(urlAnalyticsOperations)
    http.Handle("/api/v1/analytics/", urlAnalyticsHandler)

    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/", fs)

    http.ListenAndServe(":8080", nil)
}

// Handler for getting long URL associated with a short one and redirecting
func urlExpandAndRedirectOperation(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[2:][1:]
    var longURLInterface interface{}
    var longURL string
    var err error
    var deleted bool

    // Check if the requested URL is in cache
    if cache.Contains(path) {
        longURLInterface, _ = cache.Get(path)
        longURL = longURLInterface.(string)
    } else {
        // If the requested URL is not cached then request it from database
        longURL, deleted, err = service.GetLongURL(path)
        if err != nil {
            setErrorResponse(w,"internal-error", err.Error(), http.StatusInternalServerError)
            return
        }
        if deleted {
            cache.Remove(path)
        } else {
            cache.Add(path, longURL)
        }
    }

    if longURL != "" { 
        http.Redirect(w, r, longURL, http.StatusFound)
    } else {
        setErrorResponse(w, "url-not-found", "Requested URL does not exist or is expired", http.StatusNotFound)
    }
}

// Handler for creating a new short URL given a long URL
func createURLOperation(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    if strings.TrimSpace(r.FormValue("url_long")) == "" {
        http.Redirect(w, r, "/index.html?error="+url.QueryEscape("no-long-url-provided") , http.StatusTemporaryRedirect)
    }

     shortURL, errService, _ := service.CreateAndSaveShortURL(r.FormValue("url_long"), r.FormValue("url_alias"), r.FormValue("url_exp"))
    if errService != nil {
        log.Println("error creating URL")
        http.Redirect(w, r, "/index.html?error="+url.QueryEscape(errService.Error()) , http.StatusTemporaryRedirect)
    }

    http.Redirect(w, r, "/index.html?created="+shortURL , http.StatusFound)
}

// Handler for REST APIs
// GET /v1/url
// POST /v1/url
func urlShortenAPIOperations(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        shortURL := strings.TrimPrefix(r.URL.Path, apiURL)

        if len(shortURL) < 1 {
            log.Println("Url param 'shortURL' is missing")
            setErrorResponse(w,"short-url-missing", "No short URL provided", http.)
            return
        }

        longURL, _, err := service.GetLongURL(shortURL)
        if err != nil {
            setErrorResponse(w, "not-found", err.Error(), http.StatusNotFound)
            return
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
        return
        }

        var urlInfo urlShort

        decoder := json.NewDecoder(r.Body)
        decoder.DisallowUnknownFields()
        err := decoder.Decode(&urlInfo)

        if err != nil {
           log.Println("urlShortenAPIOperations: malformed json in request body: ", err)    
           setErrorResponse(w, "malformed-json", "Error parsing request." + err.Error(), http.StatusBadRequest)
           return
        }

        log.Println("creating shortURL")

        // Check if a long URL is provided, this field is required
        if strings.TrimSpace(urlInfo.LongURL) == "" {
            setErrorResponse(w, "missing-long-url", "Long URL is required to create a short URL", http.StatusBadRequest)
            return
        }
                
        shortURL, errService, status := service.CreateAndSaveShortURL(urlInfo.LongURL, urlInfo.Alias, urlInfo.Expiration)
        if errService != nil {
            log.Println("urlShortenAPIOperations: Failed to create short URL for url: ", urlInfo.LongURL, " error: ", errService)        
            setErrorResponse(w, errService.Error(), "Failed to create short url", status)
            return
        }

        response := make(map[string]string)
        response["shortURL"] = shortURL
        response["longURL"] = urlInfo.LongURL
        
        setResponse(w, response, 201)

    case "DELETE":
        shortURL := strings.TrimPrefix(r.URL.Path, apiURL)

        if len(shortURL) < 1 {
            log.Println("Url param 'shortURL' is missing")
            return
        }
           
        rowsAffected, err := dbService.DeleteShortURL(shortURL)
        if err != nil || rowsAffected < 1 {
            log.Println("Failed to delete short URL:", shortURL, "error: ", err)
            setErrorResponse(w, err.Error(), "An internal error occurred while deleting short URL", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}

// Handler for analytics requests
// GET /analytics
func urlAnalyticsOperations(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    paramURL := query["shortURL"]
    paramHour := query["hours"]
    paramDay := query["days"]
    paramTop := query["top"]

    // If no parameters are specified, provide the usage
    if(len(paramURL) == 0 && len(paramHour) == 0 && len(paramDay) == 0 && len(paramTop) == 0) {
        setErrorResponse(w, "no-params-found", "Usage: /analytics?top=n, /analytics?shortURL=url&hours=n, /analytics?shortURL=url&days=n", http.StatusBadRequest)
        return
    }

    // top n URLs visited
    if(len(paramTop) != 0) {
        urlCount, err := dbService.Top(paramTop[0])

        if err != nil {
            log.Println("urlAnalyticsOperations: Failed to get top five most visited URLs")
            setErrorResponse(w, err.Error(), "Failed to get top " + paramTop[0] + "visited URLs", http.StatusInternalServerError)
        }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    jsonResponse, _ := json.Marshal(urlCount)
    w.Write(jsonResponse)
    return
    } 

    if len(paramURL) < 1 {
        log.Println("Url param 'shortURL' is missing")
        setErrorResponse(w, "shortURL-not-found", "short URL is required", http.StatusBadRequest)
        return
    }

    // get URL count by hours
    if len(paramHour) != 0 {
        urlCountHour, err := dbService.GetURLCountPastnHours(paramURL[0], paramHour[0])
        if err != nil {
            setErrorResponse(w, err.Error(), "Failed to URL calls by hour", http.StatusInternalServerError)
            return
        }

        response := make(map[string]int)
        response[paramURL[0]] = urlCountHour

        setResponseForAnalytics(w, response, http.StatusOK)
        return
    }

    // get URL count by days
    if len(paramDay) != 0 {
        urlCountDay, err := dbService.GetURLCountPastnDays(paramURL[0], paramDay[0])
        if err != nil {
            setErrorResponse(w, err.Error(), "Failed to URL calls by hour", http.StatusInternalServerError)
            return
        }

        response := make(map[string]int)
        response[paramURL[0]] = urlCountDay

        setResponseForAnalytics(w, response, http.StatusOK)
        return
    } 

}

// Sets the response JSON
func setResponse(w http.ResponseWriter, response map[string]string, httpStatusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(httpStatusCode)
    jsonResponse, _ := json.Marshal(response)
    w.Write(jsonResponse)
}

// Sets response JSON for analytics related requests
func setResponseForAnalytics(w http.ResponseWriter, response map[string]int, httpStatusCode int) {
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



