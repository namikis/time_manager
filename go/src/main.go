package main

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func errorResponse(w http.ResponseWriter, errMsg error, statusCode int) {
	if errMsg != nil {
		log.Println(errMsg)
	}
	w.WriteHeader(statusCode)
}

func main() {
	// クライアント
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

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
				case "ping":
					if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("pong", false)); err != nil {
						errorResponse(w, err, http.StatusInternalServerError)
						return
					}
				case "そんなこと言っても":
					if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("しょうがないじゃないかぁ", false)); err != nil {
						errorResponse(w, err, http.StatusInternalServerError)
						return
					}
				}
			}
		}
	})

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
