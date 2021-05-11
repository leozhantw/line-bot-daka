package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

func (s *Server) HandleCallBack(w http.ResponseWriter, req *http.Request) {
	events, err := s.line.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch {
				case strings.EqualFold(message.Text, "到"):
					var (
						content string
						err     error
					)
					now := time.Now().In(s.location)

					record, err := s.record.GetByUser(event.Source.UserID, now)
					if err != nil {
						if !errors.Is(err, gorm.ErrRecordNotFound) {
							log.Fatalln("failed to get record", err)
						}

						if err := s.record.Create(&dao.Record{UserID: event.Source.UserID, WorkDate: now}); err != nil {
							log.Fatalln("failed to create record", err)
						}

						content = fmt.Sprintf("%s 打卡上班", now.In(s.location).Format("15:04"))
					} else {
						content = fmt.Sprintf("今日已於 %s 打卡上班", record.WorkedAt.In(s.location).Format("15:04"))
					}

					if _, err = s.line.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
						log.Fatalln("failed to reply message", err)
					}
				}
			}
		}
	}
}
