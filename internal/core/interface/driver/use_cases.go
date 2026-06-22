package driver

import "http-server/internal/core/model"

type GetSiteList func() ([]*models.Site, error)
type AddSite func(url string) (*models.Site, error)
type UpdateSite func(id models.SiteID, url string) error
type RemoveSite func(id models.SiteID) error
