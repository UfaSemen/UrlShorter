package main

import (
	"context"
	"fmt"
	"math/rand"

	proto "github.com/UfaSemen/urlShorter/urlShorter/proto"

	micro "github.com/micro/go-micro"
)

type test struct {
	option int
	url    string
}

const letterBytes = "./abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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
		{1, randStringBytes(10)},
		{2, "myURLShorter.com/qqq"},
		{1, "bash.im"},
		{2, "myURLShorter.com/s"},
	}

	for _, t := range tests {
		if t.option == 1 {
			rsp, err := urlShorter.GetShort(context.TODO(), &proto.ShortRequest{Url: t.url})
			if err != nil {
				fmt.Println("GetShort", err)
			}
			fmt.Println(rsp.Url)
		} else if t.option == 2 {
			rsp2, err := urlShorter.GetFull(context.TODO(), &proto.FullRequest{Url: t.url})
			if err != nil {
				fmt.Println("GetFull", err)
			}
			fmt.Println(rsp2.Url)
		}
	}
}
