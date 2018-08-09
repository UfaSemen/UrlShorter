package main

import (
	"context"
	"fmt"

	proto "urlShorter/urlShorter/proto"

	micro "github.com/micro/go-micro"
)

type test struct {
	option int
	url    string
}

func main() {

	service := micro.NewService(micro.Name("urlShorter.client"))
	service.Init()
	urlShorter := proto.NewUrlShorterService("urlShorter", service.Client())

	tests := []test{
		{1, "bash.im"},
		{2, "myURLShorter.com/q"},
		{1, "https://www.youtube.com/"},
		{2, "bash.im"},
		{1, "https://www.gismeteo.ru/weather-sankt-peterburg-4079/month/"},
		{2, "myURLShorter.com/qqq"},
		{1, "bash.im"},
		{2, "myURLShorter.com/s"},
	}

	for _, t := range tests {
		if t.option == 1 {
			rsp, err := urlShorter.GetShort(context.TODO(), &proto.ShortRequest{Url: t.url})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(rsp.Url)
		} else if t.option == 2 {
			rsp2, err2 := urlShorter.GetFull(context.TODO(), &proto.FullRequest{Url: t.url})
			if err2 != nil {
				fmt.Println(err2)
			}
			fmt.Println(rsp2.Url)
		}
	}
}
