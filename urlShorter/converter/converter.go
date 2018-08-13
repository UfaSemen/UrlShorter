package converter

import (
	proto "URLShorter/URLShorter/proto"
	"URLShorter/URLShorter/repository"
	"context"
	"fmt"
	"strings"
)

const (
	alphabet = "abcdfghjkmnpqrstvwxyzABCDFGHJKLMNPQRSTVWXYZ-_0123456789"
	base     = len(alphabet)
	domen    = "myURLShorter.com/"
)

type URLShorter struct {
	rep *repository.Repository
}

func MakeURLShorter(rep *repository.Repository) *URLShorter {
	shorter := new(URLShorter)
	shorter.rep = rep
	return shorter
}

func (g *URLShorter) generateShortURL() string {
	var str string
	num := g.rep.GetMaxID() + 1
	for num > 0 {
		str = string(alphabet[num%base]) + str
		num /= base
	}
	return str
}

func (g *URLShorter) getShortLoc(req string) string {
	if short := g.rep.GetShortByLongURL(req); short != "" {
		return domen + short
	}
	short := g.generateShortURL()
	g.rep.InsertURL(short, req)
	return domen + short

}

func (g *URLShorter) getFullLoc(req string) string {
	if strings.HasPrefix(req, domen) {
		return g.rep.GetLongByShortURL(req[len(domen):])
	}

	fmt.Println("wrong input")
	return ""
}

func (g *URLShorter) GetShort(ctx context.Context, req *proto.ShortRequest, rsp *proto.ShortResponse) error {
	rsp.Url = g.getShortLoc(req.Url)
	return nil
}

func (g *URLShorter) GetFull(ctx context.Context, req *proto.FullRequest, rsp *proto.FullResponse) error {
	rsp.Url = g.getFullLoc(req.Url)
	return nil
}
