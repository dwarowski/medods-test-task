package main

import (
	"fmt"

	app "github.com/dwarowski/medods-test-task/src"
)

func main() {
	r, err := app.SetupApp()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.Run(":8080")
}
