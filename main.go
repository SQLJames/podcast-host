package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	logfile, err := os.Create("gin.log")
	if err != nil {
		log.Fatal(err)
	}
	ProgramConfiguration, err := initializeConfig("./Config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	ProgramConfiguration.CreateFeeds()
	//log to standard out and the file at the same time
	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)
	router := gin.New()
	//pass in the logging function format
	router.Use(gin.LoggerWithFormatter(logging))
	router.Use(gin.Recovery())
	router.StaticFS("/images", http.Dir(ProgramConfiguration.Static.Images))
	router.StaticFS("/episodes", http.Dir(ProgramConfiguration.Static.EpisodeLocation))
	router.StaticFile("/favicon.ico", ProgramConfiguration.Static.Favicon)
	files, err := ioutil.ReadDir(ProgramConfiguration.RSS.SearchFolder)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {

		filename := filepath.Join(filepath.Join(ProgramConfiguration.RSS.SearchFolder, f.Name()), "rss.xml")
		router.StaticFile("/"+f.Name(), filename)
	}

	router.Run(":80")
}
