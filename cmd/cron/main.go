package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/leozhantw/line-bot-daka/pkg/scheduler"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Args struct {
	ChannelSecret string `long:"channel-secret" env:"CHANNEL_SECRET" required:"true"`
	ChannelToken  string `long:"channel-token" env:"CHANNEL_TOKEN" required:"true"`
	PGHost        string `long:"pg-host" env:"PG_HOST" required:"true"`
	PGUser        string `long:"pg-user" env:"PG_USER" required:"true"`
	PGPassword    string `long:"pg-password" env:"PG_PASSWORD" required:"true"`
	PGName        string `long:"pg-name" env:"PG_NAME" required:"true"`
	PGPort        string `long:"pg-port" env:"PG_PORT" required:"true"`
}

func main() {
	var args Args
	if _, err := flags.NewParser(&args, flags.Default).Parse(); err != nil {
		log.Fatal(err)
	}

	line, err := linebot.New(args.ChannelSecret, args.ChannelToken)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", args.PGHost, args.PGUser, args.PGPassword, args.PGName, args.PGPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	recordDAO := dao.NewPGRecordDAO(db)

	s := scheduler.New(
		recordDAO,
		line,
	)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
