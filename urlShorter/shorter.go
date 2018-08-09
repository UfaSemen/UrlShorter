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

type urlShorter struct{}

type url struct {
	id       int
	shorturl string
	longurl  string
}

const (
	alphabet = "abcdfghjkmnpqrstvwxyzABCDFGHJKLMNPQRSTVWXYZ-_0123456789"
	base     = len(alphabet)
	domen    = "myURLShorter.com/"
)

func (g *urlShorter) GetShort(ctx context.Context, req *proto.ShortRequest, rsp *proto.ShortResponse) error {
	db, err := sql.Open("mysql", "root:12345@/urlshorterdb")
	if err != nil {
		fmt.Println("sql.Open", err)
	}
	defer db.Close()
	/*
		Search for existing short url in database,
		if there isn't one insert full url in it,
		compute short url from it's id and update inserted row
	*/
	p := url{}
	if err := db.QueryRow("SELECT * FROM urlshorterdb.url WHERE longurl = ?", req.Url).Scan(&p.id, &p.shorturl, &p.longurl); err == nil {
		rsp.Url = domen + p.shorturl

	} else if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO urlshorterdb.url (shorturl, longurl) VALUES('', ?)", req.Url)
		if err != nil {
			fmt.Println("db.Exec", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("LastInsertId()", err)
		}

		str := ""
		num := int(id)
		for num > 0 {
			str = string(alphabet[num%base]) + str
			num /= base
		}

		_, err = db.Exec("UPDATE urlshorterdb.url SET shorturl = ? WHERE id = ?", str, id)
		if err != nil {
			fmt.Println(" db.Exec UPDATE", err)
		}

		rsp.Url = domen + str

	} else {
		fmt.Println("db.QueryRow", err)
	}

	return nil
}

func (g *urlShorter) GetFull(ctx context.Context, req *proto.FullRequest, rsp *proto.FullResponse) error {

	db, err := sql.Open("mysql", "root:12345@/urlshorterdb")
	if err != nil {
		fmt.Println("sql.Open", err)
	}
	defer db.Close()
	/*
		compute id from url and look for it in database
	*/
	if strings.HasPrefix(req.Url, domen) {
		var num int
		code := req.Url[17:]
		for _, v := range code {
			cInd := strings.IndexRune(alphabet, v)

			num = num*base + cInd
		}

		p := url{}
		err = db.QueryRow("SELECT * FROM urlshorterdb.url WHERE id = ?", num).Scan(&p.id, &p.shorturl, &p.longurl)
		if err != nil {
			fmt.Println("db.QueryRow", err)
		}

		rsp.Url = p.longurl
	}
	return nil
}

func main() {

	service := micro.NewService(
		micro.Name("urlShorter"),
		micro.Version("latest"),
	)
	service.Init()
	proto.RegisterUrlShorterHandler(service.Server(), new(urlShorter))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
