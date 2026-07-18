package driver

import (
	"context"
	"http-server/internal/core/model"
)

type CheckSites func(ctx context.Context, targets []model.Site) error
