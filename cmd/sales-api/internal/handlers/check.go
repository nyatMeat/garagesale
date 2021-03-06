package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nyatmeat/garagesale/internal/platform/database"
	"github.com/nyatmeat/garagesale/internal/platform/web"
)

type Check struct {
	DB *sqlx.DB
}

//Healt responds with a 200 if service is healthy and ready for trafic
func (c *Check) Health(ctxt context.Context, w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(ctxt, c.DB); err != nil {
		health.Status = "DB is not ready"
		return web.Respond(ctxt, w, health, http.StatusInternalServerError)
	}
	health.Status = "OK"
	return web.Respond(ctxt, w, health, http.StatusOK)
}
