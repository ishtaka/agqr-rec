package recorder

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ishtaka/agqr-rec/internal/pkg/config"
	"github.com/ishtaka/agqr-rec/internal/pkg/playlist"
)

const hlsURL = "https://www.uniqueradio.jp/agplayer5/hls/mbr-0-cdn.m3u8"
const format = "20060102"
const ext = ".mp4"

type Recorder struct {
	recs     []*config.Rec
	saveDir  string
	location *time.Location
}

func NewRecorder(saveDir string, l *time.Location) (*Recorder, error) {
	recs, err := config.NewRecs("configs/rec.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create recorder %v", err)
	}

	return &Recorder{
		recs:     recs,
		saveDir:  saveDir,
		location: l,
	}, nil
}

func (r *Recorder) Start() error {
	now := time.Now().In(r.location).Add(1 * time.Minute)

	var recs []*config.Rec
	for _, rec := range r.recs {
		if rec.IsValid(now) {
			recs = append(recs, rec)
		}
	}

	if len(recs) == 0 {
		return nil
	}

	pl, err := playlist.NewPlayList(hlsURL)
	if err != nil {
		return fmt.Errorf("failed to get master playlist %v", err)
	}

	uri, err := pl.URI()
	if err != nil {
		return fmt.Errorf("failed to get media URI %v", err)
	}

	for _, rec := range recs {
		err := r.Rec(rec, uri, now)
		if err != nil {
			return fmt.Errorf("failed to rec:%v", err)
		}
	}

	return nil
}

func (r *Recorder) Rec(rec *config.Rec, uri string, t time.Time) error {
	path := fmt.Sprintf("%s/%s", r.saveDir, rec.Name)
	if ok := makeDirIfNotExists(path); !ok {
		return nil
	}

	args := []string{
		"-i", uri,
		"-t", strconv.Itoa(rec.Length * 60),
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		fmt.Sprintf("%s/%s%s", path, t.Format(format), ext),
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		command := fmt.Sprintf("ffmpeg %s", args)
		return fmt.Errorf("failed to execute command %v\nout: %v\n, err: %v\n, %v", command, out.String(), stderr.String(), err)
	}

	return nil
}

func makeDirIfNotExists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		return true
	}

	return false
}
