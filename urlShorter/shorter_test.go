package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

const letterBytes = "./abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type testpair struct {
	url    string
	result string
}

var tests = []testpair{
	{"myURLShorter.com/q", "bash.im"},
	{"bash.im", ""},
	{"myURLShorter.com/qqq", ""},
}

func TestGetFullLoc(t *testing.T) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	for _, pair := range tests {
		v := us.getFullLoc(pair.url)
		if v != pair.result {
			t.Error(
				"For", pair.url,
				"expected", pair.result,
				"got", v,
			)
		}
	}
}

func TestGetExistingShortLoc(t *testing.T) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	v := us.getShortLoc("bash.im")
	if v != "myURLShorter.com/q" {
		t.Error(
			"For bash.im expected myURLShorter.com/q got", v,
		)
	}
}

func TestGetNewShortLoc(t *testing.T) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	v := us.getShortLoc(RandStringBytes(25))
	if !strings.HasPrefix(v, "myURLShorter.com") {
		t.Error(
			"For bash.im expected myURLShorter.com/* got", v,
		)
	}
}

//BenchmarkGetNewShort - benchark test
func BenchmarkGetNewShort(b *testing.B) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	for n := 0; n < b.N; n++ {
		us.getShortLoc(RandStringBytes(10))
	}
}

//BenchmarkGetExistingShort - benchark test
func BenchmarkGetExistingShort(b *testing.B) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	for n := 0; n < b.N; n++ {
		us.getShortLoc("https://www.youtube.com/")
	}
}

//BenchmarkGetExistingLong - benchark test
func BenchmarkGetExistingLong(b *testing.B) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	for n := 0; n < b.N; n++ {
		us.getShortLoc("myURLShorter.com/p")
	}
}

func BenchmarkGenerateShortURL(b *testing.B) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	for n := 0; n < b.N; n++ {
		us.generateShortURL()
	}
}

func BenchmarkGetMaxID(b *testing.B) {
	rep, err := makeRepos()
	if err != nil {
		fmt.Println("error in makeRepos", err)
	}
	defer rep.Close()
	us := makeURLShorter(rep)
	for n := 0; n < b.N; n++ {
		us.rep.GetMaxID()
	}
}
