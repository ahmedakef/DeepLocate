package python

import (
	"encoding/json"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

var pythonDirectory = "../ImageCaptioningAndKeyWordExtraction/"

// ExecuteScript run python script and decode its stdout to object
func ExecuteScript(scriptName string, parameters string, object interface{}) error {
	cmd := exec.Command("python", "-W ignore", pythonDirectory+scriptName, parameters)

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

var filesContent map[string]map[string]float32

func main() {

	ExecuteScript("Extract.py", "/home/ahmed/Downloads/cloud computing/", &filesContent)

	log.Info(filesContent)

	// var person Person
	// ExecuteScript("foo.py", "funcky", &person)
	// log.Infof("%s is %s on the disk\n", person.Name, person.FileName)

}
