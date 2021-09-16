# URL Shortener

URL Shortener is an application designed to provide a short and descriptive URL for a given long URL.

## Design and implementation

This application is primarily designed using Golang, uses SQLite for persisting URLs and related analytics. The user interface is a simple form with fields:

- Long URL (required)
- Alias (optional)
- Expiration (optional)

If no `Alias` is provided a random alphanumeric string is generated for the short URL.
If no `Expiration` is provided the URL lives forever.

The application provides two ways of creating and accessing the short URLs.

 - Form
 - REST API

### Implementation

The application is implemented broadly in four parts:

- User interface - takes user input through a form
- HTTP server - handles HTTP requests for URL creation, getting long URL and redirecting and analytics
- Service layer - handles any processing required before sending off the data to database, e.g: checking if URL already exists
- Persistent layer - handles all the database related queries

## Build and deploy URL Shortener

### Build from GitHub

- Check out the project from: https://github.com/mvish/url-shortener

- If the current directory is not url-shortener, cd to that directory

	```> cd url-shortener```

- Run the following command to build and run the application

	```go build && go run .```

### Build from container

To build and deploy using container, make sure [Docker Desktop](https://www.docker.com/products/docker-desktop) installed and running.

- Change to the url-shortener directory
	
	```cd url-shortener```

- Build the application:

	```docker build -t url-shortener .```

- Run the application:

	```docker run -p 8080:8080 -it url-shortener```

## Run URL Shortener

- Go to `http://localhost:8080` to launch the application home page

- The home page contains:
  
  - A brief decription of the application
  - A link to the form that creates short URLs

- To create a short URL:

  - Click on the button "Create short URL"
  - In the form provide a long URL, optionally an alias and an expiration
  - Click on `Create Short URL`
  - The created short URL should appear at the bottom, click on it to re-direct to original URL

## URL shortener API

An API to create a short URL, redirect short URL to original URL and delete short URLs.

## URL[/api/v1/url/]

### GET

Gets and redirects to original URL given a short URL.

+ Parameters

    - shorturl(required, string) - short URL to be visited

+ Request `GET /api/v1/url/{newurl}`

+ Response 200 OK (application/json)

    ```json
	   {
          "longURL": "https://medium.com/wesionary-team/map-types-in-golang-24591abbafc6",
          "shortURL": "newurl"
	   }
	```

+ Response 404 (application/json)

  If no short URL is provided:

    ```{"errorCode": "short-url-missing"}```

  If no long URL is found:

    ```{"errorCode": "url-not-found"}```


### POST

Creates a unique random or a custom short URL for a long given URL.

Fields `longURL`, `alias` and `expiration` are used to create the short URL. Following are the default values used for these fields:

- `originalURL` - no default
- `customName` - a randomly generated string
- `expiration` - 1 year from the date of creation

+ Request `POST /api/v1/url` (application/json)

	 ```json
	 {
	    "longURL": "http://www.mywebsite.com/share?us=2HujKemgLsuI13",
	    "alias": "my-web",
	    "expiration": "10-15-2021"
	 }
	 ```

+ Response 201 (application/json)

	```json
	{
		"shortURL": "http://myweb/my-web",
		"longURL": "http://www.mywebsite.com/share?us=2HujKemgLsuI13"
	}
	```

+ Response 400 (application/json)
	
	If the request body is not valid:

    ```{"errorCode": "malformed-json"}```

    If the long URL is missing:

    ```{"errorCode": "missing-long-url"}```

### DELETE

Deletes a short URL.

+ Request `DELETE /api/v1/url/{shortURL}`

+ Response 200

+ Response 404 (application/json)

URL to be deleted does not exist:

```json
{"errorCode": "missing-short-url"}
```


## URL shortener analytics API

An API to get top URLs, total number of times a URL is called in past "n" hours and total number of times a URL is called in past "n" days.

## URL[/api/v1/analytics/top/{limit}]

### GET

Gets the top "n" called short URLs. If no `limit` is provided, top five URLs are returned.

+ Request `GET /api/v1/analytics/top/3`

+ Response 200

    ```json
    [
      {
        "Url": "newurl",
        "TotalCalls": 5
      },
      {
        "Url": "id7sydksj",
        "TotalCalls": 2
      },
      {
        "Url": "b79d4e33",
        "TotalCalls": 1
      }
    ]
    ```
+ Response 500

If no rows exists or a database failure occurs:

    ```{errorCode: "failure:get-top-five-row"}```

## URL[/api/v1/analytics/{shortURL}]

### GET

Gets the number of times a short URL has been visited in the past "n hours" or "n days".

+ Parameters
    
    - shortURL(required, string) - short URL for which information needs to be returned
    - hours(optional, string) - represents number of hours in the past from current hour
    - days(optional, string) - represents number of hours in the past from current day

 If no paramters are provided, the total count of calls from past 1 day is returned as the default.

+ Request `GET /api/v1/analytics/{shortURL}?hours=5`, `GET /api/v1/analytics/{shortURL}?days=2`

+ Response 200 (application/json)

    ```json
        {
           "newurl": 11
        }
    ```

+ Response 400 (application/json)

If no parameters are provided:

    ```{"errorCode": "no-params-found"}```

If no short URL is provided:

    ```{"errorCode": "shortURL-not-found"}```


## References:

- HTTP server in go: https://gowebexamples.com/http-server/
