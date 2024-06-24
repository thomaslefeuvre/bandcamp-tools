package main

import (
	"fmt"
	"os"

	"github.com/thomaslefeuvre/bandcamp-tools/bandcamp"
)

const (
	username         = "?"
	wishlistLocation = "?"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	command := os.Args[1]

	if command == "count" {
		// items, err := bandcamp.LoadWishlist(wishlistLocation)
		// if err != nil {
		// 	fmt.Println("Error loading wishlist:", err)
		// 	return
		// }

		// fmt.Println("Number of items:", len(items))
		return
	}

	if command == "sync" {
		err := bandcamp.SyncWishlist(username, wishlistLocation)
		if err != nil {
			fmt.Println("Error syncing wishlist:", err)
		}
		return
	}

	// command not matched
	usage()
}

func usage() {
	fmt.Println("Usage: bc <command>")
}
