package main

import (
	. "microgram/pkg/api"
)

func main() {
	theApp := App{}

	//Initialize with DB
	theApp.Initialize("microgram", "yourpassword", "microgram")

	theApp.Run(":8080")
}
