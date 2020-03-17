package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type HttpEndpointNotifier struct {
	Enabled     bool
	ClusterName string `json:"cluster-name"`
	BaseURL     string `json:"base-url"`
	Endpoint    string `json:"endpoint"`
	Text        string `json:"text"`
}

// NotifierName provides name for notifier selection
func (notifier *HttpEndpointNotifier) NotifierName() string {
	return "http-endpoint"
}

func (notifier *HttpEndpointNotifier) Copy() Notifier {
	n := *notifier
	return &n
}

//Notify sends messages to the endpoint notifier
func (notifier *HttpEndpointNotifier) Notify(messages Messages) bool {

	for _, message := range messages {
		text := fmt.Sprintf("%s:%s:%s", message.Node, message.Service, message.Status)
		notifier.Text = text
		data, err := json.Marshal(notifier)
		if err != nil {
		log.Println("Unable to encode POST data")
		return false
	}
		b := bytes.NewBuffer(data)
		endpoint := fmt.Sprintf("%s%s", notifier.BaseURL, notifier.Endpoint)
		if res, err := http.Post(endpoint, "application/json", b); err != nil {
			log.Println("Unable to send data to endpoint:", err)
			return false
		} else {
			defer res.Body.Close()
			statusCode := res.StatusCode
			if statusCode != 200 {
				body, _ := ioutil.ReadAll(res.Body)
				log.Println("Unable to notify slack:", string(body))
				return false
			} else {
				log.Println("Slack notification sent.")
				return true
			}
		}
	}
	return true
}
