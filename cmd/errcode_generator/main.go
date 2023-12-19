package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ErrorStatus represents the structure of the input JSON data.
type ErrorStatus struct {
	StatusName string `json:"status_name"`
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// GenerateGoCode generates Go code from the provided JSON data.
func GenerateGoCode(data []ErrorStatus) string {
	var codeBuilder strings.Builder

	codeBuilder.WriteString("package errcode\n\n")
	codeBuilder.WriteString("var (\n")

	for _, status := range data {
		line := fmt.Sprintf("\t%s = NewError(%s, \"%s\")\n", status.StatusName, status.StatusCode, status.StatusMsg)
		codeBuilder.WriteString(line)
	}

	codeBuilder.WriteString(")\n")

	return codeBuilder.String()
}

func findJsonFile(directory string) []string {

	// Slice to store JSON file paths
	var jsonFiles []string

	// Use the Walk function to traverse all files in the folder
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// Check if it is a JSON file
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			// Append the JSON file path to the slice
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
	return jsonFiles
}

func err_checker(prevMsg string, err error) {
	if err != nil {
		fmt.Println(prevMsg, err)
		os.Exit(1)
	}
}

func fileNameProcessor(path string) string {
	fileName := filepath.Base(path)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func main() {
	// Check if the folder path parameter is provided
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go /path/to/your/directory /path/output/directory")
		os.Exit(1)
	}

	// User-provided folder path
	directory := os.Args[1]
	targetFolder := os.Args[2]
	// Find all json files in the directory
	jsonFiles := findJsonFile(directory)

	for _, file := range jsonFiles {
		fmt.Println("Read json file:", file)
		var statusData []ErrorStatus

		fileReader, fr_err := os.Open(file)
		err_checker("Read file err:", fr_err)
		defer fileReader.Close()

		byteValue, _ := ioutil.ReadAll(fileReader)
		unmarshal_err := json.Unmarshal([]byte(byteValue), &statusData)
		err_checker("Error parsing JSON:", unmarshal_err)

		goCode := GenerateGoCode(statusData)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(goCode)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

		// Save code generation in target directory
		outputFile := fmt.Sprintf("%s/%s.go", targetFolder, fileNameProcessor(file))
		fileWriter, fw_err := os.Create(outputFile)
		err_checker("Write file err:", fw_err)
		defer fileWriter.Close()

		_, w_err := fileWriter.WriteString(goCode)
		err_checker("Write err:", w_err)
	}
}
