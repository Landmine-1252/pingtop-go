package main

import (
	"os"

	"pingtop/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
