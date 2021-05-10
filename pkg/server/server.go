package server

import (
	"net/http"
	"time"

	"github.com/leozhantw/line-bot-daka/pkg/dao"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Server struct {
	location *time.Location
	line     *linebot.Client
	record   dao.RecordDAO
}

func New(location *time.Location, line *linebot.Client, record dao.RecordDAO) *Server {
	return &Server{
		location: location,
		line:     line,
		record:   record,
	}
}

func (s *Server) Routes() {
	http.HandleFunc("/callback", s.HandleCallBack)
}
