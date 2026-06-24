// FILE: apps/api/internal/atlas/graph/resolver/nutrition_graphql_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Atlas nutrition executable GraphQL schema wiring for factual daily food-log operations and product/template apply extensions.
//   SCOPE: gqlgen request execution through root query and mutation resolvers with service doubles; excludes repository, service, and HTTP middleware integration behavior.
//   DEPENDS: apps/api/internal/atlas/graph/generated, apps/api/internal/atlas/graph/resolver, apps/api/internal/atlas/models, apps/api/internal/atlas/service, gqlgen test client.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestDailyNutritionGraphQL_AddEntryReturnsSnapshotTotals - Proves addDailyNutritionEntry exposes snapshot entry macros and aggregate totals.
//   TestDailyNutritionGraphQL_AddEntryWithoutUserReturnsAuthError - Proves GraphQL auth mapping without invoking the daily service.
//   TestNutritionGraphQL_NutritionProductsDelegatesExecutableSchema - Proves existing nutrition root queries delegate instead of panicking.
//   TestNutritionGraphQL_ProductAllAndRestoreUseProductManagementService - Proves archived products are listed and restore delegates to the product service.
//   TestNutritionGraphQL_ApplyTemplateToWeekReturnsDateStatuses - Proves applyNutritionTemplateToWeek maps enum input/output and per-date statuses.
//   TestNutritionGraphQL_ApplyTemplateToWeekAuthErrorKeepsValidMode - Proves error results do not serialize an invalid empty enum.
// END_MODULE_MAP

package resolver_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/graph/generated"
	"monorepo-template/apps/api/internal/atlas/graph/resolver"
	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

type dailyNutritionGraphQLService struct {
	atlasService.DailyNutritionLogService
	addEntryFn func(ctx context.Context, userID string, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLog, error)
}

func (s *dailyNutritionGraphQLService) AddEntry(ctx context.Context, userID string, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLog, error) {
	return s.addEntryFn(ctx, userID, input)
}

type nutritionProductGraphQLService struct {
	atlasService.NutritionProductService
	listActiveFn func(ctx context.Context, userID string) ([]models.NutritionProduct, error)
	listAllFn    func(ctx context.Context, userID string) ([]models.NutritionProduct, error)
	restoreFn    func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error)
}

func (s *nutritionProductGraphQLService) ListActive(ctx context.Context, userID string) ([]models.NutritionProduct, error) {
	return s.listActiveFn(ctx, userID)
}

func (s *nutritionProductGraphQLService) ListAll(ctx context.Context, userID string) ([]models.NutritionProduct, error) {
	return s.listAllFn(ctx, userID)
}

func (s *nutritionProductGraphQLService) Restore(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
	return s.restoreFn(ctx, userID, id)
}

type nutritionTemplateApplyGraphQLService struct {
	atlasService.NutritionTemplateApplyService
	applyFn func(ctx context.Context, userID string, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error)
}

func (s *nutritionTemplateApplyGraphQLService) ApplyToWeek(ctx context.Context, userID string, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error) {
	return s.applyFn(ctx, userID, templateID, mode)
}

