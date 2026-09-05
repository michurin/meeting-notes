package common

import (
	"context"
	"fmt"
)

type Case struct {
	repo CityRepositoryInterface
}

func NewCase(repo CityRepositoryInterface) *Case {
	return &Case{repo: repo}
}

func (c *Case) Add(ctx context.Context, cities []string) error {
	for _, ct := range cities {
		if err := c.repo.Add(ctx, ct); err != nil {
			return fmt.Errorf("insert city: %q: %w", ct, err)
		}
	}
	return nil
}
