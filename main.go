package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type ConfigsModel struct {
	inputFile  string
	outputFile string
	moveAction string
}

func createConfigsModelFromEnvs() (ConfigsModel, error) {

	ret := ConfigsModel{
		inputFile:  os.Getenv("input_file"),
		outputFile: os.Getenv("output_file"),
		moveAction: os.Getenv("move_action"),
	}
	return ret, nil
}

func (configs ConfigsModel) print() {
	fmt.Println()
	log.Printf("Configs:")
	log.Printf("%+v", configs)
}

func (configs ConfigsModel) validate() error {
	// required
	if configs.inputFile == "" {
		return errors.New("no inputfile parameter specified")
	}
	if configs.outputFile == "" {
		return errors.New("no outputFile parameter specified")
	}

	if configs.moveAction != "move" && configs.moveAction != "copy" {
		return errors.New("wrong move action: " + configs.moveAction)
	}

	return nil
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	containingDir := filepath.Dir(dest)

	err = os.MkdirAll(containingDir, os.ModePerm)

	if err != nil {
		return err
	}

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err == nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func copyFolder(source string, dest string) (err error) {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	if err != nil {
		return err
	}

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = copyFolder(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		}

	}
	return
}

func main() {
	configs, err := createConfigsModelFromEnvs()

	if err != nil {
		log.Fatalf("Issue with input: %s", err)
	}

	configs.print()
	if err = configs.validate(); err != nil {
		log.Fatalf("Issue with input: %s", err)
	}

	// check if target path exists

	if ex, err := exists(configs.outputFile); err != nil && ex == false {
		os.MkdirAll(configs.outputFile, os.ModePerm)
	}

	fi, err := os.Stat(configs.inputFile)

	if err != nil {
		log.Fatalf("Error : %s", err)
	}

	if fi.IsDir() {
		err = copyFolder(configs.inputFile, configs.outputFile)
	} else {
		err = copyFile(configs.inputFile, configs.outputFile)
	}

	if err != nil {
		log.Fatalf("Copying failure: %s", err)
	}

	if configs.moveAction == "move" {
		err = os.Remove(configs.inputFile)
		if err != nil {
			log.Fatalf("Removing file failure: %s", err)
		}

	}

}
