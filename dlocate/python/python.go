package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("python", "foo.py", "khaled.pdf")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var person struct {
		Name     string
		FileName string
	}
	if err := json.NewDecoder(stdout).Decode(&person); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is %s in the disk\n", person.Name, person.FileName)
}
