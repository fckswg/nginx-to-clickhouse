package nginx

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

var (
	NginxTimeFormat = "02/Jan/2006:15:04:05 -0700"
	ClickhouseTimeFormat = "2006-01-02 15:04:05"
	validate *validator.Validate
)

type LogLine struct {
	TimeLocal string `json:"time_local" validate:"required"`
	RemoteAddress string `json:"remote_addr" validate:"required"`
	HttpHost string `json:"http_host" validate:"required"`
	Request string `json:"request" validate:"required"`
	Status int `json:"status" validate:"required"`
	HttpReferrer string `json:"http_referrer" validate:"required"`
	HttpUserAgent string `json:"http_user_agent" validate:"required"`
	BodyBytesSend int `json:"body_bytes_sent"`
	RequestBody string `json:"request_body" validate:"required"`
	RequestTime float32 `json:"request_time"`
}

func (l *LogLine) Validate() error {
	validate := validator.New()
	err := validate.Struct(l)
	return err
}

func (l *LogLine) ConvertTime() (string, error) {
	fmtTime, err := time.Parse(NginxTimeFormat, l.TimeLocal)
	if err != nil {
		logrus.Warnf("Nginx log line time convertation error: %s", err.Error())
		return l.TimeLocal, err
	}
	return fmtTime.Format(ClickhouseTimeFormat), err
}


