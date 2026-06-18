package driver

import "http-server/internal/domain/models"

type GetSiteList func() ([]*models.Site, error)
type AddSite func(url string) (*models.Site, error)
type UpdateSite func(id models.SiteID, url string) error
type RemoveSite func(id models.SiteID) error
