package main

import (
	"flag"
	"sacco/whatsapp"
)

func main() {
	var phoneNumber string
	var debug bool

	flag.StringVar(&phoneNumber, "p", phoneNumber, "phone number")

	flag.Parse()

	if phoneNumber == "" {
		panic("Missing required phoneNumber")
	}

	whatsapp.Main(phoneNumber, debug)
}
