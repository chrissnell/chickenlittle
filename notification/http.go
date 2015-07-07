package notification

import (
	"fmt"
	"net/http"
	"net/url"
)

// CallHTTP will execute an POST request against the provided URL.
func (e *Engine) CallHTTP(u *url.URL, subject string, message string, uuid string) error {
	resp, err := http.PostForm(u.String(), url.Values{
		"subject": []string{subject},
		"message": []string{message},
		"uuid":    []string{uuid},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP Request failed: %s", resp.Status)
	}
	return nil
}
