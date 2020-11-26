package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Rec struct {
	Name   string
	Week   string
	Time   string
	Length int
}

func (r *Rec) IsValid(t time.Time) bool {
	week := t.Weekday().String()
	nowTime := t.Format("15:04")
	return r.Week == week && r.Time == nowTime
}

func NewRecs(path string) ([]*Rec, error) {
	yml, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open rec.yaml %v", err)
	}
	defer func() {
		_ = yml.Close()
	}()

	bs, err := ioutil.ReadAll(yml)
	if err != nil {
		return nil, fmt.Errorf("failed to read rec.yaml %v", err)
	}

	var recs []*Rec
	if err := yaml.Unmarshal(bs, &recs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rec config %v", err)
	}

	return recs, nil
}
