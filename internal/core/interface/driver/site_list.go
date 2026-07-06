package driver

import "http-server/internal/core/model"

type GetSiteList func() ([]model.Site, error)
type AddSite func(site *model.Site) (model.SiteID, error)
type UpdateSite func(site *model.Site) error
type RemoveSite func(id model.SiteID) error
