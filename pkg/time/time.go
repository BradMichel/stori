package time

import "time"

type Time = time.Time

var UTC = time.UTC

var Date = time.Date

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Now() Time {
	return time.Now()
}
