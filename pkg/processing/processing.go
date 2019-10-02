package processing

import (
	"github.com/Sirupsen/logrus"
	"github.com/papertrail/go-tail/follower"
	"io"
	"nginx-to-ch/config"
	"nginx-to-ch/pkg/clickhouse"
	"nginx-to-ch/pkg/nginx"
)

var (
	t *follower.Follower
	Line *nginx.LogLine
	)


func Reader(conf *config.Config) (*follower.Follower, error) {
	t, err := follower.New(conf.Nginx.LogPath, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})
	return t, err
}


func Read(c <- chan *nginx.LogLine, config *config.Config) {
	batch := make([]nginx.LogLine, 0, config.App.BatchSize)
	for elem := range c {
		batch = append(batch, *elem)
		if len(batch) == config.App.BatchSize {
			logrus.Infof("Collected %v items.", config.App.BatchSize)
			go func() {
				logrus.Info("Connecting to clickhouse.")
				cnx, _ := clickhouse.Connect(config)
				logrus.Infof("Inserting %v items to clickhouse.", config.App.BatchSize)
				clickhouse.Insert(cnx, batch, config)
				batch = nil
			}()
		}
	}
}
