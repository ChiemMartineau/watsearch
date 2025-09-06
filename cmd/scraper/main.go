package main

import (
	"fmt"
	"log"
)

func main() {
	courses, err := getKualiCourses()
	if err != nil {
		log.Fatalf("failed to get courses: %v", err)
	}

	fmt.Printf("%v", courses)
}
