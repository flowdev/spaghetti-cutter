package main

import (
	"log"
	"os"
)

func main() {
	log.Printf("INFO - this is the main package, args: %q", os.Args[1:])
