package middleware

import (
	"erp_monolith/backend/internal/platform/response"
	"erp_monolith/backend/internal/shared/audit"
	"erp_monolith/backend/internal/shared/types"
	"net/http"

	"github.com/labstack/echo/v5"
)

const (
	HeaderCompanyID = "X-Company-ID"
	HeaderBranchID  = "X-Branch-ID"
)

// TenantScope middleware ensures every authenticated request has valid CompanyID & BranchID
func TenantScope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			companyHeader := c.Request().Header.Get(HeaderCompanyID)
			branchHeader := c.Request().Header.Get(HeaderBranchID)

			if companyHeader == "" {
				return response.Error(c, http.StatusBadRequest, "MISSING_TENANT", "Header X-Company-ID wajib disertakan.")
			}

			companyID, err := types.ParseID(companyHeader)
			if err != nil {
				return response.Error(c, http.StatusBadRequest, "INVALID_TENANT", "Nilai X-Company-ID tidak valid.")
			}

			branchID := types.EmptyID()
			if branchHeader != "" {
				if parsedBranch, err := types.ParseID(branchHeader); err == nil {
					branchID = parsedBranch
				}
			}

			// Attach to context
			actor := audit.Actor{
				CompanyID: companyID,
				BranchID:  branchID,
				IP:        c.RealIP(),
				UserAgent: c.Request().UserAgent(),
			}

			ctx := audit.WithActor(c.Request().Context(), actor)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