func TestDailyNutritionGraphQL_AddEntryReturnsSnapshotTotals(t *testing.T) {
	res := &resolver.Resolver{}
	setResolverServiceIfPresent(t, res, "DailyNutritionLogService", &dailyNutritionGraphQLService{
		addEntryFn: func(ctx context.Context, userID string, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLog, error) {
			assert.Equal(t, "test-uid", userID)
			assert.Equal(t, "2026-06-24", input.Date.String())
			assert.Equal(t, "product-1", input.ProductID)
			assert.Equal(t, 150.0, input.AmountGrams)
			assert.Equal(t, int32(0), input.Position)

			return &models.DailyNutritionLog{
				ID:     "log-1",
				UserID: "test-uid",
				Date:   "2026-06-24",
				Entries: []models.DailyNutritionEntry{{
					ID:                      "entry-1",
					DailyLogID:              "log-1",
					ProductID:               "product-1",
					ProductNameSnapshot:     "Oats",
					CaloriesPer100gSnapshot: 120,
					ProteinPer100gSnapshot:  10,
					FatPer100gSnapshot:      4,
					CarbsPer100gSnapshot:    30,
					AmountGrams:             150,
					Position:                0,
					Macros: models.NutritionMacros{
						Calories: 180,
						Protein:  15,
						Fat:      6,
						Carbs:    45,
					},
				}},
				Totals: models.NutritionMacros{
					Calories: 180,
					Protein:  15,
					Fat:      6,
					Carbs:    45,
				},
			}, nil
		},
	})

	var response struct {
		AddDailyNutritionEntry struct {
			DailyNutritionLog *struct {
				Date   string `json:"date"`
				Totals struct {
					Calories float64 `json:"calories"`
					Protein  float64 `json:"protein"`
					Fat      float64 `json:"fat"`
					Carbs    float64 `json:"carbs"`
				} `json:"totals"`
				Entries []struct {
					ProductNameSnapshot string  `json:"productNameSnapshot"`
					AmountGrams         float64 `json:"amountGrams"`
					Macros              struct {
						Calories float64 `json:"calories"`
						Protein  float64 `json:"protein"`
						Fat      float64 `json:"fat"`
						Carbs    float64 `json:"carbs"`
					} `json:"macros"`
				} `json:"entries"`
			} `json:"dailyNutritionLog"`
			ValidationError *struct {
				Message string `json:"message"`
				Code    string `json:"code"`
			} `json:"validationError"`
		} `json:"addDailyNutritionEntry"`
	}

	query := `mutation($input: AddDailyNutritionEntryInput!) {
		addDailyNutritionEntry(input: $input) {
			dailyNutritionLog {
				date
				totals { calories protein fat carbs }
				entries {
					productNameSnapshot
					amountGrams
					macros { calories protein fat carbs }
				}
			}
			validationError { message code }
		}
	}`

	err := atlasGraphQLClient(t, res, "test-uid").Post(query, &response, client.Var("input", map[string]any{
		"date":        "2026-06-24",
		"productId":   "product-1",
		"amountGrams": 150,
	}))

	require.NoError(t, err)
	require.NotNil(t, response.AddDailyNutritionEntry.DailyNutritionLog)
	assert.Equal(t, "2026-06-24", response.AddDailyNutritionEntry.DailyNutritionLog.Date)
	assert.Equal(t, 180.0, response.AddDailyNutritionEntry.DailyNutritionLog.Totals.Calories)
	assert.Equal(t, 15.0, response.AddDailyNutritionEntry.DailyNutritionLog.Totals.Protein)
	assert.Equal(t, 6.0, response.AddDailyNutritionEntry.DailyNutritionLog.Totals.Fat)
	assert.Equal(t, 45.0, response.AddDailyNutritionEntry.DailyNutritionLog.Totals.Carbs)
	require.Len(t, response.AddDailyNutritionEntry.DailyNutritionLog.Entries, 1)
	assert.Equal(t, "Oats", response.AddDailyNutritionEntry.DailyNutritionLog.Entries[0].ProductNameSnapshot)
	assert.Equal(t, 150.0, response.AddDailyNutritionEntry.DailyNutritionLog.Entries[0].AmountGrams)
	assert.Equal(t, 180.0, response.AddDailyNutritionEntry.DailyNutritionLog.Entries[0].Macros.Calories)
	assert.Nil(t, response.AddDailyNutritionEntry.ValidationError)
}

func TestDailyNutritionGraphQL_AddEntryWithoutUserReturnsAuthError(t *testing.T) {
	res := &resolver.Resolver{}
	setResolverServiceIfPresent(t, res, "DailyNutritionLogService", &dailyNutritionGraphQLService{
		addEntryFn: func(ctx context.Context, userID string, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLog, error) {
			t.Fatalf("daily nutrition service must not be called without an Atlas user")
			return nil, nil
		},
	})

	var response struct {
		AddDailyNutritionEntry struct {
			DailyNutritionLog *struct{} `json:"dailyNutritionLog"`
			AuthError         *struct {
				Message string `json:"message"`
				Code    string `json:"code"`
			} `json:"authError"`
		} `json:"addDailyNutritionEntry"`
	}

	query := `mutation($input: AddDailyNutritionEntryInput!) {
		addDailyNutritionEntry(input: $input) {
			dailyNutritionLog { id }
			authError { message code }
		}
	}`

	err := atlasGraphQLClient(t, res, "").Post(query, &response, client.Var("input", map[string]any{
		"date":        "2026-06-24",
		"productId":   "product-1",
		"amountGrams": 100,
	}))

	require.NoError(t, err)
	assert.Nil(t, response.AddDailyNutritionEntry.DailyNutritionLog)
	require.NotNil(t, response.AddDailyNutritionEntry.AuthError)
	assert.Equal(t, "AUTH_ERROR", response.AddDailyNutritionEntry.AuthError.Code)
}

