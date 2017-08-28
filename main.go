package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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

// from stackoverflow answer https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
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

	err = copyFileContents(configs.inputFile, configs.outputFile)
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
