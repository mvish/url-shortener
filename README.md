## URL Shortener

The URL Shortener is an application designed to provide a short and descriptive URL for a given long URL.

## Design and implementation

This application is primarily designed using Golang, uses SQLite for persisting URLs and related analytics. The user interface is a simple form with fields:

- Long URL (required)
- Alias (optional)
- Expiration (optional)

If no `Alias` is provided a random alphanumeric string is generated for the short URL.
If no `Expiration` is provided the URL lives forever.

The application provides two ways of creating and accessing the short URLs.

 - HTML form
 - REST API

 Both methods create a short URL of the form `http://localhost:8080/u/short-url`

## Build and deploy URL Shortener

### Build from GitHub

- Check out the project from: https://github.com/mvish/url-shortener

- If the current directory is not url-shortener, cd to that directory

	```> cd url-shortener```

- Run the following command to build and run the application

	```go build && go run .```

### Build from container

To build and deploy using container, make sure [Docker Desktop](https://www.docker.com/products/docker-desktop) installed.

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

### General notes and definitions

URL[/url]

GET

Gets and redirects to original URL given a short URL.

	* Parameters

		- shorturl(string) - short URL to be visited

	* Request `GET /url?shorturl=http://myweb/3sTjKlDo`

	* Response 200 OK (application/json)

	{
	  "originalURL": "http://www.mywebsite.com/share?us=2HujKemgLsuI13"
	}
	
	* Response 400 (application/json)

	{"errorCode": "invalid-url"}

	* Response 404 (application/json)

	If no URL exists.

	{"errorCode": "url-not-found"}


POST

Creates a unique random or a custom short URL for a given URL.

Fields `originalURL`, `customName` and `expiration` are used to create the short URL. Following are the default values used for these fields:

	- `originalURL` - no default
	- `customName` - a randomly generated string
	- `expiration` - 1 year from the date of creation

* Request `POST /url` (application/json)

	 {

	    "originalURL": "http://www.mywebsite.com/share?us=2HujKemgLsuI13",
	    "custonName": "my-web",
	    "expiration": "10-15-2021"

	 }

	

* Response 200 (application/json)

	{
		"shortURL": "http://myweb/my-web"
	}	


* Response 400 (application/json)

       {"errorCode": "malformed-json"}

DELETE

Deletes a short URL.

	* Parameters
		- shortURL - URL to be deleted

* Request `DELETE /url?shortURL=http://myweb/my-web`

* Response 200

* Response 404 (application/json)

URL to be delete does not exist.

	{"errorCode": "url-not-found"}		

   
	
		






## References:

- HTTP server in go: https://gowebexamples.com/http-server/
