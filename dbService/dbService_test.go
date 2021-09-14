package dbService

import (
		"testing"
		"database/sql"
)


func TestOpenConnInitTables(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", "../database/testdb.db")

	if err != nil {
		t.Errorf("Open sql connection = %v; want db pointer", err)
	} 
}

func TestTableCreation(t *testing.T) {
	_, err := db.Exec("create table if not exists url (short_url text NOT NULL, long_url text, created text, expiration text, total_calls, PRIMARY KEY (short_url))")

	if err != nil {
		t.Errorf("Create table = %v; want nil", err)
	}
}

func TestSaveShortURL(t *testing.T) {
	longURL := "https://golang.org/doc/tutorial/getting-started"
	shortURL := "golang"
	expiration := "2021-09-15"
  
	res, err := SaveShortURL(shortURL, longURL, expiration)

	if res != 1 || err != nil {
		t.Errorf("SaveShortURL(shortURL, longURL, expiration) = %d; want 1", res)
	}
}

func TestGetURLTotalCalls(t * testing.T) {
	shortURL := "golang"
	calls, err := getURLTotalCalls(shortURL)

	if calls != 0 || err != nil {
		t.Errorf("getURLTotalCalls(shortURL) = %d; want 0", calls)
	}
}

func TestUpdateTotalCalls(t *testing.T) {
	shortURL := "golang"
	updatedCall, err := UpdateTotalCalls(shortURL)

	if updatedCall != 1 || err != nil {
		t.Errorf("UpdateTotalCalls(shortURL) = %d; want 1", updatedCall)
	}

}

func TestURLExistsQ(t *testing.T) {
	shortURL := "golang"
	exists, err := URLExistsQ(shortURL)

	if !exists || err != nil {
		t.Errorf("RandomKeyExistsQ(shortURL) = %v; want true", err)
	}
}	

func TestGetLongURL(t *testing.T) {
	shortURL := "golang"
	longURL, err := GetLongURL(shortURL)

	if longURL == "" || err != nil {
		t.Errorf("GetLongURL(shortURL) = %s; want https://golang.org/doc/tutorial/getting-started" , longURL)
	}
 }

func TestDeleteURL(t *testing.T) {
	shortURL := "golang"
	deleted, err := DeleteShortURL(shortURL)

	if deleted == 0 || err != nil {
		t.Errorf("DeleteShortURL(shortURL) = %d; want 1", deleted)
	}
}	
