package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// clean up old runs
	clean()
	// clone buckets
	buckets := []string{
		"https://github.com/ScoopInstaller/Main.git",
		"https://github.com/lukesampson/scoop-extras.git",
	}
	for i := 0; i < len(buckets); i++ {
		clone(buckets[i])
	}
	// get a list of all json files
	files, err := filepath.Glob("./*/bucket/*.json")
	catch(err, "", "")
	// read files
	filesArray := readFilesToArray(files)
	filesString := strings.Join(filesArray, ",")
	filesString = "[" + filesString + "]"
	// write to file
	log.Println("Writing manifests to file")
	write("manifests.json", filesString)
	log.Println("Done")
}

func write(filename string, content string) {
	file, err := os.Create("docs/" + filename)
	catch(err, "", "")
	defer file.Close()
	_, err = io.WriteString(file, content)
	catch(err, "", "")
}

func clone(url string) {
	log.Println("Cloning bucket repository " + url)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "clone", url)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	catch(err, stdout.String(), stderr.String())
}

func readFilesToArray(files []string) []string {
	log.Println("Reading manifests")
	var result []string
	for i := 0; i < len(files); i++ {
		dat, err := ioutil.ReadFile(files[i])
		catch(err, "", "")
		name, bucket := extractManifestDetails(files[i])
		manifest := addManifestDetails(string(dat), name, bucket)
		result = append(result, manifest)
	}
	return result
}

func addManifestDetails(manifest string, name string, bucket string) string {
	runes := []rune(manifest)
	manifest = string(runes[2 : len(runes)-1])
	manifest = "{ \"name\": \"" + name + "\", \"bucket\": \"" + bucket + "\", " + manifest
	return manifest
}

func extractManifestDetails(path string) (string, string) {
	pathParts := strings.Split(path, "\\")
	bucket := pathParts[0]
	nameWithJSON := pathParts[len(pathParts)-1]
	jsonParts := strings.Split(nameWithJSON, ".json")
	name := jsonParts[0]
	return name, bucket
}

func clean() {
	log.Println("Cleaning up previous runs (if any)")
	files, err := filepath.Glob("./*/*")
	catch(err, "", "")
	for i := 0; i < len(files); i++ {
		// don't delete our .git
		if !strings.HasPrefix(files[i], ".git") && !strings.HasPrefix(files[i], "docs") {
			err := os.RemoveAll(files[i])
			catch(err, "", "")
		}
	}
	log.Println("Cleaned up " + strconv.Itoa(len(files)) + " paths")
}

func catch(err error, stdout string, stderr string) {
	if err != nil {
		log.Println("ERROR")
		log.Println(err)
		if stdout != "" {
			log.Println(stdout)
		}
		if stderr != "" {
			log.Println(stderr)
		}
		log.Fatalln("Exiting shovel because of an error. Check the logging output above.")
	}
}
