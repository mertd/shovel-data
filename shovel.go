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
	cleanOldRuns()
	cloneBuckets()
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

type Bucket struct {
	name string
	url  string
}

func cloneBuckets() {
	buckets := []Bucket{
		Bucket{"main", "https://github.com/ScoopInstaller/Main"},
		Bucket{"extras", "https://github.com/lukesampson/scoop-extras"},
		Bucket{"versions", "https://github.com/ScoopInstaller/Versions"},
		Bucket{"nightlies", "https://github.com/ScoopInstaller/Nightlies"},
		Bucket{"nirsoft", "https://github.com/kodybrown/scoop-nirsoft"},
		Bucket{"php", "https://github.com/ScoopInstaller/PHP"},
		Bucket{"nerd-fonts", "https://github.com/matthewjberger/scoop-nerd-fonts"},
		Bucket{"nonportable", "https://github.com/TheRandomLabs/scoop-nonportable"},
		Bucket{"java", "https://github.com/ScoopInstaller/Java"},
		Bucket{"games", "https://github.com/Calinou/scoop-games"},
		Bucket{"jetbrains", "https://github.com/Ash258/Scoop-JetBrains"},
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
	cmd := exec.Command("git", "clone", bucket.url, bucket.name)
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
	pathParts := strings.Split(path, "\\")
	bucket := pathParts[0]
	nameWithJSON := pathParts[len(pathParts)-1]
	jsonParts := strings.Split(nameWithJSON, ".json")
	name := jsonParts[0]
	return name, bucket
}

func cleanOldRuns() {
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
