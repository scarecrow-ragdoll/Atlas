// FILE: apps/api/internal/atlas/graph/resolver/nutrition_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for nutrition GraphQL resolver error mapping.
//   SCOPE: Focused nutrition resolver union mapping regressions; excludes service and repository behavior.
//   DEPENDS: apps/api/internal/atlas/graph/resolver, apps/api/internal/atlas/service, apps/api/internal/atlas/models, apps/api/internal/atlas/middleware.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

package resolver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/graph/resolver"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

type mockNutritionTemplateItemService struct {
	atlasService.NutritionTemplateItemService
	createFn func(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error)
}

func (m *mockNutritionTemplateItemService) Create(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error) {
	return m.createFn(ctx, userID, input)
}

func TestCreateNutritionTemplateItem_ProductNotFoundMapsToTypedUnion(t *testing.T) {
	r := &resolver.Resolver{
		NutritionTemplateItemService: &mockNutritionTemplateItemService{
			createFn: func(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "product-1", input.ProductID)
				return nil, atlasService.ErrProductNotFound
			},
		},
	}

	result, err := r.CreateNutritionTemplateItem(userCtx("test-uid"), models.CreateTemplateItemInput{
		TemplateID:  "template-1",
		ProductID:   "product-1",
		AmountGrams: 100,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.NotFoundErr)
	assert.Equal(t, "product not found", result.NotFoundErr.Message)
	assert.Equal(t, models.NutritionErrorNotFound, result.NotFoundErr.Code)
	assert.Nil(t, result.NutritionTemplateItem)
	assert.Nil(t, result.ValidationErr)
	assert.Nil(t, result.AuthErr)
}
