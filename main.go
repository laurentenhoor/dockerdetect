package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// DockerResponse is an API response
type DockerResponse struct {
	Installed               bool
	InstallationURL         string
	InstallationInstruction string
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/docker", CheckDocker).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

// CheckDocker will return whether Docker is installed, furthermore it will give installation relevant installation instructions
func CheckDocker(w http.ResponseWriter, r *http.Request) {

	if dockerIsNotInstalled() {
		if runtime.GOOS == "windows" {
			var response = DockerResponse{Installed: false, InstallationURL: "https://hub.docker.com/editions/community/docker-ce-desktop-windows"}
			json.NewEncoder(w).Encode(response)

		} else if runtime.GOOS == "darwin" {
			var response = DockerResponse{Installed: false, InstallationURL: "https://hub.docker.com/editions/community/docker-ce-desktop-mac"}
			json.NewEncoder(w).Encode(response)

		} else {
			var response = DockerResponse{Installed: false, InstallationInstruction: "Installation script for Linux started in the background."}
			json.NewEncoder(w).Encode(response)
			exec.Command("curl -fsSL https://get.docker.com -o get-docker.sh")
			exec.Command("sudo sh get-docker.sh")

		}
	} else {
		var response = DockerResponse{Installed: true, InstallationInstruction: "Docker is already installed."}
		json.NewEncoder(w).Encode(response)
	}
}

func dockerIsNotInstalled() bool {
	out, err := exec.Command("docker", "version").Output()

	if out != nil && err == nil {
		fmt.Println("Docker is installed.")
		return false
	}
	fmt.Println("Docker is not installed.")
	return true
}