func TestNutritionGraphQL_NutritionProductsDelegatesExecutableSchema(t *testing.T) {
	res := &resolver.Resolver{
		NutritionProductService: &nutritionProductGraphQLService{
			listActiveFn: func(ctx context.Context, userID string) ([]models.NutritionProduct, error) {
				assert.Equal(t, "test-uid", userID)
				return []models.NutritionProduct{{
					ID:       "product-1",
					UserID:   "test-uid",
					Name:     "Rice",
					IsActive: true,
				}}, nil
			},
		},
	}

	var response struct {
		NutritionProducts struct {
			Products []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				IsActive bool   `json:"isActive"`
			} `json:"products"`
			AuthError *struct {
				Code string `json:"code"`
			} `json:"authError"`
		} `json:"nutritionProducts"`
	}

	query := `query {
		nutritionProducts {
			products { id name isActive }
			authError { code }
		}
	}`

	err := atlasGraphQLClient(t, res, "test-uid").Post(query, &response)

	require.NoError(t, err)
	require.Len(t, response.NutritionProducts.Products, 1)
	assert.Equal(t, "Rice", response.NutritionProducts.Products[0].Name)
	assert.True(t, response.NutritionProducts.Products[0].IsActive)
	assert.Nil(t, response.NutritionProducts.AuthError)
}

func TestNutritionGraphQL_ProductAllAndRestoreUseProductManagementService(t *testing.T) {
	res := &resolver.Resolver{
		NutritionProductService: &nutritionProductGraphQLService{
			listAllFn: func(ctx context.Context, userID string) ([]models.NutritionProduct, error) {
				assert.Equal(t, "test-uid", userID)
				return []models.NutritionProduct{
					{ID: "product-active", UserID: "test-uid", Name: "Milk", IsActive: true},
					{ID: "product-archived", UserID: "test-uid", Name: "Old snack", IsActive: false},
				}, nil
			},
			restoreFn: func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "product-archived", id)
				return &models.NutritionProduct{ID: id, UserID: userID, Name: "Old snack", IsActive: true}, nil
			},
		},
	}
	gql := atlasGraphQLClient(t, res, "test-uid")

	var listResponse struct {
		NutritionProductsAll struct {
			Products []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				IsActive bool   `json:"isActive"`
			} `json:"products"`
		} `json:"nutritionProductsAll"`
	}

	listQuery := `query {
		nutritionProductsAll {
			products { id name isActive }
		}
	}`

	require.NoError(t, gql.Post(listQuery, &listResponse))
	require.Len(t, listResponse.NutritionProductsAll.Products, 2)
	assert.Equal(t, "product-active", listResponse.NutritionProductsAll.Products[0].ID)
	assert.True(t, listResponse.NutritionProductsAll.Products[0].IsActive)
	assert.Equal(t, "product-archived", listResponse.NutritionProductsAll.Products[1].ID)
	assert.False(t, listResponse.NutritionProductsAll.Products[1].IsActive)

	var restoreResponse struct {
		RestoreNutritionProduct struct {
			NutritionProduct *struct {
				ID       string `json:"id"`
				IsActive bool   `json:"isActive"`
			} `json:"nutritionProduct"`
			NotFoundError *struct {
				Code string `json:"code"`
			} `json:"notFoundError"`
		} `json:"restoreNutritionProduct"`
	}

	restoreMutation := `mutation($id: ID!) {
		restoreNutritionProduct(id: $id) {
			nutritionProduct { id isActive }
			notFoundError { code }
		}
	}`

	require.NoError(t, gql.Post(restoreMutation, &restoreResponse, client.Var("id", "product-archived")))
	require.NotNil(t, restoreResponse.RestoreNutritionProduct.NutritionProduct)
	assert.Equal(t, "product-archived", restoreResponse.RestoreNutritionProduct.NutritionProduct.ID)
	assert.True(t, restoreResponse.RestoreNutritionProduct.NutritionProduct.IsActive)
	assert.Nil(t, restoreResponse.RestoreNutritionProduct.NotFoundError)
}

