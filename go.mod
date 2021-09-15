module mili.photos/url-shortener

go 1.17

replace mili.photos/service => ./service

replace mili.photos/dbService => ./dbService

require (
	mili.photos/dbService v0.0.0-00010101000000-000000000000
	mili.photos/service v0.0.0-00010101000000-000000000000
)

require (
	github.com/flyaways/golang-lru v0.0.0-20190617091412-ec8b77fcae6c // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
)
