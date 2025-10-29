package main

import "fmt"

func main() {
	for _, p := range getProjects() {
		fmt.Printf("%+v\n", p)
	}
}
