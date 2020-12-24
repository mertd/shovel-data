package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// configure
var workDir = ".work"

func main() {
	prepareWorkDir()
	cloneBuckets()
	// get a list of all json files
	files, err := filepath.Glob(workDir + "/*/bucket/*.json")
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

// A Bucket consists of its name and a git url
type Bucket struct {
	name string
	url  string
}

func cloneBuckets() {
	buckets := []Bucket{
		Bucket{"main", "https://github.com/ScoopInstaller/Main"},
		Bucket{"extras", "https://github.com/lukesampson/scoop-extras"},
		Bucket{"versions", "https://github.com/ScoopInstaller/Versions"},
		// Bucket{"nightlies", "https://github.com/ScoopInstaller/Nightlies"}, contains faulty manifest
		Bucket{"nirsoft", "https://github.com/kodybrown/scoop-nirsoft"},
		Bucket{"php", "https://github.com/ScoopInstaller/PHP"},
		Bucket{"nonportable", "https://github.com/TheRandomLabs/scoop-nonportable"},
		Bucket{"java", "https://github.com/ScoopInstaller/Java"},
		Bucket{"games", "https://github.com/Calinou/scoop-games"},
	}
	for i := 0; i < len(buckets); i++ {
		clone(buckets[i])
	}
}

func write(filename string, content string) {
	file, err := os.Create("docs/" + filename)
	catch(err, "", "")
	defer file.Close()
	_, err = io.WriteString(file, content)
	catch(err, "", "")
}

func clone(bucket Bucket) {
	log.Println("Cloning bucket repository " + bucket.name)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "clone", bucket.url, workDir+"/"+bucket.name)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	catch(err, stdout.String(), stderr.String())
}

func readFilesToArray(files []string) []string {
	log.Println("Reading manifests and gathering additional data")
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
	separator := string(os.PathSeparator)
	pathParts := strings.Split(path, separator)
	bucket := pathParts[0]
	nameWithJSON := pathParts[len(pathParts)-1]
	jsonParts := strings.Split(nameWithJSON, ".json")
	name := jsonParts[0]
	return name, bucket
}

func prepareWorkDir() {
	log.Println("Preparing work directory")
	removeErr := os.RemoveAll(workDir)
	if !os.IsNotExist(removeErr) {
		catch(removeErr, "", "")
	}
	createErr := os.Mkdir(workDir, 0755)
	if !os.IsExist(createErr) {
		catch(createErr, "", "")
	}
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
