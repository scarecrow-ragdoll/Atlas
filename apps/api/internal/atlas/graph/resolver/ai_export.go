// FILE: apps/api/internal/atlas/graph/resolver/ai_export.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-07 AiExport queries and mutations.
//   SCOPE: getAiExport, listAiExports, createAiExportPrompt, generateAiExport, deleteAiExport.
//   DEPENDS: apps/api/internal/atlas/service.AiExportService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, apps/api/internal/appconfig, generated ai_export.resolvers.go.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package resolver

import (
	"context"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
)

func (r *Resolver) GetAiExport(ctx context.Context, id string) (*models.AiExportResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiExportResult{
			AuthErr: &models.AiExportAuthErr{
				Message: "unauthorized",
				Code:    models.AiExportErrorAuth,
			},
		}, nil
	}

	export, err := r.AiExportService.GetByID(ctx, userID, id)
	if err != nil {
		return nil, nil
	}

	return &models.AiExportResult{Export: export}, nil
}

func (r *Resolver) ListAiExports(ctx context.Context) (*models.AiExportsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiExportsResult{
			AuthErr: &models.AiExportAuthErr{
				Message: "unauthorized",
				Code:    models.AiExportErrorAuth,
			},
		}, nil
	}

	exports, err := r.AiExportService.List(ctx, userID)
	if err != nil {
		return nil, nil
	}

	return &models.AiExportsResult{Exports: exports}, nil
}

func (r *Resolver) CreateAiExportPrompt(ctx context.Context, input models.CreateAiExportInput) (*models.AiExportResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiExportResult{
			AuthErr: &models.AiExportAuthErr{
				Message: "unauthorized",
				Code:    models.AiExportErrorAuth,
			},
		}, nil
	}

	cfg := r.AiExportConfig
	export, prompt, err := r.AiExportService.Generate(ctx, userID, input, cfg.MaxRangeDays, cfg.MaxExportSizeBytes, cfg.BasePath)
	if err != nil {
		return nil, nil
	}

	export.GeneratedPrompt = prompt

	return &models.AiExportResult{Export: export}, nil
}

func (r *Resolver) GenerateAiExport(ctx context.Context, id string) (*models.AiExportResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiExportResult{
			AuthErr: &models.AiExportAuthErr{
				Message: "unauthorized",
				Code:    models.AiExportErrorAuth,
			},
		}, nil
	}

	export, err := r.AiExportService.GetByID(ctx, userID, id)
	if err != nil {
		return nil, nil
	}

	return &models.AiExportResult{Export: export}, nil
}

func (r *Resolver) DeleteAiExport(ctx context.Context, id string) (*models.AiExportResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiExportResult{
			AuthErr: &models.AiExportAuthErr{
				Message: "unauthorized",
				Code:    models.AiExportErrorAuth,
			},
		}, nil
	}

	export, err := r.AiExportService.Delete(ctx, userID, id)
	if err != nil {
		return nil, nil
	}

	return &models.AiExportResult{Export: export}, nil
}