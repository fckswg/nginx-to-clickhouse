package main

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"nginx-to-ch/config"
	"nginx-to-ch/pkg/clickhouse"
	"nginx-to-ch/pkg/nginx"
	"nginx-to-ch/pkg/processing"
	"strings"
)

var warn int
var auth *config.AuthRequest
var reg *config.RegRequest
var restore *config.RestoreRequest

func main() {

	c := make(chan *nginx.LogLine, 50000)

	conf := config.Read()
	if err := config.Validate(conf); err != nil {
		logrus.Fatalf("Config validation failed: %s\n", err.Error())
	}

	cnx, _ := clickhouse.Connect(conf)
	clickhouse.Prepare(cnx, conf)

	t, err := processing.Reader(conf)
	if err != nil {
		logrus.Fatalf("Cant start follow nginx config: %s", err.Error())
	}

	// read channel with logline structs and write chunked records to clickhouse
	go processing.Read(c, conf)

	for l := range t.Lines() {
		err := json.Unmarshal(l.Bytes(), &processing.Line)
		if err != nil {
			logrus.Warnf("Nginx log line unmarshal error: %s", err.Error())
		}
		// validate logline struct and reinit fields if empty or other reason
		err = processing.Line.Validate()
		if err != nil {
			if processing.Line.HttpReferrer == "" {
				processing.Line.HttpReferrer = "empty"
			}
			if processing.Line.HttpHost == "" {
				processing.Line.HttpHost = "empty"
			}
			if processing.Line.RequestBody == "" {
				processing.Line.RequestBody = "empty"
			}
			if strings.Contains(processing.Line.RequestBody, "Content-Type: image") {
				processing.Line.RequestBody = "image upload"
			}

			warn += 1
			if warn >= conf.Nginx.WarnCount {
				logrus.Warnf("Nginx log line validation failed: %s", err.Error())
				logrus.Warnf("Nginx log line validation failed: More than %v warnings occurred.", warn)
				warn = 0
			}
		}
		c <- processing.Line
	}
}
