package crontab

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Crontab struct {
	Rule string
	Fun  func(c context.Context)
}

var (
	crontabs = make([]Crontab, 0)
)

func InitCrontab(c *cron.Cron, ctx context.Context) {
	for _, v := range crontabs {
		go func(v Crontab) {
			c.AddFunc(v.Rule, func() {
				v.Fun(ctx)
			})
		}(v)
	}
}
