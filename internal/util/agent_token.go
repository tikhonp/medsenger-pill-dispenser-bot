// Package util provides utility functions and middleware for the Medsenger Pill Dispenser Bot.
package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/maigo"
)

type agentTokenModel struct {
	AgentToken string `json:"agent_token" validate:"required"`
}

func processAgentToken(agentToken string, c echo.Context, client *maigo.Client, roles []maigo.RequestRole) error {
	data, err := client.DecodeAgentJWT(agentToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid jwt key.")
	}
	// Check if the token has at least one of the required roles
	for _, role := range data.Roles {
		if !slices.Contains(roles, role) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid jwt key role.")
		}
	}
	if data.ContractID != nil {
		c.Set("contract_id", *data.ContractID)
	}
	c.Set("agent_token", agentToken)
	return nil
}

func AgentTokenJSON(client *maigo.Client, roles ...maigo.RequestRole) echo.MiddlewareFunc {
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

			data := new(agentTokenModel)
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON.")
			}
			if err := c.Validate(data); err != nil {
				return err
			}
			if err := processAgentToken(data.AgentToken, c, client, roles); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func AgentTokenGetParam(client *maigo.Client, roles ...maigo.RequestRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			agentToken := c.QueryParam("agent_token")
			if err := processAgentToken(agentToken, c, client, roles); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func AgentTokenForm(client *maigo.Client, roles ...maigo.RequestRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			agentToken := c.FormValue("agent-token")
			if err := processAgentToken(agentToken, c, client, roles); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func GetContractID(c echo.Context) (int, error) {
	contractID, ok := c.Get("contract_id").(int)
	if ok {
		return contractID, nil
	}
	return 0, errors.New("no contract ID in context")
}
