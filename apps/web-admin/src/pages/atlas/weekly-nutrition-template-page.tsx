// FILE: apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the Atlas Weekly Plan route as an editable weekly nutrition template.
//   SCOPE: Loads one weekly template by week start, lists active products, edits template header and planned product rows, saves template-only changes, applies saved templates to empty factual days, and displays calculated planned totals; excludes product CRUD, factual daily entry editing, replace-week apply mode, body tracking, and mock/reference UI.
//   DEPENDS: @tanstack/react-query, react, react-router, apps/web-admin/src/app/i18n.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts, apps/web-admin/src/styles/atlas.css, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - API-backed Weekly Plan route content for editable weekly nutrition templates.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Hardened current-template missing state, text clearing, refetch dirty guard, and localized entry labels.
// END_CHANGE_SUMMARY

import { useQuery, useQueryClient } from '@tanstack/react-query';
import { type FormEvent, useEffect, useMemo, useRef, useState } from 'react';
import { Link } from 'react-router';
import { useI18n } from '../../app/i18n';
import '../../styles/atlas.css';
import {
  applyAtlasNutritionTemplateToWeek,
  AtlasNutritionApiError,
  createAtlasNutritionTemplate,
  createAtlasNutritionTemplateItem,
  deleteAtlasNutritionTemplateItem,
  getAtlasNutritionTemplateCurrent,
  listAtlasNutritionProducts,
  updateAtlasNutritionTemplate,
  updateAtlasNutritionTemplateItem,
  type AtlasNutritionMacros,
  type AtlasNutritionProduct,
  type AtlasNutritionTemplate,
  type AtlasNutritionTemplateApplyResult,
  type AtlasNutritionTemplateItem,
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Skeleton,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  Textarea,
} from '@shared/ui';

type WeeklyNutritionTemplatePageProps = {
  initialWeekStartDate?: string;
};

type DraftTemplateItem = {
  localId: string;
  id: string | null;
  originalProductId: string | null;
  productId: string;
  amountGrams: string;
  mealLabel: string;
  notes: string;
};

const productPlaceholderValue = '__select_product__';
const zeroMacros: AtlasNutritionMacros = {
  calories: 0,
  protein: 0,
  fat: 0,
  carbs: 0,
};

let draftItemCounter = 0;

function nextDraftItemId() {
  draftItemCounter += 1;
  return `new-template-item-${draftItemCounter}`;
}

function templateQueryKey(weekStartDate: string) {
  return ['atlas-weekly-nutrition-template', weekStartDate] as const;
}

