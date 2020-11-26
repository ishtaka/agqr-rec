package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	c *cron.Cron
}

func (c *Cron) Start(sig <-chan bool) {
	c.c.Start()
	defer c.c.Stop()

	<-sig
}

func (c *Cron) AddFunc(spec string, cmd func()) error {
	_, err := c.c.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("failed to add func %v", err)
	}

	return nil
}

func New(l *time.Location) (*Cron, error) {
	c := cron.New(cron.WithLocation(l))

	return &Cron{c: c}, nil
}
