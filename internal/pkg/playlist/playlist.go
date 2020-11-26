package playlist

import (
	"fmt"
	"github.com/grafov/m3u8"
	"net/http"
)

type PlayList struct {
	pl *m3u8.MasterPlaylist
}

func NewPlayList(url string) (*PlayList, error) {
	f, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download from %v: %v", url, err)
	}
	defer func() { _ = f.Body.Close() }()

	p, _, err := m3u8.DecodeFrom(f.Body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %v", err)
	}

	pl, ok := p.(*m3u8.MasterPlaylist)
	if !ok {
		return nil, fmt.Errorf("url is not master playlist %v", err)
	}

	return &PlayList{pl: pl}, nil
}

func (p *PlayList) URI() (string, error) {
	for _, variant := range p.pl.Variants {
		return variant.URI, nil
	}

	return "", fmt.Errorf("cannot find enable URI")
}
