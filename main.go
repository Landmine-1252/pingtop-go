package main

import (
	"os"

	"github.com/Landmine-1252/pingtop-go/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
