package handlers

import (
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
)

type PillDispenserHandler struct {
	Db db.ModelsFactory
}