func TestNutritionGraphQL_ApplyTemplateToWeekReturnsDateStatuses(t *testing.T) {
	res := &resolver.Resolver{}
	setResolverServiceIfPresent(t, res, "NutritionTemplateApplyService", &nutritionTemplateApplyGraphQLService{
		applyFn: func(ctx context.Context, userID string, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error) {
			assert.Equal(t, "test-uid", userID)
			assert.Equal(t, "template-1", templateID)
			assert.Equal(t, models.ApplyModeSeedEmptyDays, mode)
			return &models.NutritionTemplateApplyResult{
				WeekStartDate: "2026-06-22",
				WeekEndDate:   "2026-06-28",
				Mode:          models.ApplyModeSeedEmptyDays,
				Dates: []models.NutritionTemplateApplyDateResult{
					{Date: "2026-06-22", Status: models.ApplyDateCreated, EntryCount: 2},
					{Date: "2026-06-23", Status: models.ApplyDateSkipped, EntryCount: 1, Reason: stringPointer("day has entries")},
				},
			}, nil
		},
	})

	var response struct {
		ApplyNutritionTemplateToWeek struct {
			WeekStartDate string `json:"weekStartDate"`
			WeekEndDate   string `json:"weekEndDate"`
			Mode          string `json:"mode"`
			Dates         []struct {
				Date       string  `json:"date"`
				Status     string  `json:"status"`
				EntryCount int     `json:"entryCount"`
				Reason     *string `json:"reason"`
			} `json:"dates"`
			ValidationError *struct {
				Code string `json:"code"`
			} `json:"validationError"`
		} `json:"applyNutritionTemplateToWeek"`
	}

	query := `mutation($templateId: ID!, $mode: NutritionTemplateApplyMode!) {
		applyNutritionTemplateToWeek(templateId: $templateId, mode: $mode) {
			weekStartDate
			weekEndDate
			mode
			dates { date status entryCount reason }
			validationError { code }
		}
	}`

	err := atlasGraphQLClient(t, res, "test-uid").Post(
		query,
		&response,
		client.Var("templateId", "template-1"),
		client.Var("mode", "SEED_EMPTY_DAYS"),
	)

	require.NoError(t, err)
	assert.Equal(t, "2026-06-22", response.ApplyNutritionTemplateToWeek.WeekStartDate)
	assert.Equal(t, "2026-06-28", response.ApplyNutritionTemplateToWeek.WeekEndDate)
	assert.Equal(t, "SEED_EMPTY_DAYS", response.ApplyNutritionTemplateToWeek.Mode)
	require.Len(t, response.ApplyNutritionTemplateToWeek.Dates, 2)
	assert.Equal(t, "2026-06-22", response.ApplyNutritionTemplateToWeek.Dates[0].Date)
	assert.Equal(t, "created", response.ApplyNutritionTemplateToWeek.Dates[0].Status)
	assert.Equal(t, 2, response.ApplyNutritionTemplateToWeek.Dates[0].EntryCount)
	assert.Equal(t, "skipped", response.ApplyNutritionTemplateToWeek.Dates[1].Status)
	require.NotNil(t, response.ApplyNutritionTemplateToWeek.Dates[1].Reason)
	assert.Equal(t, "day has entries", *response.ApplyNutritionTemplateToWeek.Dates[1].Reason)
	assert.Nil(t, response.ApplyNutritionTemplateToWeek.ValidationError)
}

func TestNutritionGraphQL_ApplyTemplateToWeekAuthErrorKeepsValidMode(t *testing.T) {
	res := &resolver.Resolver{}
	setResolverServiceIfPresent(t, res, "NutritionTemplateApplyService", &nutritionTemplateApplyGraphQLService{
		applyFn: func(ctx context.Context, userID string, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error) {
			t.Fatalf("template apply service must not be called without an Atlas user")
			return nil, nil
		},
	})

	var response struct {
		ApplyNutritionTemplateToWeek struct {
			Mode      string `json:"mode"`
			AuthError *struct {
				Code string `json:"code"`
			} `json:"authError"`
		} `json:"applyNutritionTemplateToWeek"`
	}

	query := `mutation($templateId: ID!, $mode: NutritionTemplateApplyMode!) {
		applyNutritionTemplateToWeek(templateId: $templateId, mode: $mode) {
			mode
			authError { code }
		}
	}`

	err := atlasGraphQLClient(t, res, "").Post(
		query,
		&response,
		client.Var("templateId", "template-1"),
		client.Var("mode", "SEED_EMPTY_DAYS"),
	)

	require.NoError(t, err)
	assert.Equal(t, "SEED_EMPTY_DAYS", response.ApplyNutritionTemplateToWeek.Mode)
	require.NotNil(t, response.ApplyNutritionTemplateToWeek.AuthError)
	assert.Equal(t, "AUTH_ERROR", response.ApplyNutritionTemplateToWeek.AuthError.Code)
}

func atlasGraphQLClient(t *testing.T, res *resolver.Resolver, userID string) *client.Client {
	t.Helper()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: res}))
	var h http.Handler = srv
	if userID != "" {
		h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := middleware.ContextWithAtlasUserID(req.Context(), userID)
			srv.ServeHTTP(w, req.WithContext(ctx))
		})
	}
	return client.New(h)
}

func setResolverServiceIfPresent(t *testing.T, res *resolver.Resolver, fieldName string, service any) {
	t.Helper()

	field := reflect.ValueOf(res).Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return
	}
	if !field.CanSet() {
		t.Fatalf("resolver field %s cannot be set", fieldName)
	}
	value := reflect.ValueOf(service)
	if !value.Type().AssignableTo(field.Type()) {
		t.Fatalf("resolver field %s expects %s, got %s", fieldName, field.Type(), value.Type())
	}
	field.Set(value)
}

func stringPointer(value string) *string {
	return &value
}
