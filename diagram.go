package main

import (
	"strings"
	"os"
	"fmt"
	"path"
	"io/ioutil"
)

var outputFile *os.File
var resultFilepath string = "output.mmd"
var projectRoot string
var sceneCount int = 0
var linesOfCode int = 0


func main() {
	currentDir, _ := os.Getwd()
	projectRoot = findProjectRoot(currentDir)
	if projectRoot == "" {
		fmt.Println("Not inside Godot project. Could not find project.godot")
		return
	}

	fmt.Printf("Scanning code in %s...\n", currentDir)

	f, err := os.Create(resultFilepath)
	outputFile = f
	if err != nil {
		fmt.Println("Cannot create output file")
		os.Exit(1)
	}
	defer outputFile.Close()

	outputFile.WriteString("classDiagram\n")
	outputFile.Sync()

	err = scan(currentDir)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Scenes: %d\n", sceneCount)
	fmt.Printf("Lines of code: %d\n", linesOfCode)

	outputFile.Sync()
	outputFile.Close()
}

func findProjectRoot(currentPath string) (string) {
	testPath := path.Join(currentPath, "project.godot")
	_, err := os.Stat(testPath)
	if err == nil {
		return currentPath
	}
	
	currentPath = path.Dir(currentPath)
	if currentPath == "." || currentPath == "/" {
		return ""
	}
	return findProjectRoot(currentPath)
}

func scan(directory string) (error) {
	file, err := os.Open(directory)
	if err != nil {
		fmt.Printf("Could not open path %s", directory)
		return nil
	}
	
	files, err := file.Readdir(0)
	for _, f := range files {
		filename := f.Name()
		extension := path.Ext(filename)
		fullPath := path.Join(directory, filename)

		if extension == ".tscn" {
			err := parseScene(fullPath)
			if err != nil {
				fmt.Println(err)
			}
		}

		if f.IsDir() {
			err := scan(fullPath)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func parseScene(scenePath string) (error) {
	//fmt.Printf("Parsing %s...\n", scenePath)

	sceneName := cleanClassName(scenePath)
	classDefinition := fmt.Sprintf("\tclass %s\n", sceneName)

	_, err := outputFile.WriteString(classDefinition)
	if err != nil {
		return err
	}

	// Open scene file
	// Parse out scripts and packed scenes
	contents, err := ioutil.ReadFile(scenePath)
	if err != nil {
		return err
	}

	sceneCount++

	lines := strings.Split(string(contents), "\n")
	foundScript := false
	for _, line := range lines {
		if strings.Contains(line, "ext_resource") {
			// Tokenize this line and get the path and type parts
			resourceType := getType(line)
			resourcePath := getPath(line)
			if resourceType == "Script" && !foundScript {
				// Only export the first script found
				parseScript(scenePath, resourcePath)
				foundScript = true
			}
			if resourceType == "Script" {
				countScriptLines(resourcePath)
			}
			if resourceType == "PackedScene" {
				internalScene := cleanClassName(resourcePath)
				classDefinition := fmt.Sprintf("\t%s <|-- %s\n", sceneName, internalScene)
				_, err := outputFile.WriteString(classDefinition)
				if err != nil {
					return err
				}
			}
		}
	}
	
	return nil
}

func countScriptLines(scriptPath string) {
	contents, err := ioutil.ReadFile(path.Join(projectRoot, scriptPath))
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(contents), "\n")
	linesOfCode += len(lines)
}

func contains(list []string, needle string) (bool) {
	for _, val := range list {
		if val == needle {
			return true
		}
	}

	return false
}

func parseScript(scenePath string, scriptPath string) {
	sceneName := cleanClassName(scenePath)
	contents, err := ioutil.ReadFile(path.Join(projectRoot, scriptPath))
	if err != nil {
		fmt.Println(err)
		return
	}
	var varnames []string
	var funcnames []string
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) > 3 && line[0:3] == "var" {
			name := line
			if strings.Contains(name, "=") {
				name = strings.Split(name, "=")[0]
				name = strings.Split(name, "var ")[1]
			} else {
				name = strings.Split(name, "var ")[1]
			}
			name = strings.Split(name, ":")[0]
			if name[0:1] != "_" && name[0:1] == strings.ToUpper(name[0:1]) && !contains(varnames, name) {
				varnames = append(varnames, name)
				outputFile.WriteString(fmt.Sprintf("\t%s: +%s\n", sceneName, name))
			}
		}

		if len(line) > 6 && line[0:6] == "export" {
			name := line
			if strings.Contains(name, "=") {
				name = strings.Split(name, "=")[0]
				name = strings.Split(name, "var ")[1]
			} else {
				name = strings.Split(name, "var ")[1]
			}
			name = strings.Split(name, ":")[0]
			varnames = append(varnames, name)
			if !contains(varnames, name) {
				outputFile.WriteString(fmt.Sprintf("\t%s: +%s\n", sceneName, name))
			}
		}

		if len(line) > 5 && line[0:5] == "func " {
			name := line
			name = strings.Replace(name, "func ", "", -1)
			name = strings.Split(name, "(")[0]
			if name[0:1] != "_" && !contains(funcnames, name) {
				funcnames = append(funcnames, name)
				outputFile.WriteString(fmt.Sprintf("\t%s: +%s()\n", sceneName, name))
			}
		}
	}
}

func getPath(line string) (string) {
	path := ""
	fields := strings.Fields(line)
	for _, f := range fields {
		if strings.Contains(f, "path=") {
			path = strings.Replace(f, "path=", "", -1)
			path = strings.Replace(path, "res:/", "", -1)
			path = path[1:len(path)-1]
			return path
		}
	}
	return path
}

func getType(line string) (string) {
	resourceType := ""
	fields := strings.Fields(line)
	for _, f := range fields {
		if strings.Contains(f, "type=") {
			resourceType = strings.Replace(f, "type=", "", -1)
			resourceType = resourceType[1:len(resourceType)-1]
			return resourceType
		}
	}
	return resourceType
}

func cleanClassName(scenePath string) (string) {
	className := path.Base(scenePath)
	className = strings.Replace(className, ".tscn", "", -1)
	className = strings.Title(className)
	className = strings.Replace(className, "-", "", -1)

	return className
}