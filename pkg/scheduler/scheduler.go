package scheduler

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/leozhantw/line-bot-daka/pkg/randompicture"

	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

const (
	workingHours      = 8
	scheduleFrequency = 10
)

type Scheduler struct {
	record dao.RecordDAO
	line   *linebot.Client
}

func New(record dao.RecordDAO, line *linebot.Client) *Scheduler {
	return &Scheduler{
		record: record,
		line:   line,
	}
}

func (s *Scheduler) Run() error {
	record, err := s.record.GetByDate(time.Now())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("failed to get by date %v", err)
	}

	now := time.Now()
	offWorkTime := record.WorkedAt.Add(time.Hour * time.Duration(workingHours))

	if now.After(offWorkTime) {
		return nil
	}

	d := offWorkTime.Sub(now)
	if d > time.Minute*time.Duration(scheduleFrequency) {
		return nil
	}

	log.Println(fmt.Sprintf("starting to count down [%s]", d.String()))

	timer := time.NewTimer(d)
	<-timer.C

	pic, err := randompicture.Random()
	if err != nil {
		return fmt.Errorf("failed to random picture %v", err)
	}

	if _, err = s.line.PushMessage(
		record.UserID,
		linebot.NewTextMessage("下班啦！ 再不走就虧大啦！！"),
		linebot.NewImageMessage(pic, pic),
	).Do(); err != nil {
		return fmt.Errorf("failed to push message %v", err)
	}

	return nil
}
