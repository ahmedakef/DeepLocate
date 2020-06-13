package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

var pythonCode = ""

func executeScript(scriptName string, parameters string, object interface{}) error {
	cmd := exec.Command("python", pythonCode+scriptName, parameters)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	if err := json.NewDecoder(stdout).Decode(&object); err != nil {
		log.Error(err)
		return err
	}
	if err := cmd.Wait(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Person ahmed
type Person struct {
	Name     string
	FileName string
}

func main() {

	var person Person
	executeScript("foo.py", "funcky", &person)
	fmt.Printf("%s is %s on the disk\n", person.Name, person.FileName)

}
