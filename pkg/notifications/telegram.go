package notifications

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Alert struct {
	Text string
	Token string
	ChatId string
}

func (a *Alert) Send() (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage",
		a.Token)
	msg := a.Text
	c := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	query := req.URL.Query()
	query.Add("chat_id", a.ChatId)
	query.Add("text", msg)
	query.Add("parse_mode", "html")
	req.URL.RawQuery = query.Encode()

	resp, err := c.Do(req)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return string(respBody), nil

}