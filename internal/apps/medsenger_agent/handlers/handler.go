package handlers

import (
	"github.com/TikhonP/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
)

type MedsengerAgentHandler struct {
	MaigoClient *maigo.Client
	Db          db.ModelsFactory
}
