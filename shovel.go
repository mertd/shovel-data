package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Jeffail/gabs/v2"
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
	filesArray := parseManifests(files)
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
	repo string
}

var gitHubURL = "https://github.com/"

func getBuckets() []Bucket {
	buckets := []Bucket{
		{"main", "ScoopInstaller/Main"},
		{"extras", "lukesampson/scoop-extras"},
		{"versions", "ScoopInstaller/Versions"},
		{"nirsoft", "kodybrown/scoop-nirsoft"},
		{"php", "ScoopInstaller/PHP"},
		{"nerd-fonts", "matthewjberger/scoop-nerd-fonts"},
		{"nonportable", "TheRandomLabs/scoop-nonportable"},
		{"java", "ScoopInstaller/Java"},
		{"games", "Calinou/scoop-games"},
	}
	return buckets
}

func cloneBuckets() {
	buckets := getBuckets()
	for i := 0; i < len(buckets); i++ {
		clone(buckets[i])
	}
}

func write(fileName string, content string) {
	file, err := os.Create("docs/" + fileName)
	catch(err, "", "")
	defer file.Close()
	_, err = io.WriteString(file, content)
	catch(err, "", "")
}

func clone(bucket Bucket) {
	log.Println("Cloning bucket repository " + bucket.name)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "clone", gitHubURL+bucket.repo, workDir+"/"+bucket.name)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	catch(err, stdout.String(), stderr.String())
}

func parseManifests(files []string) []string {
	log.Println("Parsing manifests")
	var result []string
	errorCount := 0
	successCount := 0
	for i := 0; i < len(files); i++ {
		index := gabs.New()
		// set filtered manifest keys
		manifest, err := gabs.ParseJSONFile(files[i])
		homepage, version, description, github := extractManifestContents(manifest);
		index.SetP(homepage, "homepage");
		index.SetP(version, "version");
		index.SetP(description, "description");
		if github != "" {
			index.SetP(github, "checkver.github");
		}
		// set custom keys
		name, bucket, manifestURL := extractManifestDetails(files[i])
		index.SetP(name, "name")
		index.SetP(bucket, "bucket")
		index.SetP(manifestURL, "manifestURL")
		if err == nil {
			result = append(result, index.String())
			successCount = successCount + 1
		} else {
			log.Println("Skipping", name, "from", bucket, "--", err)
			errorCount = errorCount + 1
		}
	}
	log.Println("Successfully parsed", successCount, "manifest(s).")
	log.Println("Skipped", errorCount, "erroneous manifest(s).")
	return result
}

func extractManifestContents(manifest *gabs.Container) (string, string, string, string) {
	var homepage string
	var version string
	var description string
	var github string
	var ok bool
	homepage, ok = manifest.Path("homepage").Data().(string)
	if ok == false {
		homepage = ""
	}
	version, ok = manifest.Path("version").Data().(string)
	if ok == false {
		version = ""
	}
	description, ok = manifest.Path("description").Data().(string)
	if ok == false {
		description = ""
	}
	github, ok = manifest.Path("checkver.github").Data().(string)
	if ok == false {
		github = ""
	}
	return homepage, version, description, github
}

func extractManifestDetails(path string) (string, string, string) {
	// extract from filename
	separator := string(os.PathSeparator)
	parts := strings.Split(path, separator)
	bucket := parts[1] // workDir/bucket
	nameWithJSON := parts[len(parts)-1]
	jsonParts := strings.Split(nameWithJSON, ".json")
	name := jsonParts[0]
	// extract from repository url
	buckets := getBuckets()
	var repo string
	for i := 0; i < len(buckets); i++ {
		if buckets[i].name == bucket {
			repo = buckets[i].repo

		}
	}
	manifestURL := gitHubURL + repo + "/tree/master/bucket/" + name + ".json"
	return name, bucket, manifestURL
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
