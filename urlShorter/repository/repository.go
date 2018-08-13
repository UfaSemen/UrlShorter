package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type url struct {
	id       int
	shorturl string
	longurl  string
}

//Repository struct
type Repository struct {
	db           *sql.DB
	urlList      []url
	dbDataLength int
}

//MakeRepos connects to database and gets data from it
func MakeRepos() (*Repository, error) {
	rep := new(Repository)
	var err error
	rep.db, err = sql.Open("mysql", "root:12345@/urlshorterdb")
	if err != nil {
		return nil, err
	}

	rows, err := rep.db.Query("SELECT * FROM urlshorterdb.url")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u url
		err = rows.Scan(&u.id, &u.shorturl, &u.longurl)
		if err != nil {
			return nil, err
		}
		rep.urlList = append(rep.urlList, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rep.dbDataLength = len(rep.urlList)
	return rep, err
}

//Close connects to database and gets data from it
func (rep *Repository) Close() {
	newData := rep.urlList[rep.dbDataLength:]
	for _, url := range newData {
		_, err := rep.db.Exec("INSERT INTO urlshorterdb.url (shorturl, longurl) VALUES(?, ?)", url.shorturl, url.longurl)
		if err != nil {
			fmt.Println("repository.close(), insert error: ", err)
		}
	}
	rep.db.Close()
}

//GetShortByLongURL gets short url by long
func (rep *Repository) GetShortByLongURL(long string) string {
	for i := range rep.urlList {
		if rep.urlList[i].longurl == long {
			return rep.urlList[i].shorturl
		}
	}
	return ""
}

//GetLongByShortURL gets long url by short
func (rep *Repository) GetLongByShortURL(short string) string {
	for i := range rep.urlList {
		if rep.urlList[i].shorturl == short {
			return rep.urlList[i].longurl
		}
	}
	return ""
}

//GetMaxID gets maximum id
func (rep *Repository) GetMaxID() int {
	return rep.urlList[len(rep.urlList)-1].id
}

//InsertURL adds new url pair
func (rep *Repository) InsertURL(short string, long string) {
	u := url{rep.GetMaxID() + 1, short, long}
	rep.urlList = append(rep.urlList, u)
}
