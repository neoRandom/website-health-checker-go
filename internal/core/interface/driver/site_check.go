package driver

import "http-server/internal/core/model"

type CheckSites func(targets []model.Site) error
