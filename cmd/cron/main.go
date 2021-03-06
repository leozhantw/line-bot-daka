package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/leozhantw/line-bot-daka/pkg/scheduler"
	"github.com/leozhantw/line-bot-daka/pkg/utils/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Args struct {
	ChannelSecret string `long:"channel-secret" env:"CHANNEL_SECRET"`
	ChannelToken  string `long:"channel-token" env:"CHANNEL_TOKEN"`
	PatientID     string `long:"patient-id" env:"PATIENT_ID"`
	GormConfig    gorm.Config
}

func main() {
	var args Args
	if _, err := flags.NewParser(&args, flags.Default).Parse(); err != nil {
		log.Fatalln("failed to new parser", err)
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

	s := scheduler.New(
		recordDAO,
		line,
		args.PatientID,
	)

	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
