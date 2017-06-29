package lib

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func forwardLeaderPut(addr string, path string, tsv *TimeStampValue) error {
	url := fmt.Sprintf("http://%s%s", addr, path)
	return forwardPut(url, tsv.data)
}

func forwardFollowerPut(addr string, path string, tsv *TimeStampValue) error {
	url := fmt.Sprintf("http://%s%s/%d", addr, path, tsv.timestamp.UnixNano())
	return forwardPut(url, tsv.data)
}

func forwardPut(url string, b []byte) (err error) {
	log.Printf("Forward put %s", url)

	contentType := ""
	body := bytes.NewBuffer(b)
	resp, err := http.Post(url, contentType, body)

	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Error forwarding %s: %d", url, resp.StatusCode)
	}

	return
}

func forwardLeaderDelete(addr string, path string) error {
	url := fmt.Sprintf("http://%s%s", addr, path)
	return forwardDelete(url)
}

func forwardFollowerDelete(addr string, path string, tsv *TimeStampValue) error {
	url := fmt.Sprintf("http://%s%s/%d", addr, path, tsv.timestamp.UnixNano())
	return forwardDelete(url)
}

func forwardDelete(url string) (err error) {
	log.Printf("Forward delete %s", url)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Error forwarding %s: %d", url, resp.StatusCode)
	}

	return
}
