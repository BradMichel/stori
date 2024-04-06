//go:build wireinject

package v1

import "github.com/google/wire"

func Initialize() (*Handler, error) {
	wire.Build(stdSet)
	return &Handler{}, nil
}
