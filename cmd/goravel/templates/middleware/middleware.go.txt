package middleware

import (
	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/models"
)

type Middleware struct {
	App    *goravel.Goravel
	Models *models.Models
}
