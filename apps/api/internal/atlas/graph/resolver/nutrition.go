// FILE: apps/api/internal/atlas/graph/resolver/nutrition.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement all WAVE-05 Nutrition GraphQL resolvers with PIN auth guard and union error returns following the existing cardio resolver pattern.
//   SCOPE: All nutrition queries and mutations: NutritionProduct CRUD, NutritionTemplate CRUD, NutritionTemplateItem CRUD, DailyNutritionOverride CRUD, DailyNutritionOverrideItem CRUD, NutritionMacros calculation. Each resolver checks middleware.GetAtlasUserID for auth, maps sentinel service errors to union result fields.
//   DEPENDS: middleware for auth, models for result types, service for business logic.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added product list-all and restore GraphQL adapter methods for archived product management.
// END_CHANGE_SUMMARY

package resolver

import (
	"context"
	"errors"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) GetNutritionProducts(ctx context.Context) (*models.NutritionProductsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductsResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	products, err := r.NutritionProductService.ListActive(ctx, userID)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionProductsResult{Products: products}, nil
}

func (r *Resolver) GetNutritionProductsAll(ctx context.Context) (*models.NutritionProductsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductsResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	products, err := r.NutritionProductService.ListAll(ctx, userID)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionProductsResult{Products: products}, nil
}

func (r *Resolver) GetNutritionProduct(ctx context.Context, id string) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrProductNotFound) {
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) CreateNutritionProduct(ctx context.Context, input models.CreateProductInput) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrProductNameRequired),
			errors.Is(err, atlasService.ErrProductMacroNegative),
			errors.Is(err, atlasService.ErrProductNameTooLong):
			return &models.NutritionProductResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) UpdateNutritionProduct(ctx context.Context, id string, input models.UpdateProductInput) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrProductNameRequired),
			errors.Is(err, atlasService.ErrProductMacroNegative),
			errors.Is(err, atlasService.ErrProductNameTooLong):
			return &models.NutritionProductResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrProductNotFound):
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) DeleteNutritionProduct(ctx context.Context, id string) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrProductNotFound) {
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) RestoreNutritionProduct(ctx context.Context, id string) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Restore(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrProductNotFound) {
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) GetNutritionTemplates(ctx context.Context, startDate, endDate string) (*models.NutritionTemplatesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplatesResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	templates, err := r.NutritionTemplateService.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionTemplatesResult{Templates: templates}, nil
}

func (r *Resolver) GetNutritionTemplate(ctx context.Context, id string) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateNotFound) {
			return &models.NutritionTemplateResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) GetNutritionTemplateCurrent(ctx context.Context, weekStartDate string) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.GetCurrent(ctx, userID, weekStartDate)
	if err != nil {
		return nil, nil
	}
	if tmpl == nil {
		return &models.NutritionTemplateResult{}, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) CreateNutritionTemplate(ctx context.Context, input models.CreateTemplateInput) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrTemplateWeekRequired):
			return &models.NutritionTemplateResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) UpdateNutritionTemplate(ctx context.Context, id string, input models.UpdateTemplateInput) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.Update(ctx, userID, id, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateNotFound) {
			return &models.NutritionTemplateResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) DeleteNutritionTemplate(ctx context.Context, id string) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateNotFound) {
			return &models.NutritionTemplateResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) CreateNutritionTemplateItem(ctx context.Context, input models.CreateTemplateItemInput) (*models.NutritionTemplateItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.NutritionTemplateItemService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrTemplateItemAmountInvalid):
			return &models.NutritionTemplateItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrTemplateNotFound):
			return &models.NutritionTemplateItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: "template not found", Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrProductNotFound):
			return &models.NutritionTemplateItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionTemplateItemResult{NutritionTemplateItem: item}, nil
}

func (r *Resolver) UpdateNutritionTemplateItem(ctx context.Context, id string, input models.UpdateTemplateItemInput) (*models.NutritionTemplateItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.NutritionTemplateItemService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrTemplateItemAmountInvalid):
			return &models.NutritionTemplateItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrTemplateItemNotFound):
			return &models.NutritionTemplateItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "item not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionTemplateItemResult{NutritionTemplateItem: item}, nil
}

func (r *Resolver) DeleteNutritionTemplateItem(ctx context.Context, id string) (*models.NutritionTemplateItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.NutritionTemplateItemService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateItemNotFound) {
			return &models.NutritionTemplateItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "item not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateItemResult{NutritionTemplateItem: item}, nil
}

func (r *Resolver) GetDailyNutritionOverrides(ctx context.Context, startDate, endDate string) (*models.DailyNutritionOverridesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverridesResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	overrides, err := r.DailyNutritionOverrideService.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, nil
	}

	return &models.DailyNutritionOverridesResult{Overrides: overrides}, nil
}

func (r *Resolver) GetDailyNutritionOverride(ctx context.Context, id string) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideNotFound) {
			return &models.DailyNutritionOverrideResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) GetDailyNutritionOverrideByDate(ctx context.Context, date string) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, nil
	}
	if override == nil {
		return &models.DailyNutritionOverrideResult{}, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) CreateDailyNutritionOverride(ctx context.Context, input models.CreateOverrideInput) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrOverrideDateRequired):
			return &models.DailyNutritionOverrideResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) UpdateDailyNutritionOverride(ctx context.Context, id string, input models.UpdateOverrideInput) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.Update(ctx, userID, id, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideNotFound) {
			return &models.DailyNutritionOverrideResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) DeleteDailyNutritionOverride(ctx context.Context, id string) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideNotFound) {
			return &models.DailyNutritionOverrideResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) CreateDailyNutritionOverrideItem(ctx context.Context, input models.CreateOverrideItemInput) (*models.DailyNutritionOverrideItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.DailyNutritionOverrideService.CreateItem(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrOverrideNotFound):
			return &models.DailyNutritionOverrideItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: "override not found", Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrOverrideItemAmountInvalid),
			errors.Is(err, atlasService.ErrOverrideItemOperationInvalid):
			return &models.DailyNutritionOverrideItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.DailyNutritionOverrideItemResult{DailyNutritionOverrideItem: item}, nil
}

func (r *Resolver) UpdateDailyNutritionOverrideItem(ctx context.Context, id string, input models.UpdateOverrideItemInput) (*models.DailyNutritionOverrideItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.DailyNutritionOverrideService.UpdateItem(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrOverrideItemAmountInvalid),
			errors.Is(err, atlasService.ErrOverrideItemOperationInvalid):
			return &models.DailyNutritionOverrideItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrOverrideItemNotFound):
			return &models.DailyNutritionOverrideItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override item not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.DailyNutritionOverrideItemResult{DailyNutritionOverrideItem: item}, nil
}

func (r *Resolver) DeleteDailyNutritionOverrideItem(ctx context.Context, id string) (*models.DailyNutritionOverrideItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.DailyNutritionOverrideService.DeleteItem(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideItemNotFound) {
			return &models.DailyNutritionOverrideItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override item not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideItemResult{DailyNutritionOverrideItem: item}, nil
}

func (r *Resolver) GetNutritionMacros(ctx context.Context, weekStartDate string, date *string) (*models.NutritionMacrosResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionMacrosResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	d := ""
	if date != nil {
		d = *date
	}

	macros, err := r.NutritionMacroService.Calculate(ctx, userID, weekStartDate, d)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionMacrosResult{Macros: macros}, nil
}
