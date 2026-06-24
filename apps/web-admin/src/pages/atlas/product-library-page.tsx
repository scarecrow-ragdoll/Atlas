// FILE: apps/web-admin/src/pages/atlas/product-library-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the Atlas nutrition product library backed by the factual Atlas nutrition API.
//   SCOPE: Loads user-owned nutrition products, filters active/archived records, validates product form input, and performs create/update/archive/restore mutations; excludes daily food-log entry editing and weekly template management.
//   DEPENDS: @tanstack/react-query, react, apps/web-admin/src/app/i18n.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts, apps/web-admin/src/styles/atlas.css, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - API-backed Product Library route content for private nutrition products.
// END_MODULE_MAP

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { type FormEvent, useMemo, useState } from 'react';
import { useI18n } from '../../app/i18n';
import '../../styles/atlas.css';
import {
  archiveAtlasNutritionProduct,
  createAtlasNutritionProduct,
  listAtlasNutritionProducts,
  restoreAtlasNutritionProduct,
  updateAtlasNutritionProduct,
  type AtlasNutritionProduct,
  type CreateAtlasNutritionProductInput,
} from './nutrition-api';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  Label,
  Skeleton,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  Textarea,
} from '@shared/ui';

type ProductStatusFilter = 'active' | 'archived';

type ProductFormState = {
  name: string;
  caloriesPer100g: string;
  proteinPer100g: string;
  fatPer100g: string;
  carbsPer100g: string;
  notes: string;
};

type MutationKind = 'create' | 'update' | 'archive' | 'restore';

const PRODUCTS_QUERY_KEY = ['atlas-nutrition-products', 'with-archived'] as const;

const emptyForm: ProductFormState = {
  name: '',
  caloriesPer100g: '',
  proteinPer100g: '',
  fatPer100g: '',
  carbsPer100g: '',
  notes: '',
};

function formFromProduct(product: AtlasNutritionProduct): ProductFormState {
  return {
    name: product.name,
    caloriesPer100g: String(product.caloriesPer100g),
    proteinPer100g: String(product.proteinPer100g),
    fatPer100g: String(product.fatPer100g),
    carbsPer100g: String(product.carbsPer100g),
    notes: product.notes ?? '',
  };
}

function formatMacro(value: number): string {
  return Number.isInteger(value) ? String(value) : value.toFixed(1);
}

function errorMessageFromUnknown(error: unknown): string {
  return error instanceof Error ? error.message : 'Request failed';
}

function updateProductInList(
  products: AtlasNutritionProduct[] | undefined,
  product: AtlasNutritionProduct,
): AtlasNutritionProduct[] {
  const existingProducts = products ?? [];
  const existingIndex = existingProducts.findIndex((candidate) => candidate.id === product.id);

  if (existingIndex === -1) {
    return [product, ...existingProducts];
  }

  return existingProducts.map((candidate) => (candidate.id === product.id ? product : candidate));
}

// START_CONTRACT: parseProductForm
//   PURPOSE: Validate controlled product form values and convert them to API input.
//   INPUTS: { form: ProductFormState - current form fields, messages: validation messages }
//   OUTPUTS: { input or error - parsed API input or user-visible validation error }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
// END_CONTRACT: parseProductForm
function parseProductForm(
  form: ProductFormState,
  messages: { nameRequired: string; macroInvalid: string },
): { input: CreateAtlasNutritionProductInput; error: null } | { input: null; error: string } {
  const name = form.name.trim();

  if (!name) {
    return { input: null, error: messages.nameRequired };
  }

  const macroValues = {
    caloriesPer100g: Number(form.caloriesPer100g || 0),
    proteinPer100g: Number(form.proteinPer100g || 0),
    fatPer100g: Number(form.fatPer100g || 0),
    carbsPer100g: Number(form.carbsPer100g || 0),
  };

  if (Object.values(macroValues).some((value) => !Number.isFinite(value) || value < 0)) {
    return { input: null, error: messages.macroInvalid };
  }

  return {
    input: {
      name,
      caloriesPer100g: macroValues.caloriesPer100g,
      proteinPer100g: macroValues.proteinPer100g,
      fatPer100g: macroValues.fatPer100g,
      carbsPer100g: macroValues.carbsPer100g,
      notes: form.notes.trim() || null,
    },
    error: null,
  };
}

