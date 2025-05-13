package main

import (
	"github.com/movie-app/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.Module).Run()
}
