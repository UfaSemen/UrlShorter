package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	proto "urlShorter/urlShorter/proto"

	micro "github.com/micro/go-micro"

	_ "github.com/go-sql-driver/mysql"
)

const (
	alphabet = "abcdfghjkmnpqrstvwxyzABCDFGHJKLMNPQRSTVWXYZ-_0123456789"
	base     = len(alphabet)
	domen    = "myURLShorter.com/"
)

type urlShorter struct {
	rep *repository
}

type url struct {
	id       int
	shorturl string
	longurl  string
}

type repository struct {
	db           *sql.DB
	urlList      []url
	dbDataLength int
}

func makeURLShorter(rep *repository) *urlShorter {
	shorter := new(urlShorter)
	shorter.rep = rep
	return shorter
}

func makeRepos() (*repository, error) {
	rep := new(repository)
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

func (rep *repository) Close() {
	newData := rep.urlList[rep.dbDataLength:]
	for _, url := range newData {
		_, err := rep.db.Exec("INSERT INTO urlshorterdb.url (shorturl, longurl) VALUES(?, ?)", url.shorturl, url.longurl)
		if err != nil {
			fmt.Println("repository.close(), insert error: ", err)
		}
	}
	rep.db.Close()
}

func (rep *repository) GetShortByLongURL(long string) string {
	for i := range rep.urlList {
		if rep.urlList[i].longurl == long {
			return rep.urlList[i].shorturl
		}
	}
	return ""
}

func (rep *repository) GetLongByShortURL(short string) string {
	for i := range rep.urlList {
		if rep.urlList[i].shorturl == short {
			return rep.urlList[i].longurl
		}
	}
	return ""
}

func (rep *repository) GetMaxID() int {
	return rep.urlList[len(rep.urlList)-1].id
}

func (rep *repository) InsertURL(short string, long string) {
	u := url{rep.GetMaxID() + 1, short, long}
	rep.urlList = append(rep.urlList, u)
}

func (g *urlShorter) generateShortURL() string {
	var str string
	num := g.rep.GetMaxID() + 1
	for num > 0 {
		str = string(alphabet[num%base]) + str
		num /= base
	}
	return str
}

func (g *urlShorter) getShortLoc(req string) string {
	if short := g.rep.GetShortByLongURL(req); short != "" {
		return domen + short
	}
	short := g.generateShortURL()
	g.rep.InsertURL(short, req)
	return domen + short

}

func (g *urlShorter) getFullLoc(req string) string {
	if strings.HasPrefix(req, domen) {
		return g.rep.GetLongByShortURL(req[len(domen):])
	}

	fmt.Println("wrong input")
	return ""
}

func (g *urlShorter) GetShort(ctx context.Context, req *proto.ShortRequest, rsp *proto.ShortResponse) error {
	rsp.Url = g.getShortLoc(req.Url)
	return nil
}

func (g *urlShorter) GetFull(ctx context.Context, req *proto.FullRequest, rsp *proto.FullResponse) error {
	rsp.Url = g.getFullLoc(req.Url)
	return nil
}

func main() {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()

	service := micro.NewService(
		micro.Name("urlShorter"),
		micro.Version("latest"),
	)
	service.Init()
	proto.RegisterUrlShorterHandler(service.Server(), makeURLShorter(rep))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
