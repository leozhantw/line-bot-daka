package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/leozhantw/line-bot-daka/pkg/server"
	"github.com/leozhantw/line-bot-daka/pkg/utils/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Args struct {
	Timezone      string      `long:"timezone" env:"TIMEZONE" default:"Asia/Taipei"`
	ChannelSecret string      `long:"channel-secret" env:"CHANNEL_SECRET" required:"true"`
	ChannelToken  string      `long:"channel-token" env:"CHANNEL_TOKEN" required:"true"`
	Port          string      `long:"port" env:"PORT" required:"true"`
	GormConfig    gorm.Config `group:"gorm" namespace:"gorm" env-namespace:"gorm"`
}

func main() {
	var args Args
	if _, err := flags.NewParser(&args, flags.Default).Parse(); err != nil {
		log.Fatalln("failed to new parser", err)
	}

	loc, err := time.LoadLocation(args.Timezone)
	if err != nil {
		log.Fatalln("failed to load location", err)
	}

	line, err := linebot.New(args.ChannelSecret, args.ChannelToken)
	if err != nil {
		log.Fatalln("failed to new line bot", err)
	}

	db, err := gorm.New(args.GormConfig)
	if err != nil {
		log.Fatalln("failed to new gorm", err)
	}

	recordDAO := dao.NewPGRecordDAO(db)

	s := server.New(
		loc,
		line,
		recordDAO,
	)
	s.Routes()

	if err := http.ListenAndServe(":"+args.Port, nil); err != nil {
		log.Fatalln("failed to listen and serve", err)
	}
}
