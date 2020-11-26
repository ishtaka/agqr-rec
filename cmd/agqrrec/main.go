package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ishtaka/agqr-rec/internal/pkg/recorder"
	"github.com/ishtaka/agqr-rec/pkg/cron"
)

const SaveDir = "rec"
const Location = "Asia/Tokyo"

func main() {
	l, err := time.LoadLocation(Location)
	if err != nil {
		log.Fatal(err)
	}

	r, err := recorder.NewRecorder(SaveDir, l)
	if err != nil {
		log.Fatal(err)
	}

	c, err := cron.New(l)
	if err != nil {
		log.Fatal(err)
	}

	err = c.AddFunc("29,59 * * * *", func() {
		time.Sleep(53 * time.Second)
		if err := r.Start(); err != nil {
			log.Fatal(err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	waitCh := func() <-chan bool {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)

		boolCh := make(chan bool)
		go func(sig <-chan os.Signal) {
			defer close(boolCh)
			<-sig
			boolCh <- true
		}(sig)

		return boolCh
	}

	c.Start(waitCh())
}
