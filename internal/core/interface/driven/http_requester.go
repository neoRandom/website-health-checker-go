package driven

import "http-server/internal/core/model"

type HTTPRequester interface {
	CheckSite(s *model.Site) (*model.Result, error)
}
