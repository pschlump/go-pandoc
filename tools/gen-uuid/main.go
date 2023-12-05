package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	u := uuid.New().String()
	fmt.Printf("%s\n", u)
}
