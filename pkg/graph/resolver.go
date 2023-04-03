package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/KarnaTushar/videoApp-backend/pkg/models"
	"github.com/KarnaTushar/videoApp-backend/utils"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is used for state management
type Resolver struct {
	DB     *models.Database
	Logger *utils.Logger
}
