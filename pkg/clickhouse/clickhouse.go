package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/kshvakov/clickhouse"
	"nginx-to-ch/config"
	"nginx-to-ch/pkg/nginx"
	"strconv"
	"strings"
)

func Connect(c *config.Config) (*sql.DB, error) {
	dataSource := fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s",
		c.Clickhouse.Connection.Host,
		strconv.Itoa(c.Clickhouse.Connection.Port),
		c.Clickhouse.Credentials.User,
		c.Clickhouse.Credentials.Password,
		c.Clickhouse.Connection.Db)
	cnx, err := sql.Open("clickhouse", dataSource)
	if err != nil {
		logrus.Fatalf("Clickhouse connection error: %s", err.Error())
	}
	return cnx, err
}

func Callback(c *sql.DB) error {
	defer c.Close()
	if err := c.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			logrus.Fatalf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			logrus.Fatal(err)
		}
		return err
	}
	return nil
}

func Prepare(c *sql.DB, conf *config.Config) error {
	defer c.Close()
	tableName := conf.Clickhouse.Connection.Table
	logrus.Infof("Trying to create table %s if not exists", tableName)
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			TimeLocal DateTime,
			Date Date DEFAULT toDate(TimeLocal),
			RemoteAddr String,
			RequestMethod String,
			Request String,
			RequestBody String,
			BytesSent Int64,
			Status Int32,
			HttpHost String,
			HttpReferrer String,
			HttpUserAgent String,
			RequestTime Float32
		) ENGINE = MergeTree(Date, (Status, RemoteAddr, Date), 8192)
	`, tableName)
	_, err := c.Exec(query)
	if err != nil {
		logrus.Fatalf("Clickhouse: create table error: %s", err.Error())
	}
	return err
}

func Insert(c *sql.DB, lines []nginx.LogLine, conf *config.Config) () {
	defer c.Close()
	tx, _ := c.Begin()
	query := fmt.Sprintf(`INSERT INTO %s
					(TimeLocal,
					RemoteAddr,
					RequestMethod,
					Request,
					RequestBody,
					BytesSent,
					Status,
					HttpHost,
					HttpReferrer,
					HttpUserAgent,
					RequestTime) VALUES   (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					conf.Clickhouse.Connection.Table)
	stmt, _ := tx.Prepare(query)
	defer stmt.Close()


	for _, l := range lines {

		fmtTime, _ := l.ConvertTime()

		if _, err := stmt.Exec(
			fmtTime,
			l.RemoteAddress,
			strings.Split(l.Request, " ")[0],
			strings.Split(l.Request, " ")[1],
			l.RequestBody,
			l.BodyBytesSend,
			l.Status,
			l.HttpHost,
			l.HttpReferrer,
			l.HttpUserAgent,
			l.RequestTime,
		); err != nil {
			logrus.Warnf("Clickhouse prepared statement execution error: %s", err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.Warn("Clickhouse commit error: %s", err.Error())
	} else {
		logrus.Infof("Clickhouse: commited %v items", len(lines))
	}
}
