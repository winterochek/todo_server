package main

import (
	"log"

	api "github.com/winterochek/todo-server/internal/app/server"
)

func main(){
	err := api.Start()
	if err != nil {
		log.Fatal(err)
	}
}