function toDateString(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function dateFromString(dateString: string) {
  return new Date(`${dateString}T12:00:00`);
}

function addDays(dateString: string, days: number) {
  const date = dateFromString(dateString);
  date.setDate(date.getDate() + days);
  return toDateString(date);
}

function getWeekStartDate(date: Date) {
  const weekDate = new Date(date);
  const day = weekDate.getDay();
  const mondayOffset = day === 0 ? -6 : 1 - day;
  weekDate.setDate(weekDate.getDate() + mondayOffset);
  return toDateString(weekDate);
}

function getCurrentWeekStartDate() {
  return getWeekStartDate(new Date());
}

function formatLongDate(dateString: string, language: 'en' | 'ru') {
  return new Intl.DateTimeFormat(language === 'ru' ? 'ru-RU' : 'en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(dateFromString(dateString));
}

function formatWeekRange(weekStartDate: string, language: 'en' | 'ru') {
  const weekEndDate = addDays(weekStartDate, 6);
  return `${formatLongDate(weekStartDate, language)} - ${formatLongDate(weekEndDate, language)}`;
}

function formatNumber(value: number, maximumFractionDigits = 1) {
  return value.toLocaleString('en-US', {
    maximumFractionDigits,
    minimumFractionDigits: 0,
  });
}

function formatCalories(value: number) {
  return `${formatNumber(value, 0)} kcal`;
}

function formatGrams(value: number) {
  return `${formatNumber(value)} g`;
}

function errorMessageFromUnknown(error: unknown): string {
  return error instanceof Error ? error.message : 'Request failed';
}

function isNotFoundError(error: unknown) {
  return error instanceof AtlasNutritionApiError ||
    (typeof error === 'object' && error !== null && 'type' in error)
    ? (error as { type?: string }).type === 'not_found'
    : false;
}

async function getCurrentTemplateOrNull(
  weekStartDate: string,
): Promise<AtlasNutritionTemplate | null> {
  try {
    return await getAtlasNutritionTemplateCurrent(weekStartDate);
  } catch (error) {
    if (isNotFoundError(error)) {
      return null;
    }
    throw error;
  }
}

function draftFromTemplateItem(item: AtlasNutritionTemplateItem): DraftTemplateItem {
  return {
    localId: item.id,
    id: item.id,
    originalProductId: item.productId,
    productId: item.productId,
    amountGrams: String(item.amountGrams),
    mealLabel: item.mealLabel ?? '',
    notes: item.notes ?? '',
  };
}

function emptyDraftItem(): DraftTemplateItem {
  return {
    localId: nextDraftItemId(),
    id: null,
    originalProductId: null,
    productId: '',
    amountGrams: '',
    mealLabel: '',
    notes: '',
  };
}

function parsePositiveGrams(value: string) {
  const grams = Number(value);
  return Number.isFinite(grams) && grams > 0 ? grams : null;
}

function optionalCreateText(value: string) {
  return value.trim() || null;
}

function updateText(value: string) {
  return value.trim();
}

function productName(
  product: AtlasNutritionProduct | undefined,
  productId: string,
  unknownProductLabel: string,
  plannedEntryLabel: string,
) {
  return product?.name ?? (productId ? `${unknownProductLabel} (${productId})` : plannedEntryLabel);
}

function calculateItemMacros(
  item: DraftTemplateItem,
  product: AtlasNutritionProduct | undefined,
): AtlasNutritionMacros {
  const amountGrams = parsePositiveGrams(item.amountGrams) ?? 0;

  if (!product) {
    return zeroMacros;
  }

  return {
    calories: (product.caloriesPer100g * amountGrams) / 100,
    protein: (product.proteinPer100g * amountGrams) / 100,
    fat: (product.fatPer100g * amountGrams) / 100,
    carbs: (product.carbsPer100g * amountGrams) / 100,
  };
}

function addMacros(left: AtlasNutritionMacros, right: AtlasNutritionMacros): AtlasNutritionMacros {
  return {
    calories: left.calories + right.calories,
    protein: left.protein + right.protein,
    fat: left.fat + right.fat,
    carbs: left.carbs + right.carbs,
  };
}

function statusCount(result: AtlasNutritionTemplateApplyResult, status: string) {
  return result.dates.filter((date) => date.status.toLowerCase() === status).length;
}

// START_CONTRACT: WeeklyNutritionTemplatePage
//   PURPOSE: Render and mutate one weekly nutrition template while keeping factual day seeding as an explicit separate action.
//   INPUTS: { initialWeekStartDate?: string - optional test/default selected week start in YYYY-MM-DD format }
//   OUTPUTS: { JSX.Element - weekly template editor with loading, error, empty, validation, success, apply result, and no body-weight states }
//   SIDE_EFFECTS: Sends template create/update/delete item requests on Save Template and sends seed_empty_days apply only from Apply to Week.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
// END_CONTRACT: WeeklyNutritionTemplatePage
export default function WeeklyNutritionTemplatePage({
  initialWeekStartDate,
}: WeeklyNutritionTemplatePageProps) {
  const { language, t } = useI18n();
  const queryClient = useQueryClient();
  const initializedWeekStartDate = useRef<string | null>(null);
  const [selectedWeekStartDate, setSelectedWeekStartDate] = useState(
    initialWeekStartDate ?? getCurrentWeekStartDate(),
  );
  const [title, setTitle] = useState('');
  const [notes, setNotes] = useState('');
  const [draftItems, setDraftItems] = useState<DraftTemplateItem[]>([]);
  const [deletedItemIds, setDeletedItemIds] = useState<string[]>([]);
  const [editorError, setEditorError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [applyResult, setApplyResult] = useState<AtlasNutritionTemplateApplyResult | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [isApplying, setIsApplying] = useState(false);
  const [isDirty, setIsDirty] = useState(false);

  const templateQuery = useQuery({
    queryKey: templateQueryKey(selectedWeekStartDate),
    queryFn: () => getCurrentTemplateOrNull(selectedWeekStartDate),
  });

  const productsQuery = useQuery({
    queryKey: ['atlas-nutrition-products', 'active'],
    queryFn: () => listAtlasNutritionProducts(),
  });

  const template = templateQuery.data ?? null;
  const products = productsQuery.data ?? [];
  const productsById = useMemo(
    () => new Map(products.map((product) => [product.id, product])),
    [products],
  );
  const isLoading = templateQuery.isLoading || productsQuery.isLoading;
  const hasLoadError = templateQuery.isError || productsQuery.isError;
  const weekLabel = `${t('nutrition.weekOf')} ${formatWeekRange(selectedWeekStartDate, language)}`;
  const totals = useMemo(
    () =>
      draftItems.reduce(
        (currentTotals, item) =>
          addMacros(currentTotals, calculateItemMacros(item, productsById.get(item.productId))),
        zeroMacros,
      ),
    [draftItems, productsById],
  );

  const summaryCards = useMemo(
    () => [
      { label: t('nutrition.calories'), value: formatCalories(totals.calories) },
      { label: t('nutrition.protein'), value: formatGrams(totals.protein) },
      { label: t('nutrition.fat'), value: formatGrams(totals.fat) },
      { label: t('nutrition.carbs'), value: formatGrams(totals.carbs) },
    ],
    [t, totals],
  );

  useEffect(() => {
    if (templateQuery.data === undefined) {
      return;
    }

    if (initializedWeekStartDate.current === selectedWeekStartDate && isDirty) {
      return;
    }

    initializedWeekStartDate.current = selectedWeekStartDate;
    setTitle(templateQuery.data?.title ?? '');
    setNotes(templateQuery.data?.notes ?? '');
    setDraftItems((templateQuery.data?.items ?? []).map(draftFromTemplateItem));
    setDeletedItemIds([]);
    setEditorError(null);
    setIsDirty(false);
  }, [isDirty, selectedWeekStartDate, templateQuery.data]);

  function changeWeek(nextWeekStartDate: string) {
    setSelectedWeekStartDate(nextWeekStartDate);
    setEditorError(null);
    setSuccessMessage(null);
    setApplyResult(null);
    setIsDirty(false);
  }

  function updateDraftItem(localId: string, field: keyof DraftTemplateItem, value: string) {
    setDraftItems((currentItems) =>
      currentItems.map((item) => (item.localId === localId ? { ...item, [field]: value } : item)),
    );
    setEditorError(null);
    setIsDirty(true);
  }

  function addDraftItem() {
    setDraftItems((currentItems) => [...currentItems, emptyDraftItem()]);
    setEditorError(null);
    setSuccessMessage(null);
    setIsDirty(true);
  }

  function deleteDraftItem(item: DraftTemplateItem) {
    setDraftItems((currentItems) =>
      currentItems.filter((currentItem) => currentItem.localId !== item.localId),
    );
    if (item.id) {
      setDeletedItemIds((currentIds) => [...currentIds, item.id as string]);
    }
    setEditorError(null);
    setSuccessMessage(null);
    setIsDirty(true);
  }

  function validateDraftItems() {
    for (const item of draftItems) {
      if (!item.productId) {
        return t('nutrition.chooseProduct');
      }
      if (parsePositiveGrams(item.amountGrams) === null) {
        return t('nutrition.gramsPositive');
      }
    }

    return null;
  }

  function createItemInput(item: DraftTemplateItem) {
    return {
      productId: item.productId,
      amountGrams: parsePositiveGrams(item.amountGrams) ?? 0,
      mealLabel: optionalCreateText(item.mealLabel),
      notes: optionalCreateText(item.notes),
    };
  }

  function updateItemInput(item: DraftTemplateItem) {
    return {
      amountGrams: parsePositiveGrams(item.amountGrams) ?? 0,
      mealLabel: updateText(item.mealLabel),
      notes: updateText(item.notes),
    };
  }

  // START_CONTRACT: handleSaveTemplate
  //   PURPOSE: Persist template header and draft item CRUD changes without applying them to factual daily logs.
  //   INPUTS: { event: FormEvent<HTMLFormElement> - editor submit event }
  //   OUTPUTS: none.
  //   SIDE_EFFECTS: Prevents default form submit, sends template header/item mutations, updates React Query cache, and never calls applyAtlasNutritionTemplateToWeek.
  //   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
  // END_CONTRACT: handleSaveTemplate
  async function handleSaveTemplate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSuccessMessage(null);
    setApplyResult(null);

    const validationError = validateDraftItems();
    if (validationError) {
      setEditorError(validationError);
      return;
    }

    setIsSaving(true);
    setEditorError(null);

    try {
      const savedTemplate = template
        ? await updateAtlasNutritionTemplate(template.id, {
            title: updateText(title),
            notes: updateText(notes),
          })
        : await createAtlasNutritionTemplate({
            weekStartDate: selectedWeekStartDate,
            title: optionalCreateText(title),
            notes: optionalCreateText(notes),
          });

      const savedItems: AtlasNutritionTemplateItem[] = [];

      const itemIdsToDelete = [...deletedItemIds];

      for (const item of draftItems) {
        if (item.id && item.productId === item.originalProductId) {
          savedItems.push(await updateAtlasNutritionTemplateItem(item.id, updateItemInput(item)));
        } else {
          if (item.id) {
            itemIdsToDelete.push(item.id);
          }
          savedItems.push(
            await createAtlasNutritionTemplateItem({
              templateId: savedTemplate.id,
              ...createItemInput(item),
            }),
          );
        }
      }

      for (const deletedItemId of [...new Set(itemIdsToDelete)]) {
        await deleteAtlasNutritionTemplateItem(deletedItemId);
      }

      const nextTemplate = { ...savedTemplate, items: savedItems };
      queryClient.setQueryData<AtlasNutritionTemplate | null>(
        templateQueryKey(selectedWeekStartDate),
        nextTemplate,
      );
      setDraftItems(savedItems.map(draftFromTemplateItem));
      setDeletedItemIds([]);
      setIsDirty(false);
      setSuccessMessage(t('nutrition.templateSaved'));
    } catch (error) {
      setEditorError(errorMessageFromUnknown(error));
    } finally {
      setIsSaving(false);
    }
  }

  async function handleApplyTemplate() {
    setSuccessMessage(null);
    setApplyResult(null);

    if (!template?.id) {
      setEditorError(t('nutrition.saveBeforeApply'));
      return;
    }

    setIsApplying(true);
    setEditorError(null);

    try {
      const result = await applyAtlasNutritionTemplateToWeek(template.id, 'SEED_EMPTY_DAYS');
      setApplyResult(result);
      setSuccessMessage(t('nutrition.templateApplied'));
    } catch (error) {
      setEditorError(errorMessageFromUnknown(error));
    } finally {
      setIsApplying(false);
    }
  }

  function renderSummaryCards() {
    return (
      <div className="atlas-summary-grid">
        {summaryCards.map((card) => (
          <Card key={card.label}>
            <CardHeader className="pb-2">
              <CardDescription>{card.label}</CardDescription>
              <CardTitle>{card.value}</CardTitle>
            </CardHeader>
          </Card>
        ))}
      </div>
    );
  }

  function renderProductOption(product: AtlasNutritionProduct) {
    return (
      <SelectItem key={product.id} value={product.id}>
        {product.name} ({formatCalories(product.caloriesPer100g)} / 100 g)
      </SelectItem>
    );
  }

  function renderDraftRows() {
    return draftItems.map((item, index) => {
      const entryNumber = index + 1;
      const product = productsById.get(item.productId);
      const entryContext = `${t('nutrition.forEntry')} ${entryNumber}`;
      const name = productName(
        product,
        item.productId,
        t('nutrition.unknownProduct'),
        t('nutrition.plannedEntry'),
      );
      const macros = calculateItemMacros(item, product);

      return (
        <TableRow key={item.localId}>
          <TableCell>
            <Select
              disabled={productsQuery.isLoading}
              onValueChange={(value) =>
                updateDraftItem(
                  item.localId,
                  'productId',
                  value === productPlaceholderValue ? '' : value,
                )
              }
              value={item.productId || productPlaceholderValue}
            >
              <SelectTrigger
                aria-label={`${t('nutrition.product')} ${entryContext}`}
                id={`template-product-${item.localId}`}
              >
                <SelectValue placeholder={t('nutrition.selectProduct')} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={productPlaceholderValue}>
                  {t('nutrition.selectProduct')}
                </SelectItem>
                {item.productId && !product ? (
                  <SelectItem value={item.productId}>{name}</SelectItem>
                ) : null}
                {products.map(renderProductOption)}
              </SelectContent>
            </Select>
            <div className="mt-2 text-xs text-muted-foreground">{name}</div>
          </TableCell>
          <TableCell>
            <Input
              aria-label={`${t('nutrition.grams')} ${entryContext}`}
              inputMode="decimal"
              onChange={(event) => updateDraftItem(item.localId, 'amountGrams', event.target.value)}
              type="number"
              value={item.amountGrams}
            />
          </TableCell>
          <TableCell>
            <Input
              aria-label={`${t('nutrition.mealLabel')} ${entryContext}`}
              onChange={(event) => updateDraftItem(item.localId, 'mealLabel', event.target.value)}
              value={item.mealLabel}
            />
          </TableCell>
          <TableCell>
            <Input
              aria-label={`${t('nutrition.notes')} ${entryContext}`}
              onChange={(event) => updateDraftItem(item.localId, 'notes', event.target.value)}
              value={item.notes}
            />
          </TableCell>
          <TableCell className="atlas-macro-cell">{formatNumber(macros.calories, 0)}</TableCell>
          <TableCell className="atlas-macro-cell">{formatGrams(macros.protein)}</TableCell>
          <TableCell className="atlas-macro-cell">{formatGrams(macros.fat)}</TableCell>
          <TableCell className="atlas-macro-cell">{formatGrams(macros.carbs)}</TableCell>
          <TableCell>
            <div className="atlas-table-actions">
              <Button
                aria-label={`${t('nutrition.deleteEntry')} ${name}`}
                onClick={() => deleteDraftItem(item)}
                size="sm"
                type="button"
                variant="outline"
              >
                {t('nutrition.deleteEntry')}
              </Button>
            </div>
          </TableCell>
        </TableRow>
      );
    });
  }

  function renderEntriesEditor() {
    if (draftItems.length === 0) {
      return (
        <AdminEmptyState
          title={t('nutrition.emptyWeeklyPlanTitle')}
          description={t('nutrition.emptyWeeklyPlanDescription')}
        />
      );
    }

    return (
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t('nutrition.product')}</TableHead>
            <TableHead>{t('nutrition.grams')}</TableHead>
            <TableHead>{t('nutrition.mealLabel')}</TableHead>
            <TableHead>{t('nutrition.notes')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.calories')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.protein')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.fat')}</TableHead>
            <TableHead className="atlas-macro-cell">{t('nutrition.carbs')}</TableHead>
            <TableHead className="text-right">{t('nutrition.actions')}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>{renderDraftRows()}</TableBody>
      </Table>
    );
  }

  function renderApplyResult() {
    if (!applyResult) {
      return null;
    }

    const createdCount = statusCount(applyResult, 'created');
    const skippedCount = statusCount(applyResult, 'skipped');

    return (
      <AdminSection title={t('nutrition.applyResult')} description={weekLabel}>
        <div className="mb-4 flex flex-wrap gap-2">
          <Badge variant="default">
            {createdCount} {t('nutrition.createdCount')}
          </Badge>
          <Badge variant="secondary">
            {skippedCount} {t('nutrition.skippedCount')}
          </Badge>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('nutrition.weekOf')}</TableHead>
              <TableHead>{t('nutrition.status')}</TableHead>
              <TableHead>{t('nutrition.entryCount')}</TableHead>
              <TableHead>{t('nutrition.notes')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {applyResult.dates.map((dateResult) => (
              <TableRow key={dateResult.date}>
                <TableCell>{dateResult.date}</TableCell>
                <TableCell>{dateResult.status}</TableCell>
                <TableCell>
                  {dateResult.entryCount} {t('nutrition.entryCount')}
                </TableCell>
                <TableCell>{dateResult.reason ?? '-'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </AdminSection>
    );
  }

  return (
    <AdminPageShell className="atlas-weekly-template">
      <AdminPageHeader
        actions={
          <Button asChild variant="outline">
            <Link to="/atlas/nutrition/products">{t('nutrition.manageProducts')}</Link>
          </Button>
        }
        description={t('nutrition.weeklyPlanDescription')}
        title={t('nutrition.weeklyPlan')}
      />

      <AdminToolbar className="atlas-date-switcher">
        <div className="atlas-date-controls">
          <Button
            aria-label={t('nutrition.previousWeek')}
            onClick={() => changeWeek(addDays(selectedWeekStartDate, -7))}
            type="button"
            variant="outline"
          >
            {t('nutrition.previousWeek')}
          </Button>
          <div className="atlas-selected-date">{weekLabel}</div>
          <Button
            aria-label={t('nutrition.nextWeek')}
            onClick={() => changeWeek(addDays(selectedWeekStartDate, 7))}
            type="button"
            variant="outline"
          >
            {t('nutrition.nextWeek')}
          </Button>
        </div>
        <Button
          onClick={() => changeWeek(getCurrentWeekStartDate())}
          type="button"
          variant="secondary"
        >
          {t('nutrition.today')}
        </Button>
      </AdminToolbar>

      <AdminSection title={t('nutrition.weeklyTotals')} description={weekLabel}>
        {renderSummaryCards()}
      </AdminSection>

      {successMessage ? (
        <Alert>
          <AlertTitle>{successMessage}</AlertTitle>
          <AlertDescription>{weekLabel}</AlertDescription>
        </Alert>
      ) : null}

      {editorError ? (
        <Alert variant="destructive">
          <AlertTitle>{editorError}</AlertTitle>
          <AlertDescription>{t('nutrition.weeklyPlanDescription')}</AlertDescription>
        </Alert>
      ) : null}

      {isLoading ? (
        <Card aria-label={t('nutrition.loadingWeeklyPlan')} role="status">
          <CardHeader>
            <CardTitle>{t('nutrition.loadingWeeklyPlan')}</CardTitle>
            <CardDescription>{t('nutrition.weeklyPlanDescription')}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
          </CardContent>
        </Card>
      ) : null}

      {hasLoadError ? (
        <Alert variant="destructive">
          <AlertTitle>
            {productsQuery.isError
              ? t('nutrition.loadProductsError')
              : t('nutrition.loadWeeklyPlanError')}
          </AlertTitle>
          <AlertDescription className="space-y-3">
            <span>{errorMessageFromUnknown(productsQuery.error ?? templateQuery.error)}</span>
            <Button
              onClick={() => {
                if (templateQuery.isError) {
                  templateQuery.refetch();
                }
                if (productsQuery.isError) {
                  productsQuery.refetch();
                }
              }}
              size="sm"
              type="button"
              variant="outline"
            >
              {t('nutrition.retry')}
            </Button>
          </AlertDescription>
        </Alert>
      ) : null}

      {!isLoading && !hasLoadError ? (
        <form className="space-y-6" onSubmit={handleSaveTemplate}>
          <AdminSection title={t('nutrition.weeklyPlan')} description={weekLabel}>
            <div className="atlas-template-header-grid">
              <div className="space-y-2">
                <Label htmlFor="template-title">{t('nutrition.templateTitle')}</Label>
                <Input
                  id="template-title"
                  onChange={(event) => {
                    setTitle(event.target.value);
                    setEditorError(null);
                    setIsDirty(true);
                  }}
                  value={title}
                />
              </div>
              <div className="space-y-2 md:col-span-2">
                <Label htmlFor="template-notes">{t('nutrition.templateNotes')}</Label>
                <Textarea
                  id="template-notes"
                  onChange={(event) => {
                    setNotes(event.target.value);
                    setEditorError(null);
                    setIsDirty(true);
                  }}
                  value={notes}
                />
              </div>
            </div>
          </AdminSection>

          <AdminSection
            title={t('nutrition.entries')}
            description={t('nutrition.weeklyPlanDescription')}
          >
            <AdminToolbar className="atlas-filter-bar">
              <Button onClick={addDraftItem} type="button" variant="outline">
                {t('nutrition.addPlannedEntry')}
              </Button>
              <div className="flex flex-wrap gap-2">
                <Button disabled={isSaving} type="submit">
                  {t('nutrition.saveTemplate')}
                </Button>
                <Button
                  disabled={isSaving || isApplying || !template?.id}
                  onClick={handleApplyTemplate}
                  type="button"
                  variant="secondary"
                >
                  {t('nutrition.applyToWeek')}
                </Button>
              </div>
            </AdminToolbar>
            {renderEntriesEditor()}
          </AdminSection>
        </form>
      ) : null}

      {renderApplyResult()}
    </AdminPageShell>
  );
}
