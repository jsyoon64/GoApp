package main

import (
	"fmt"

	"go-multiple-directories-example/models"
	"go-multiple-directories-example/routes"
)

func main() {
	fmt.Println("Main package - main file")
	models.AllUsers()
	routes.APIPostRoute()
}
