package main

import (
	"fmt"
	"os"
	"os/exec"
	"picture-frame/frame"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const configName string = ".frame.yaml"
const urlKey string = "PHOTOS_URL"

var browserURL string

// Logger is the main logger object
// var Logger *log.Logger

func init() {
	gopath := os.Getenv("GOPATH")
	os.Chdir(fmt.Sprintf("%s/src/picture-frame", gopath))
}

// func initLogs() {
// 	timeFormat := "2006-01-02_15:04:05"
// 	logFile := fmt.Sprintf("frame-%s.log", time.Now().Format(timeFormat))

// 	var file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		fmt.Println("Could Not Open Log File : " + err.Error())
// 	}

// 	// create logger
// 	Logger = &log.Logger{
// 		Out:       file,
// 		Formatter: new(log.TextFormatter),
// 	}
// }

func updateRepo() {
	cmd := exec.Command("git", "pull")
	if err := cmd.Run(); err != nil {
		log.Warnf("error updating project repo:\n\t%s\n", err)
		return
	}
	log.Info("updated project repo")
}

func loadURL() {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("no config file found at:\n\t%s", configName)
		} else {
			log.Fatalf("error reading config file:\n\t%s", err)
		}
	}

	browserURL = viper.GetString(urlKey)
	log.Printf("got browser URL from config:\n\t%s", browserURL)
}

func main() {
	// initLogs()
	updateRepo()
	log.Info("running main")
	loadURL()

	browser := &frame.Browser{URL: browserURL}
	browser.SetDefaults()

	browser.RunForever()
}
