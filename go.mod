module mili.photos/url-shortener

go 1.17

replace mili.photos/service => ./service

replace mili.photos/dbService => ./dbService

require (
	mili.photos/dbService v0.0.0-00010101000000-000000000000
	mili.photos/service v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
)
