package docs

import (
	"github.com/swaggo/swag"
)

go:generate swag init -g ../cmd/main.go -o ./
