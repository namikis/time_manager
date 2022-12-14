package main

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io"
	"io/ioutil"
	"log"
	"mymodule/attendance"
	"mymodule/breaking"
	"net/http"
	"os"
	"strings"
	"time"
)

func errorResponse(w http.ResponseWriter, errMsg error, statusCode int) {
	if errMsg != nil {
		log.Println(errMsg)
	}
	w.WriteHeader(statusCode)
}

func sendMessage(api *slack.Client, event *slackevents.AppMentionEvent, w http.ResponseWriter, text string) {
	if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText(text, false)); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
	}
}

func currentTime() string {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "cannot get TZ"
	}
	return time.Now().In(tokyo).Format("2006-01-02 15:04")
}

func main() {
	// クライアント
	var api *slack.Client = slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		// リクエスト検証
		verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
		if err != nil {
			errorResponse(w, err, http.StatusInternalServerError)
			return
		}

		bodyReader := io.TeeReader(r.Body, &verifier)
		body, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			errorResponse(w, err, http.StatusInternalServerError)
			return
		}

		if err := verifier.Ensure(); err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			errorResponse(w, err, http.StatusInternalServerError)
			return
		}

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			var res *slackevents.ChallengeResponse
			if err := json.Unmarshal(body, &res); err != nil {
				errorResponse(w, err, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			if _, err := w.Write([]byte(res.Challenge)); err != nil {
				errorResponse(w, err, http.StatusInternalServerError)
				return
			}
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch event := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				message := strings.Split(event.Text, " ")
				if len(message) < 2 {
					w.WriteHeader(http.StatusBadRequest)
					errorResponse(w, nil, http.StatusBadRequest)
					return
				}

				command := message[1]
				switch command {
				case "test":
					sendMessage(api, event, w, "ok!")
				case "testUserAndTime":
					res_text := "user: " + event.User + " time: " + currentTime()
					sendMessage(api, event, w, res_text)
				case "start":
					current_time := currentTime()
					attendance.InsertRecord(event.User, current_time)
					sendMessage(api, event, w, "start time : "+current_time)
				case "end":
					current_time := currentTime()
					start_time, end_time, working_time, breaking_time := attendance.UpdateRecord(event.User, current_time)
					sendMessage(api, event, w, "start time : "+start_time+"\nend time : "+end_time+"\ntotal working time: "+working_time+"\nbreaking time: "+breaking_time)
				case "break":
					if len(message) < 3 {
						w.WriteHeader(http.StatusBadRequest)
						errorResponse(w, nil, http.StatusBadRequest)
						return
					}
					current_time := currentTime()
					sub_command := message[2]
					response_text := sub_command + " breaking."

					var result int
					if sub_command == "start" {
						result = breaking.InsertBreak(current_time, event.User)
					} else if sub_command == "end" {
						result = breaking.UpdateBreak(current_time, event.User)
					} else {
						w.WriteHeader(http.StatusBadRequest)
						errorResponse(w, nil, http.StatusBadRequest)
						return
					}

					if result == 0 {
						response_text = "no attendance record."
					}
					sendMessage(api, event, w, response_text)
				default:
					sendMessage(api, event, w, "invalid message.")
				}
			}
		}
	})

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
