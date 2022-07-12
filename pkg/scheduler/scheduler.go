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
	workingHours      = 9
	scheduleFrequency = 10
)

type Scheduler struct {
	record    dao.RecordDAO
	line      *linebot.Client
	patientID string
}

var (
	queriesForGetOffWork = []string{
		"快逃啊",
		"大家可以回家啦",
		"下班啦",
		"不要浪費生命了",
		"下班表情包",
		"下班 梗圖",
	}
	queriesForTakeMedicine = []string{
		"該吃藥了",
		"吃藥 梗圖",
	}
)

func New(record dao.RecordDAO, line *linebot.Client, patientID string) *Scheduler {
	return &Scheduler{
		record:    record,
		line:      line,
		patientID: patientID,
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

		if record.UserID == s.patientID {
			takeMedicineTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
			if now.After(takeMedicineTime) {
				continue
			}

			d := takeMedicineTime.Sub(now)
			if d > time.Minute*time.Duration(scheduleFrequency) {
				continue
			}

			eg.Go(func() error {
				return s.countdownForTakeMedicine(record, d)
			})
		}

		offWorkTime := record.WorkedAt.Add(time.Hour * time.Duration(workingHours))
		if now.After(offWorkTime) {
			continue
		}

		d := offWorkTime.Sub(now)
		if d > time.Minute*time.Duration(scheduleFrequency) {
			continue
		}

		eg.Go(func() error {
			return s.countdownForGetOffWork(record, d)
		})
	}

	return eg.Wait()
}

func (s *Scheduler) countdownForGetOffWork(record dao.Record, d time.Duration) error {
	log.Println(fmt.Sprintf("user %s starting to countdown [%s]", record.UserID, d.String()))

	timer := time.NewTimer(d)
	<-timer.C

	pic, err := randompicture.Random(queriesForGetOffWork)
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

func (s *Scheduler) countdownForTakeMedicine(record dao.Record, d time.Duration) error {
	log.Println(fmt.Sprintf("user %s starting to countdown [%s]", record.UserID, d.String()))

	timer := time.NewTimer(d)
	<-timer.C

	pic, err := randompicture.Random(queriesForTakeMedicine)
	if err != nil {
		return fmt.Errorf("failed to random picture %v", err)
	}

	if _, err = s.line.PushMessage(
		record.UserID,
		linebot.NewTextMessage("兄弟該吃藥啦"),
		linebot.NewImageMessage(pic, pic),
	).Do(); err != nil {
		return fmt.Errorf("failed to push message %v", err)
	}

	log.Println(fmt.Sprintf("pushed the message to user %s", record.UserID))

	return nil
}
