package driver

import "http-server/internal/core/model"

type GetSiteList func() ([]*model.Site, error)
type AddSite func(url string) (*model.Site, error)
type UpdateSite func(id model.SiteID, url string) error
type RemoveSite func(id model.SiteID) error
