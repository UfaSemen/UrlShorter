package main

import (
	"URLShorter/URLShorter/converter"
	"URLShorter/URLShorter/repository"
	"fmt"

	proto "URLShorter/URLShorter/proto"

	micro "github.com/micro/go-micro"
)

func main() {
	rep, err := repository.MakeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()

	service := micro.NewService(
		micro.Name("urlShorter"),
		micro.Version("latest"),
	)
	service.Init()
	proto.RegisterUrlShorterHandler(service.Server(), converter.MakeURLShorter(rep))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
