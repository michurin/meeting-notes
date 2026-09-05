package common

import "context"

type CityRepositoryInterface interface {
	Add(ctx context.Context, name string) error
}
