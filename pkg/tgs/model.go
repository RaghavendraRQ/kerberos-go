package tgs

import "time"

type Response struct {
	SessionKey []byte
	ServerId   string
	TimeStamp  time.Time
	LifeTime   time.Duration
}
