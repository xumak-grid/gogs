package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)

var gogsHost = "http://localhost:3000"

// serviceReady waits untile the gogsHost is ready wiht 200 OK
func serviceReady() {
	statusCode := 0
	for statusCode != 200 {
		res, err := http.Get(gogsHost)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 5)
		} else {
			fmt.Println(res.StatusCode)
			statusCode = res.StatusCode
		}
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		fmt.Printf("%s:%s\n", key, defaultValue)
		return defaultValue
	}
	fmt.Printf("%s:%s", key, value)
	return value
}

const gitSetupScript = `
machine %s
login uu
password uugt
`

func main() {

	// flags
	port := flag.String("port", "8080", "gogs port")
	flag.Parse()

	// submodules repositories
	subModules := []string{"FERepo", "BERepo", "Java", "Prototypes", "Docs", "ProjectSpecs"}
	fmt.Println("--------------\n-------ENVS\n--------------")
	// environment variables
	projectName := getEnv("APP_NAME", "myProject")
	organizationName := getEnv("ORG_NAME", "myOrganization")
	ELBHost := getEnv("ELB_HOSTNAME", "localhost")
	ELBPort := getEnv("ELB_PORT", *port)
	ELBNexusHost := getEnv("ELB_NEXUS_HOST", "your.nexus.url")
	ELBDroneHost := getEnv("ELB_DRONE_HOST", "your.drone.url")
	fullNexusURL := fmt.Sprintf("http://%s/nexus/content/groups/%s/", ELBNexusHost, organizationName)

	// wait to get the service ready
	serviceReady()
	// new gogs organization
	NewGogsOrg(organizationName, "uu")

	// // new gogs repositories
	for _, rep := range subModules {
		NewGogsRepo(rep, organizationName)
	}
	// setup local git user .netrc file
	ioutil.WriteFile("/root/.netrc", []byte(fmt.Sprintf(gitSetupScript, ELBHost)), 0600)
	cmd := NewCMD("cat", "/root/.netrc")
	cmd.Exec()
	// setup local git user
	cmd = NewCMD("git", "config", "--global", "user.email", "uu@uu.com")
	cmd.Exec()
	cmd = NewCMD("git", "config", "--global", "user.name", "uu")
	cmd.Exec()

	type Data struct {
		ProjectName  string
		CompanyName  string
		NexusFullURL string
		NexusURL     string
	}
	fmt.Println("--------------\n-------JAVA\n--------------")
	// java
	os.Chdir("/tmp")
	cmd = NewCMD("git", "clone", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/Java.git")
	cmd.Exec()
	os.Chdir("/tmp/Java")
	cmd = NewCMD("curl", "-LOk", "")
	cmd.Exec()
	cmd = NewCMD("unzip", "-o", "layerx-archetype.zip")
	cmd.Exec()
	os.Remove("layerx-archetype.zip")
	err := ioutil.WriteFile("/tmp/Java/.drone.yml", javaYAML(ELBDroneHost), 0644)
	if err != nil {
		fmt.Println("Error creating .drone.yml", err)
	}
	// Templating files
	files := []string{"Resources/settings.xml", "pom.xml"}
	for _, file := range files {
		t, err := template.ParseFiles(file)
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
		}
		err = t.Execute(f, Data{projectName, organizationName, fullNexusURL, ELBNexusHost})
		if err != nil {
			fmt.Println(err)
		}
	}
	<-time.After(time.Second * 5)
	os.Chdir("/tmp/Java")
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("git", "add", ".")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "commit", "-m", "Added Layerx Archetype")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "master")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "dev", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "dev")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "qa", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "qa")
	cmd.Exec()

	fmt.Println("--------------\n-------SPECS\n--------------")
	// ProjectSpecs
	os.Chdir("/tmp")
	cmd = NewCMD("git", "clone", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/ProjectSpecs.git")
	cmd.Exec()
	os.Chdir("/tmp/ProjectSpecs")
	cmd = NewCMD("curl", "-LOk", "")
	cmd.Exec()
	cmd = NewCMD("unzip", "-o", "project-specs.zip")
	cmd.Exec()
	os.Remove("project-specs.zip")

	err = ioutil.WriteFile("/tmp/ProjectSpecs/.drone.yml", specsYAML(ELBDroneHost), 0644)
	if err != nil {
		fmt.Println("Error creating .drone.yml", err)
	}

	cmd = NewCMD("git", "add", ".")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "commit", "-m", "First commit")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "master")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "dev", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "dev")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "qa", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "qa")
	cmd.Exec()

	fmt.Println("--------------\n-------PROTOTYPES\n--------------")
	// Prototypes
	os.Chdir("/tmp")
	cmd = NewCMD("git", "clone", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/Prototypes.git")
	cmd.Exec()
	os.Chdir("/tmp/Prototypes")
	cmd = NewCMD("curl", "-LOk", "")
	cmd.Exec()
	cmd = NewCMD("unzip", "-o", "prototypes.zip")
	cmd.Exec()
	os.Remove("prototypes.zip")
	cmd = NewCMD("git", "add", ".")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "commit", "-m", "First commit")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "master")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "dev", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "dev")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "qa", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "qa")
	cmd.Exec()

	fmt.Println("--------------\n-------FEREPO\n--------------")
	// FERepo
	os.Chdir("/tmp")
	cmd = NewCMD("git", "clone", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/FERepo.git")
	cmd.Exec()
	os.Chdir("/tmp/FERepo")
	cmd = NewCMD("curl", "-LOk", "")
	cmd.Exec()
	cmd = NewCMD("unzip", "-o", "ferepo.zip")
	cmd.Exec()
	os.Remove("ferepo.zip")
	cmd = NewCMD("git", "add", ".")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "commit", "-m", "First commit")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "master")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "dev", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "dev")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "qa", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "qa")
	cmd.Exec()

	fmt.Println("--------------\n-------BEREPO\n--------------")
	// BERepo
	os.Chdir("/tmp")
	cmd = NewCMD("git", "clone", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/BERepo.git")
	cmd.Exec()
	os.Chdir("/tmp/BERepo")
	cmd = NewCMD("curl", "-LOk", "/berepo.zip")
	cmd.Exec()
	cmd = NewCMD("unzip", "-o", "berepo.zip")
	cmd.Exec()
	os.Remove("berepo.zip")
	cmd = NewCMD("git", "add", ".")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "commit", "-m", "First commit")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "master")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "dev", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "dev")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "qa", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "qa")
	cmd.Exec()

	fmt.Println("--------------\n-------MAIN REPO\n--------------")
	// Main repository
	NewGogsRepo(projectName, organizationName)

	os.Chdir("/tmp")
	cmd = NewCMD("git", "clone", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/"+projectName+".git")
	cmd.Exec()
	os.Chdir("/tmp/" + projectName)

	// adding modules
	for _, rep := range subModules {
		cmd = NewCMD("git", "submodule", "add", "http://"+ELBHost+":"+ELBPort+"/"+organizationName+"/"+rep+".git")
		cmd.Exec()
	}

	// adding initial project stuff
	cmd = NewCMD("curl", "-LOk", "")
	cmd.Exec()
	cmd = NewCMD("unzip", "-o", "project-init.zip")
	cmd.Exec()
	os.Remove("project-init.zip")
	err = ioutil.WriteFile("/tmp/"+projectName+"/.drone.yml", mainYAML(), 0644)
	if err != nil {
		fmt.Println("Error creating .drone.yml", err)
	}

	// Templating files
	files = []string{"default.yaml", "develop.yaml", "dockerHost/host.yaml", "Resources/JCRSyncr.yaml", "Resources/settings.xml"}

	for _, file := range files {
		t, err := template.ParseFiles(file)
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
		}
		err = t.Execute(f, Data{projectName, organizationName, fullNexusURL, ELBNexusHost})
		if err != nil {
			fmt.Println(err)
		}
	}

	cmd = NewCMD("git", "add", ".")
	cmd.Exec()
	cmd = NewCMD("git", "status")
	cmd.Exec()
	cmd = NewCMD("ls", "-ltrah")
	cmd.Exec()
	cmd = NewCMD("git", "commit", "-m", "Added submodules and initial stuff")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "master")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "dev", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "dev")
	cmd.Exec()
	cmd = NewCMD("git", "checkout", "-b", "qa", "master")
	cmd.Exec()
	cmd = NewCMD("git", "push", "origin", "qa")
	cmd.Exec()
}
