package driver

import "http-server/internal/core/model"

type GetLatestResults func() ([]model.Result, error)
type GetSiteLatestResults func(url string) (model.Result, error)
