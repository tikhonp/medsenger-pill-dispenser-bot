// Package util provides utility functions and middleware for the Medsenger Pill Dispenser Bot.
package util

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

type apiKeyModel struct {
	APIKey string `json:"api_key" validate:"required"`
}

func (k *apiKeyModel) isValid(cfg *config.Server) bool {
	return cfg.MedsengerAgentKey == k.APIKey
}

func APIKeyJSON(cfg *config.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Workaround to read request body twice
			req := c.Request()
			bodyBytes, _ := io.ReadAll(req.Body)
			if err := req.Body.Close(); err != nil {
				return err
			}
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.SetRequest(req)

			data := new(apiKeyModel)
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON.")
			}
			if err := c.Validate(data); err != nil {
				return err
			}
			if !data.isValid(cfg) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key.")
			}
			return next(c)
		}
	}
}

func APIKeyGetParam(cfg *config.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.QueryParam("api_key")
			if apiKey != cfg.MedsengerAgentKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key.")
			}
			return next(c)
		}
	}
}

func AgentTokenGetParam(db db.ModelsFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			agentToken := c.QueryParam("agent_token")
			contract, err := db.Contracts().GetByAgentToken(agentToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid agent token", err)
			}
			c.Set("contract", *contract)
			return next(c)
		}
	}
}

func AgentTokenForm(db db.ModelsFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			agentToken := c.FormValue("agent-token")
			contract, err := db.Contracts().GetByAgentToken(agentToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid agent token", err)
			}
			c.Set("contract", *contract)
			return next(c)
		}
	}
}

func ContractIDQueryParam(db db.ModelsFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contractIDStr := c.QueryParam("contract_id")
			if contractIDStr == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "provide contract id")
			}
			contractID, err := strconv.Atoi(contractIDStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid contract id")
			}
			contract, err := db.Contracts().Get(contractID)
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, "contract not found")
			} else if err != nil {
				return err
			}
			c.Set("contract", *contract)
			return next(c)
		}
	}
}

func GetContract(c echo.Context) (*models.Contract, error) {
	contract, ok := c.Get("contract").(models.Contract)
	if ok {
		return &contract, nil
	}
	return nil, errors.New("no contract in context")
}
