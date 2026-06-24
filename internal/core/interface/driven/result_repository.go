package driven

import "http-server/internal/core/model"

type ResultRepository interface {
	GetSiteLatest(siteId model.SiteID) (*model.Result, error)
	Save(r *model.Result) (model.ResultID, error)
}
