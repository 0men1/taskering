package logging

import (
	"log"
	"os"
)

func LogToFile(name string, data string) error {
	if !FileExists(name) {
		createFile(name)
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(data)

	return nil
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()
	return nil
}
