package main

import "github.com/azka-zaydan/article-materials/unit-testing/infras"

func main() {
	err := infras.InitDB()
	if err != nil {
		panic(err)
	}
}