// START_CONTRACT: ProductLibraryPage
//   PURPOSE: Render and mutate the private nutrition product catalog from the Atlas nutrition API.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - product library page with loading, error, empty, form, table, and success states }
//   SIDE_EFFECTS: Sends product create/update/archive/restore requests and updates React Query cache on success.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
// END_CONTRACT: ProductLibraryPage
export default function ProductLibraryPage() {
  const { t } = useI18n();
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<ProductStatusFilter>('active');
  const [searchTerm, setSearchTerm] = useState('');
  const [form, setForm] = useState<ProductFormState>(emptyForm);
  const [editingProduct, setEditingProduct] = useState<AtlasNutritionProduct | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const productsQuery = useQuery({
    queryKey: PRODUCTS_QUERY_KEY,
    queryFn: () => listAtlasNutritionProducts({ includeArchived: true }),
  });

  const products = productsQuery.data ?? [];
  const filteredProducts = useMemo(() => {
    const normalizedSearch = searchTerm.trim().toLocaleLowerCase();

    return products.filter((product) => {
      const matchesStatus = statusFilter === 'active' ? product.isActive : !product.isActive;
      const matchesSearch =
        !normalizedSearch ||
        product.name.toLocaleLowerCase().includes(normalizedSearch) ||
        (product.notes ?? '').toLocaleLowerCase().includes(normalizedSearch);

      return matchesStatus && matchesSearch;
    });
  }, [products, searchTerm, statusFilter]);

  function setCachedProduct(product: AtlasNutritionProduct) {
    queryClient.setQueryData<AtlasNutritionProduct[]>(PRODUCTS_QUERY_KEY, (currentProducts) =>
      updateProductInList(currentProducts, product),
    );
  }

  function mutationSuccess(product: AtlasNutritionProduct, kind: MutationKind) {
    setCachedProduct(product);
    setFormError(null);
    setSuccessMessage(
      kind === 'create'
        ? t('nutrition.productCreated')
        : kind === 'update'
          ? t('nutrition.productUpdated')
          : kind === 'archive'
            ? t('nutrition.productArchived')
            : t('nutrition.productRestored'),
    );

    if (kind === 'create' || kind === 'update') {
      setForm(emptyForm);
      setEditingProduct(null);
    }
  }

  function mutationError(error: unknown) {
    setSuccessMessage(null);
    setFormError(errorMessageFromUnknown(error));
  }

  const createProductMutation = useMutation({
    mutationFn: (input: CreateAtlasNutritionProductInput) => createAtlasNutritionProduct(input),
    onError: mutationError,
    onSuccess: (product) => mutationSuccess(product, 'create'),
  });

  const updateProductMutation = useMutation({
    mutationFn: ({ id, input }: { id: string; input: CreateAtlasNutritionProductInput }) =>
      updateAtlasNutritionProduct(id, input),
    onError: mutationError,
    onSuccess: (product) => mutationSuccess(product, 'update'),
  });

  const archiveProductMutation = useMutation({
    mutationFn: (id: string) => archiveAtlasNutritionProduct(id),
    onError: mutationError,
    onSuccess: (product) => mutationSuccess(product, 'archive'),
  });

  const restoreProductMutation = useMutation({
    mutationFn: (id: string) => restoreAtlasNutritionProduct(id),
    onError: mutationError,
    onSuccess: (product) => mutationSuccess(product, 'restore'),
  });

  const isSaving = createProductMutation.isPending || updateProductMutation.isPending;

  function updateField(field: keyof ProductFormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function startEdit(product: AtlasNutritionProduct) {
    setEditingProduct(product);
    setForm(formFromProduct(product));
    setFormError(null);
    setSuccessMessage(null);
  }

  function cancelEdit() {
    setEditingProduct(null);
    setForm(emptyForm);
    setFormError(null);
  }

  // START_CONTRACT: handleSubmit
  //   PURPOSE: Validate and submit create/update product form data to Atlas nutrition API mutations.
  //   INPUTS: { event: FormEvent<HTMLFormElement> - form submit event }
  //   OUTPUTS: none.
  //   SIDE_EFFECTS: Prevents default form submit and starts create or update mutation.
  //   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
  // END_CONTRACT: handleSubmit
  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSuccessMessage(null);

    const parsed = parseProductForm(form, {
      nameRequired: t('nutrition.nameValidationError'),
      macroInvalid: t('nutrition.macroValidationError'),
    });

    if (parsed.error) {
      setFormError(parsed.error);
      return;
    }

    if (!parsed.input) {
      setFormError(t('nutrition.macroValidationError'));
      return;
    }

    setFormError(null);

    if (editingProduct) {
      updateProductMutation.mutate({ id: editingProduct.id, input: parsed.input });
      return;
    }

    createProductMutation.mutate(parsed.input);
  }

  function renderTable() {
    if (filteredProducts.length === 0) {
      return (
        <AdminEmptyState
          title={t('nutrition.emptyProductLibrary')}
          description={t('nutrition.emptyProductLibraryDescription')}
        />
      );
    }

    return (
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t('nutrition.productName')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.caloriesPer100g')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.proteinPer100g')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.fatPer100g')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.carbsPer100g')}</TableHead>
            <TableHead>{t('nutrition.notes')}</TableHead>
            <TableHead>{t('nutrition.status')}</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredProducts.map((product) => (
            <TableRow key={product.id}>
              <TableCell className="font-medium">{product.name}</TableCell>
              <TableCell className="atlas-macro-cell">
                {formatMacro(product.caloriesPer100g)}
              </TableCell>
              <TableCell className="atlas-macro-cell">
                {formatMacro(product.proteinPer100g)}
              </TableCell>
              <TableCell className="atlas-macro-cell">{formatMacro(product.fatPer100g)}</TableCell>
              <TableCell className="atlas-macro-cell">
                {formatMacro(product.carbsPer100g)}
              </TableCell>
              <TableCell>{product.notes || '-'}</TableCell>
              <TableCell>
                <Badge variant={product.isActive ? 'default' : 'secondary'}>
                  {product.isActive ? t('nutrition.active') : t('nutrition.archived')}
                </Badge>
              </TableCell>
              <TableCell>
                <div className="atlas-table-actions">
                  <Button onClick={() => startEdit(product)} size="sm" variant="outline">
                    {t('nutrition.editProduct')}
                  </Button>
                  {product.isActive ? (
                    <Button
                      disabled={archiveProductMutation.isPending}
                      onClick={() => archiveProductMutation.mutate(product.id)}
                      size="sm"
                      variant="outline"
                    >
                      {t('nutrition.archiveProduct')}
                    </Button>
                  ) : (
                    <Button
                      disabled={restoreProductMutation.isPending}
                      onClick={() => restoreProductMutation.mutate(product.id)}
                      size="sm"
                      variant="outline"
                    >
                      {t('nutrition.restoreProduct')}
                    </Button>
                  )}
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    );
  }

  return (
    <AdminPageShell className="atlas-product-library">
      <AdminPageHeader
        title={t('nutrition.productLibrary')}
        description={t('nutrition.productLibraryDescription')}
      />

      {successMessage ? (
        <Alert>
          <AlertTitle>{successMessage}</AlertTitle>
          <AlertDescription>{t('nutrition.includeArchived')}</AlertDescription>
        </Alert>
      ) : null}

      {formError ? (
        <Alert variant="destructive">
          <AlertTitle>{formError}</AlertTitle>
          <AlertDescription>{t('nutrition.productLibraryDescription')}</AlertDescription>
        </Alert>
      ) : null}

      <AdminSection
        title={editingProduct ? t('nutrition.editProduct') : t('nutrition.createProduct')}
        description={t('nutrition.productLibraryDescription')}
      >
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div className="atlas-form-grid">
            <div className="space-y-2">
              <Label htmlFor="product-name">{t('nutrition.productName')}</Label>
              <Input
                id="product-name"
                onChange={(event) => updateField('name', event.target.value)}
                value={form.name}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="product-calories">{t('nutrition.caloriesPer100g')}</Label>
              <Input
                id="product-calories"
                inputMode="decimal"
                onChange={(event) => updateField('caloriesPer100g', event.target.value)}
                type="number"
                value={form.caloriesPer100g}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="product-protein">{t('nutrition.proteinPer100g')}</Label>
              <Input
                id="product-protein"
                inputMode="decimal"
                onChange={(event) => updateField('proteinPer100g', event.target.value)}
                type="number"
                value={form.proteinPer100g}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="product-fat">{t('nutrition.fatPer100g')}</Label>
              <Input
                id="product-fat"
                inputMode="decimal"
                onChange={(event) => updateField('fatPer100g', event.target.value)}
                type="number"
                value={form.fatPer100g}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="product-carbs">{t('nutrition.carbsPer100g')}</Label>
              <Input
                id="product-carbs"
                inputMode="decimal"
                onChange={(event) => updateField('carbsPer100g', event.target.value)}
                type="number"
                value={form.carbsPer100g}
              />
            </div>
            <div className="space-y-2 md:col-span-2 xl:col-span-3">
              <Label htmlFor="product-notes">{t('nutrition.notes')}</Label>
              <Textarea
                id="product-notes"
                onChange={(event) => updateField('notes', event.target.value)}
                value={form.notes}
              />
            </div>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button disabled={isSaving} type="submit">
              {editingProduct ? t('nutrition.updateProduct') : t('nutrition.saveProduct')}
            </Button>
            {editingProduct ? (
              <Button onClick={cancelEdit} type="button" variant="outline">
                {t('nutrition.cancelEdit')}
              </Button>
            ) : null}
          </div>
        </form>
      </AdminSection>

      <AdminSection title={t('nutrition.products')} description={t('nutrition.includeArchived')}>
        <AdminToolbar className="atlas-filter-bar">
          <div className="space-y-2 md:min-w-80">
            <Label htmlFor="product-search">{t('nutrition.searchProducts')}</Label>
            <Input
              id="product-search"
              onChange={(event) => setSearchTerm(event.target.value)}
              value={searchTerm}
            />
          </div>
          <div aria-label={t('nutrition.status')} className="atlas-segmented-filter" role="tablist">
            <Button
              aria-selected={statusFilter === 'active'}
              onClick={() => setStatusFilter('active')}
              role="tab"
              size="sm"
              type="button"
              variant={statusFilter === 'active' ? 'default' : 'outline'}
            >
              {t('nutrition.active')}
            </Button>
            <Button
              aria-selected={statusFilter === 'archived'}
              onClick={() => setStatusFilter('archived')}
              role="tab"
              size="sm"
              type="button"
              variant={statusFilter === 'archived' ? 'default' : 'outline'}
            >
              {t('nutrition.archived')}
            </Button>
          </div>
        </AdminToolbar>

        {productsQuery.isLoading ? (
          <Card aria-label={t('nutrition.loadingProducts')} role="status">
            <CardHeader>
              <CardTitle>{t('nutrition.loadingProducts')}</CardTitle>
              <CardDescription>{t('nutrition.productLibraryDescription')}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-2">
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-10 w-full" />
            </CardContent>
          </Card>
        ) : null}

        {productsQuery.isError ? (
          <Alert variant="destructive">
            <AlertTitle>{t('nutrition.loadProductsError')}</AlertTitle>
            <AlertDescription className="space-y-3">
              <span>{errorMessageFromUnknown(productsQuery.error)}</span>
              <Button onClick={() => productsQuery.refetch()} size="sm" variant="outline">
                {t('nutrition.retry')}
              </Button>
            </AlertDescription>
          </Alert>
        ) : null}

        {productsQuery.data && !productsQuery.isError ? renderTable() : null}
      </AdminSection>
    </AdminPageShell>
  );
}
