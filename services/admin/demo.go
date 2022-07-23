package admin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/models"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func AppIsInDemoMode() (bool, error) {
	_, err := os.Stat(filepath.Join(os.Getenv("DATA_PATH"), "demo.json"))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func GetDemoMode() (DemoMode, error) {
	file, err := os.Open(filepath.Join(os.Getenv("DATA_PATH"), "demo.json"))
	if err != nil {
		return DemoMode{}, err
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return DemoMode{}, err
	}

	var demoMode DemoMode
	err = json.Unmarshal(fileData, &demoMode)
	if err != nil {
		return DemoMode{}, err
	}

	return demoMode, nil
}

func SetupDemoMode(demoMode DemoMode) error {
	data, err := json.MarshalIndent(demoMode, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(os.Getenv("DATA_PATH"), "demo.json"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func StartDemoMode(tx *gorm.DB) error {
	demoMode, err := GetDemoMode()
	if err != nil {
		return err
	}

	log.Println("DEMO MODE Enabled. The server will reset every: ", demoMode.ResetInterval)

	c := cron.New()
	err = c.AddFunc(demoMode.ResetInterval, func() { ResetServer(tx) })
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
