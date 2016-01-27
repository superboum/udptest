package main

import (
	"log"
)

/* A Simple function to verify error */
func CheckError(err error) bool {
	if err != nil {
		log.Println("Error: ", err)
	}
	return err == nil
}
