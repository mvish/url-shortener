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

   
	
		



