package admin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/robfig/cron"
)

func AppIsInDemoMode() (bool, IServiceError) {
	_, err := os.Stat(filepath.Join(os.Getenv("DATA_PATH"), "demo.json"))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, NewServiceError(http.StatusInternalServerError, err)
}

func GetDemoMode() (DemoMode, IServiceError) {
	file, err := os.Open(filepath.Join(os.Getenv("DATA_PATH"), "demo.json"))
	if err != nil {
		return DemoMode{}, NewServiceError(http.StatusInternalServerError, err)
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return DemoMode{}, NewServiceError(http.StatusInternalServerError, err)
	}

	var demoMode DemoMode
	err = json.Unmarshal(fileData, &demoMode)
	if err != nil {
		return DemoMode{}, NewServiceError(http.StatusInternalServerError, err)
	}

	return demoMode, nil
}

func SetupDemoMode(demoMode DemoMode) IServiceError {
	data, err := json.MarshalIndent(demoMode, "", "\t")
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}

	err = os.WriteFile(filepath.Join(os.Getenv("DATA_PATH"), "demo.json"), data, 0644)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}

	return nil
}

func StartDemoMode(db *Database) IServiceError {
	demoMode, serviceError := GetDemoMode()
	if serviceError != nil {
		return serviceError
	}

	log.Println("DEMO MODE Enabled. The server will reset every: ", demoMode.ResetInterval)

	c := cron.New()
	err := c.AddFunc(demoMode.ResetInterval, func() { ResetServer(db) })
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}
	c.Start()
	return nil
}
