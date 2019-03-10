package dockerdetect

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

type dockerResponse struct {
	Installed               bool
	InstallationURL         string
	InstallationInstruction string
}

// InitAPI creates the docker detection API endpoint
func InitAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/dockerdetect", checkDocker).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func checkDocker(w http.ResponseWriter, r *http.Request) {
	var response dockerResponse

	if dockerIsNotInstalled() {
		if runtime.GOOS == "windows" {
			response = dockerResponse{Installed: false, InstallationURL: "https://hub.docker.com/editions/community/docker-ce-desktop-windows"}

		} else if runtime.GOOS == "darwin" {
			response = dockerResponse{Installed: false, InstallationURL: "https://hub.docker.com/editions/community/docker-ce-desktop-mac"}

		} else {
			response = dockerResponse{Installed: false, InstallationInstruction: "Installation script for Linux started in the background."}
			exec.Command("curl -fsSL https://get.docker.com -o get-docker.sh")
			exec.Command("sudo sh get-docker.sh")

		}
	} else {
		response = dockerResponse{Installed: true, InstallationInstruction: "Docker is already installed."}
	}
	enableCors(&w)
	json.NewEncoder(w).Encode(response)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
