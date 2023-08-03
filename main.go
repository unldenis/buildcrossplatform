package main

import (
	"encoding/csv"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	records, err := readCsv()
	if err != nil {
		fmt.Print(err)
		wait()
		return
	}

	var dir string
	getInput(&dir, "project folder")

	var nameFile string
	getInput(&nameFile, "build project name")

	var version string
	getInput(&version, "build project version")

	bar := progressbar.Default(43)

	var errors []string

	for i := 1; i < len(records); i++ {
		osVal := records[i][0]
		archVal := records[i][1]

		build(&dir, &nameFile, &version, &osVal, &archVal, &errors)

		err = bar.Add(1)
		if err != nil {
			fmt.Printf("Error adding to progress bar %v\n", err)
			wait()
			return
		}

	}
	for _, e := range errors {
		fmt.Println(e)
	}
	wait()
}

func wait() {
	fmt.Println()
	var wait string
	fmt.Scanln(&wait)
}

func getInput(variable *string, name string) {
	fmt.Printf("Enter the %s > ", name)
	_, err := fmt.Scanln(variable)
	if err != nil {
		log.Fatal(err)
	}
}

func build(dir, nameFile, versionFile, osVal, archVal *string, errors *[]string) {
	err := os.Setenv("GOOS", *osVal)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("GOOS Error in %v-%v: %v", *osVal, *archVal, err))
		return
	}
	err = os.Setenv("GOARCH", *archVal)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("GOARCH Error in %v-%v: %v", *osVal, *archVal, err))
		return
	}

	targetFile := fmt.Sprintf("./build/%sv%s-%s_%s", *nameFile, *versionFile, *osVal, *archVal)
	if *osVal == "windows" {
		targetFile += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", targetFile)
	cmd.Dir = *dir
	_, err = cmd.Output()
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("Fatal Error in %v-%v: %v", *osVal, *archVal, err))
	}
}

func readCsv() ([][]string, error) {
	csvReader := csv.NewReader(strings.NewReader(osarch))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

const osarch = `
os,arch
aix,ppc64
android,386
android,amd64
android,arm
android,arm64
darwin,amd64
darwin,arm64
dragonfly,amd64
freebsd,386
freebsd,amd64
freebsd,arm
illumos,amd64
ios,arm64
js,wasm
linux,386
linux,amd64
linux,arm
linux,arm64
linux,loong64
linux,mips
linux,mipsle
linux,mips64
linux,mips64le
linux,ppc64
linux,ppc64le
linux,riscv64
linux,s390x
netbsd,386
netbsd,amd64
netbsd,arm
openbsd,386
openbsd,amd64
openbsd,arm
openbsd,arm64
plan9,386
plan9,amd64
plan9,arm
solaris,amd64
wasip1,wasm
windows,386
windows,amd64
windows,arm
windows,arm64
`
