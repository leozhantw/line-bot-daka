package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/leozhantw/line-bot-daka/pkg/randompicture"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/sync/errgroup"
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
	var eg errgroup.Group

	now := time.Now()
	records, err := s.record.ListByDate(now)
	if err != nil {
		return fmt.Errorf("failed to list by date %v", err)
	}

	for _, record := range *records {
		record := record
		offWorkTime := record.WorkedAt.Add(time.Hour * time.Duration(workingHours))

		if now.After(offWorkTime) {
			continue
		}

		d := offWorkTime.Sub(now)
		if d > time.Minute*time.Duration(scheduleFrequency) {
			continue
		}

		eg.Go(func() error {
			return s.countdown(record, d)
		})
	}

	return eg.Wait()
}

func (s *Scheduler) countdown(record dao.Record, d time.Duration) error {
	log.Println(fmt.Sprintf("user %s starting to countdown [%s]", record.UserID, d.String()))

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

	log.Println(fmt.Sprintf("pushed the message to user %s", record.UserID))

	return nil
}
