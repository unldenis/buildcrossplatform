package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	records, err := readCsvFile("target.csv")
	if err != nil {
		log.Fatal(err)
	}

	var dir string
	getInput(&dir, "project folder")

	var nameFile string
	getInput(&nameFile, "build project name")

	var version string
	getInput(&version, "build project version")

	for i := 1; i < len(records); i++ {
		osVal := records[i][0]
		archVal := records[i][1]

		build(&dir, &nameFile, &version, &osVal, &archVal)
	}
}

func getInput(variable *string, name string) {
	fmt.Printf("Enter the %s > ", name)
	_, err := fmt.Scanln(variable)
	if err != nil {
		log.Fatal(err)
	}
}

func build(dir, nameFile, versionFile, osVal, archVal *string) {
	err := os.Setenv("GOOS", *osVal)
	if err != nil {
		fmt.Printf("GOOS Error in %v-%v: %v\n", *osVal, *archVal, err)
		return
	}
	err = os.Setenv("GOARCH", *archVal)
	if err != nil {
		fmt.Printf("GOARCH Error in %v-%v: %v\n", *osVal, *archVal, err)
		return
	}

	targetFile := fmt.Sprintf("./build/%sv%s-%s_%s", *nameFile, *versionFile, *osVal, *archVal)
	if *osVal == "windows" {
		targetFile += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", targetFile)
	cmd.Dir = *dir
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Fatal Error in %v-%v: %v\n", *osVal, *archVal, err)
	} else {
		fmt.Printf("%s", out)
	}
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Join(errors.New("Unable to read input file "+filePath), err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Join(errors.New("Unable to parse file as CSV for "+filePath), err)
	}

	return records, nil
}
