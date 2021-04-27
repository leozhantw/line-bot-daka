package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Args struct {
	ChannelSecret string `long:"channel-secret" env:"CHANNEL_SECRET" required:"true"`
	ChannelToken  string `long:"channel-token" env:"CHANNEL_TOKEN" required:"true"`
	Port          string `long:"port" env:"PORT" required:"true"`
	PGHost        string `long:"pg-host" env:"PG_HOST" required:"true"`
	PGUser        string `long:"pg-user" env:"PG_USER" required:"true"`
	PGPassword    string `long:"pg-password" env:"PG_PASSWORD" required:"true"`
	PGName        string `long:"pg-name" env:"PG_NAME" required:"true"`
	PGPort        string `long:"pg-port" env:"PG_PORT" required:"true"`
}

type Server struct {
	bot    *linebot.Client
	record dao.RecordDAO
}

func NewServer(args Args) (*Server, error) {
	bot, err := linebot.New(args.ChannelSecret, args.ChannelToken)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", args.PGHost, args.PGUser, args.PGPassword, args.PGName, args.PGPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})
	if err != nil {
		log.Fatalln(err)
	}

	recordDAO := dao.NewPGRecordDAO(db)

	return &Server{
		bot:    bot,
		record: recordDAO,
	}, nil
}

func main() {
	var args Args
	if _, err := flags.NewParser(&args, flags.Default).Parse(); err != nil {
		log.Fatalln(err)
	}

	server, err := NewServer(args)
	if err != nil {
		log.Fatalln(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := server.bot.ParseRequest(req)
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
						now := time.Now()

						record, err := server.record.GetByDate(now)
						if err != nil {
							if errors.Is(err, gorm.ErrRecordNotFound) {
								if content, err = server.daka(now); err != nil {
									log.Println(err)
								}
							}

							log.Println(err)
						} else {
							content = fmt.Sprintf("已於 %s 打卡上班", record.WorkedAt.Format("15:04"))
						}

						if _, err = server.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+args.Port, nil); err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) daka(time time.Time) (string, error) {
	if err := s.record.Create(&dao.Record{WorkDate: time}); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s 打卡上班", time.Format("15:04")), nil
}