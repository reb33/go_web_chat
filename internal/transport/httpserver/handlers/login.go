package handlers

import (
	"log/slog"
	"net/http"
	"time"
	"github.com/google/uuid"
	"web_chat/internal/transport/httpserver/ports"
	"web_chat/internal/transport/httpserver/responses"
	"web_chat/internal/usecase/auth"

	"github.com/labstack/echo/v4"
)

const (
	binRequestLength  = 8
	countryCodeLength = 3
	locBody           = "body"
	fieldBin          = "bin"
	fieldCountryCode  = "country_code"
)

type loginRequest struct {
	Username string `json:"username"`
}

type loginResponse struct {
	Id        uuid.UUID
	Username  string
	CreatedAt time.Time
}

func Login(uc ports.AuthUseCase, logger *slog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req loginRequest
		if err := c.Bind(&req); err != nil {
			resp := responses.CustomError(responses.TypeJSONInvalid, []string{locBody}, "JSON decode error", nil)
			return c.JSON(http.StatusUnprocessableEntity, resp)
		}

		ctx := c.Request().Context()
		res, err := uc.Execute(ctx, auth.Request{Username: req.Username})
		if err != nil {
			logger.ErrorContext(ctx, "authorization failed", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, responses.DetailResponse{Detail: "internal error"})
		}
		
		return c.JSON(http.StatusOK, loginResponse{
			Id: res.ID,
			Username: res.Username,
			CreatedAt: res.CreatedAt,
		})
	}
}
