package scheduler

import (
	"errors"
	"time"

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

		return err
	}

	now := time.Now()
	offWorkTime := record.WorkedAt.Add(time.Hour * time.Duration(workingHours))

	if now.After(offWorkTime) {
		return nil
	}

	d := offWorkTime.Sub(now)
	if d > time.Minute*time.Duration((scheduleFrequency)) {
		return nil
	}

	timer := time.NewTimer(d)
	<-timer.C

	content := "下班啦！ 再不走就虧大啦！！"
	if _, err = s.line.PushMessage(record.UserID, linebot.NewTextMessage(content)).Do(); err != nil {
		return err
	}

	return nil
}
