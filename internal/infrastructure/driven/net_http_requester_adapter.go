package driven

import (
	"fmt"
	"http-server/internal/core/model"
	"net/http"
	"time"
)

type NetHttpRequesterAdapter struct {}

func NewNetHttpRequesterAdapter() *NetHttpRequesterAdapter {
	return &NetHttpRequesterAdapter{}
}

func (hr *NetHttpRequesterAdapter) CheckSite(s *model.Site) (*model.Result, error) {
	if s.Url == "" {
		return nil, fmt.Errorf("url for site '%v' required, none provided", s.Id)
	}

	req, err := http.NewRequest("HEAD", s.Url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	checkedAt := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &model.Result{
		Id: 0,
		SiteId: s.Id,
		StatusCode: resp.StatusCode,
		CheckedAt: checkedAt,
	}
	
	return res, nil
}
