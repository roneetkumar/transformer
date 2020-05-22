package main

import (
	"fmt"
)

func main() {

	out, err := primitive("./header.jpg", "out.png", 100, ellipse)

	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
