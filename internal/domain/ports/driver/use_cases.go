package driver

import "http-server/internal/domain/models"

type AddSite func(url string) (*models.Site, error)
type GetSiteList func() ([]*models.Site, error)
