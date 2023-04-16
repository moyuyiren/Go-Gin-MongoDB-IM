package main

import "GoIm/router"

func main() {
	r := router.Router()
	r.Run(":8080")
}
