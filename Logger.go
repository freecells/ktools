/*
 * @Author: Feng
 * @version: v1.0.0
 * @Date: 2020-07-03 13:50:49
 * @LastEditors: Keven
 * @LastEditTime: 2021-09-28 13:11:19
 */
package main

import (
	"keven/ktools/tfile"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogConf struct {
	Channel  string
	MaxDay   int
	Out      string
	SavePath string
}

//ZLog 实例
func ZLog(config LogConf, debug bool) {

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logName := "kos.log"
	// fmt.Println(logName)

	logFiles, err := tfile.GetAllFiles(config.SavePath)

	if err != nil {

		log.Fatal().Err(err).Caller()
	}
	sort.Strings(logFiles)

	if len(logFiles) > config.MaxDay {

		for i := 0; i < len(logFiles)-config.MaxDay; i++ {
			os.Remove(logFiles[i])
		}
	}

	logChannel := config.Channel

	if config.Out == "file" && debug == false {

		if logChannel == "daily" {

			logName = time.Now().Format("2006-01-02") + ".log"

		}

		logFile, _ := os.OpenFile(config.SavePath+logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)

		log.Logger = log.Output(logFile)

	} else {

		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	}

}
