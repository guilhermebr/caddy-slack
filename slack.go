package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

type Slack struct {
	Next      httpserver.Handler
	ErrorFunc func(http.ResponseWriter, *http.Request, int) // failover error handler
	Rules     []Rule
}

type Rule struct {
	Token string
}

// ServerHTTP is the HTTP handler for this middleware
func (s Slack) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	for _, rule := range s.Rules {
		// Record the response
		responseRecorder := httpserver.NewResponseRecorder(w)

		// Attach the Replacer we'll use so that other middlewares can
		// set their own placeholders if they want to.
		rep := httpserver.NewReplacer(r, responseRecorder, CommonLogEmptyValue)
		responseRecorder.Replacer = rep

		// Bon voyage, request!
		status, err := s.Next.ServeHTTP(responseRecorder, r)

		if status >= 400 {
			// There was an error up the chain, but no response has been written yet.
			// The error must be handled here so the log entry will record the response size.
			if s.ErrorFunc != nil {
				s.ErrorFunc(responseRecorder, r, status)
			} else {
				// Default failover error handler
				responseRecorder.WriteHeader(status)
				fmt.Fprintf(responseRecorder, "%d %s", status, http.StatusText(status))
			}
			status = 0
		}

		// Write log entry
		s.Log(rule.Token, rep.Replace(CommonLogFormat))

		return status, err
	}
	return s.Next.ServeHTTP(w, r)
}

func (s Slack) Log(token, msg string) error {
	data, _ := json.Marshal(map[string]string{"text": msg})
	client := &http.Client{}
	resp, err := client.Post("https://hooks.slack.com/services/"+token, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
	return nil
}

const (
	CommonLogEmptyValue = "-"
	CommonLogFormat     = `{remote} ` + CommonLogEmptyValue + ` [{when}] "{method} {uri} {proto}" {status} {size}`
)
