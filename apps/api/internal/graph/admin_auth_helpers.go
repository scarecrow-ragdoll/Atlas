// FILE: apps/api/internal/graph/admin_auth_helpers.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide shared GraphQL admin-auth resolver helper functions.
//   SCOPE: Admin principal guard logging and service-to-GraphQL admin model mapping; excludes generated resolver method wiring and service behavior.
//   DEPENDS: apps/api/internal/middleware, apps/api/internal/service, apps/api/internal/graph/model, libs/go/logger.
//   LINKS: M-API / M-GRAPHQL-SCHEMA / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: LOCALS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   requireAdmin - Converts request-scoped admin principal into a service.Admin or auth error.
//   mapAdmin - Converts service.Admin into GraphQL AdminUser.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Extracted admin resolver helpers from gqlgen resolver output.
// END_CHANGE_SUMMARY

package graph

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/graph/model"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func requireAdmin(ctx context.Context) (*service.Admin, error) {
	log := logger.FromContext(ctx)
	principal, ok := middleware.GetAdminPrincipal(ctx)
	if !ok {
		log.Warn("[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL] admin principal missing")
		return nil, fmt.Errorf("admin authentication required")
	}
	log.Debug("[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL] admin principal accepted", zap.String("admin_id", principal.ID))
	return &service.Admin{
		ID:        principal.ID,
		Email:     principal.Email,
		Name:      principal.Name,
		Role:      principal.Role,
		IsActive:  true,
		CreatedAt: principal.CreatedAt,
		UpdatedAt: principal.UpdatedAt,
	}, nil
}

func mapAdmin(admin *service.Admin) *model.AdminUser {
	if admin == nil {
		return nil
	}
	return &model.AdminUser{
		ID:        admin.ID,
		Email:     admin.Email,
		Name:      admin.Name,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}
