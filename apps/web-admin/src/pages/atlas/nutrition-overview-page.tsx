// FILE: apps/web-admin/src/pages/atlas/nutrition-overview-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the Atlas Nutrition route as an API-backed factual daily food log.
//   SCOPE: Loads one Daily Nutrition Log by date, lists active products, adds/edits/deletes product gram entries, and displays calculated totals; excludes product CRUD, weekly plan editing, body tracking, and legacy override editing.
//   DEPENDS: @tanstack/react-query, react, react-router, apps/web-admin/src/app/i18n.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts, apps/web-admin/src/styles/atlas.css, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - API-backed Nutrition route content for one factual Daily Nutrition Log.
// END_MODULE_MAP

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { type FormEvent, useMemo, useState } from 'react';
import { Link } from 'react-router';
import { useI18n } from '../../app/i18n';
import '../../styles/atlas.css';
import {
  addAtlasDailyNutritionEntry,
  deleteAtlasDailyNutritionEntry,
  getAtlasDailyNutritionLog,
  listAtlasNutritionProducts,
  updateAtlasDailyNutritionEntry,
  type AtlasDailyNutritionEntry,
  type AtlasDailyNutritionLog,
  type AtlasNutritionMacros,
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

type EntryFormState = {
  productId: string;
  amountGrams: string;
  mealLabel: string;
  notes: string;
};

type NutritionOverviewPageProps = {
  initialDate?: string;
};

type SuccessMessage = {
  date: string;
  message: string;
};

const productPlaceholderValue = '__select_product__';

const emptyEntryForm: EntryFormState = {
  productId: '',
  amountGrams: '',
  mealLabel: '',
  notes: '',
};

const zeroMacros: AtlasNutritionMacros = {
  calories: 0,
  protein: 0,
  fat: 0,
  carbs: 0,
};

function dailyLogQueryKey(date: string) {
  return ['atlas-daily-nutrition-log', date] as const;
}

function getTodayDateString() {
  return toDateString(new Date());
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

function formatLongDate(dateString: string, language: 'en' | 'ru') {
  return new Intl.DateTimeFormat(language === 'ru' ? 'ru-RU' : 'en-US', {
    day: 'numeric',
    month: 'long',
    weekday: 'long',
    year: 'numeric',
  }).format(dateFromString(dateString));
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

function formFromEntry(entry: AtlasDailyNutritionEntry): EntryFormState {
  return {
    productId: entry.productId,
    amountGrams: String(entry.amountGrams),
    mealLabel: entry.mealLabel ?? '',
    notes: entry.notes ?? '',
  };
}

function parsePositiveGrams(value: string) {
  const grams = Number(value);
  return Number.isFinite(grams) && grams > 0 ? grams : null;
}

// START_CONTRACT: NutritionOverviewPage
//   PURPOSE: Render and mutate one factual Daily Nutrition Log from Atlas nutrition API data.
//   INPUTS: { initialDate?: string - optional test/default selected date in YYYY-MM-DD format }
//   OUTPUTS: { JSX.Element - nutrition page with date switcher, totals, add-entry form, entry table, and page states }
//   SIDE_EFFECTS: Sends daily nutrition add/update/delete requests and updates React Query cache on success.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
// END_CONTRACT: NutritionOverviewPage
export default function NutritionOverviewPage({ initialDate }: NutritionOverviewPageProps) {
  const { language, t } = useI18n();
  const queryClient = useQueryClient();
  const [selectedDate, setSelectedDate] = useState(initialDate ?? getTodayDateString());
  const [form, setForm] = useState<EntryFormState>(emptyEntryForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<SuccessMessage | null>(null);
  const [editingEntryId, setEditingEntryId] = useState<string | null>(null);
  const [editForm, setEditForm] = useState<EntryFormState>(emptyEntryForm);
  const [rowError, setRowError] = useState<{ entryId: string; message: string } | null>(null);

  const dailyLogQuery = useQuery({
    queryKey: dailyLogQueryKey(selectedDate),
    queryFn: () => getAtlasDailyNutritionLog(selectedDate),
  });

  const productsQuery = useQuery({
    queryKey: ['atlas-nutrition-products', 'active'],
    queryFn: () => listAtlasNutritionProducts(),
  });

  const dailyLog = dailyLogQuery.data;
  const totals = dailyLog?.totals ?? zeroMacros;
  const products = productsQuery.data ?? [];
  const selectedDateLabel = formatLongDate(selectedDate, language);
  const hasLoadError = dailyLogQuery.isError || productsQuery.isError;
  const isLoading = dailyLogQuery.isLoading || productsQuery.isLoading;
  const visibleSuccessMessage =
    successMessage?.date === selectedDate ? successMessage.message : null;

  const summaryCards = useMemo(
    () => [
      { label: t('nutrition.calories'), value: formatCalories(totals.calories) },
      { label: t('nutrition.protein'), value: formatGrams(totals.protein) },
      { label: t('nutrition.fat'), value: formatGrams(totals.fat) },
      { label: t('nutrition.carbs'), value: formatGrams(totals.carbs) },
    ],
    [t, totals],
  );

  function setDailyLogCache(log: AtlasDailyNutritionLog) {
    queryClient.setQueryData(dailyLogQueryKey(log.date), log);
  }

  function changeDate(nextDate: string) {
    setSelectedDate(nextDate);
    setFormError(null);
    setSuccessMessage(null);
    setEditingEntryId(null);
    setRowError(null);
  }

  function mutationSuccess(log: AtlasDailyNutritionLog, message: string) {
    setDailyLogCache(log);
    setSuccessMessage({ date: log.date, message });
    setFormError(null);
    setRowError(null);
  }

  const addEntryMutation = useMutation({
    mutationFn: (input: {
      date: string;
      productId: string;
      amountGrams: number;
      mealLabel: string | null;
      notes: string | null;
    }) => addAtlasDailyNutritionEntry(input),
    onError: (error) => {
      setSuccessMessage(null);
      setFormError(errorMessageFromUnknown(error));
    },
    onSuccess: (log) => {
      mutationSuccess(log, t('nutrition.entryAdded'));
      setForm(emptyEntryForm);
    },
  });

  const updateEntryMutation = useMutation({
    mutationFn: ({
      entry,
      input,
    }: {
      entry: AtlasDailyNutritionEntry;
      input: {
        dailyLogId: string;
        amountGrams: number;
        mealLabel: string | null;
        notes: string | null;
        position: number;
      };
    }) => updateAtlasDailyNutritionEntry(entry.id, input),
    onError: (error, variables) => {
      setSuccessMessage(null);
      setRowError({ entryId: variables.entry.id, message: errorMessageFromUnknown(error) });
    },
    onSuccess: (log) => {
      mutationSuccess(log, t('nutrition.entryUpdated'));
      setEditingEntryId(null);
      setEditForm(emptyEntryForm);
    },
  });

  const deleteEntryMutation = useMutation({
    mutationFn: (entry: AtlasDailyNutritionEntry) => deleteAtlasDailyNutritionEntry(entry.id),
    onError: (error, entry) => {
      setSuccessMessage(null);
      setRowError({ entryId: entry.id, message: errorMessageFromUnknown(error) });
    },
    onSuccess: (log) => mutationSuccess(log, t('nutrition.entryDeleted')),
  });

  function updateForm(field: keyof EntryFormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
    setFormError(null);
  }

  function updateEditForm(field: keyof EntryFormState, value: string) {
    setEditForm((current) => ({ ...current, [field]: value }));
    setRowError(null);
  }

  function handleAddEntry(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!form.productId) {
      setFormError(t('nutrition.chooseProduct'));
      return;
    }

    const amountGrams = parsePositiveGrams(form.amountGrams);
    if (amountGrams === null) {
      setFormError(t('nutrition.gramsPositive'));
      return;
    }

    addEntryMutation.mutate({
      date: selectedDate,
      productId: form.productId,
      amountGrams,
      mealLabel: form.mealLabel.trim() || null,
      notes: form.notes.trim() || null,
    });
  }

  function startEdit(entry: AtlasDailyNutritionEntry) {
    setEditingEntryId(entry.id);
    setEditForm(formFromEntry(entry));
    setRowError(null);
    setSuccessMessage(null);
  }

  function saveEdit(entry: AtlasDailyNutritionEntry) {
    const amountGrams = parsePositiveGrams(editForm.amountGrams);
    if (amountGrams === null) {
      setRowError({ entryId: entry.id, message: t('nutrition.gramsPositive') });
      return;
    }

    updateEntryMutation.mutate({
      entry,
      input: {
        dailyLogId: entry.dailyLogId,
        amountGrams,
        mealLabel: editForm.mealLabel.trim() || null,
        notes: editForm.notes.trim() || null,
        position: entry.position,
      },
    });
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

  function renderEntryForm() {
    return (
      <AdminSection title={t('nutrition.addFood')} description={t('nutrition.dailyLogDescription')}>
        <form className="space-y-4" onSubmit={handleAddEntry}>
          {formError ? (
            <Alert variant="destructive">
              <AlertTitle>{formError}</AlertTitle>
            </Alert>
          ) : null}
          <div className="atlas-entry-form">
            <div className="space-y-2">
              <Label htmlFor="nutrition-product">{t('nutrition.product')}</Label>
              <Select
                disabled={productsQuery.isLoading}
                onValueChange={(value) =>
                  updateForm('productId', value === productPlaceholderValue ? '' : value)
                }
                value={form.productId || productPlaceholderValue}
              >
                <SelectTrigger aria-label={t('nutrition.product')} id="nutrition-product">
                  <SelectValue placeholder={t('nutrition.selectProduct')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value={productPlaceholderValue}>
                    {t('nutrition.selectProduct')}
                  </SelectItem>
                  {products.map((product) => (
                    <SelectItem key={product.id} value={product.id}>
                      {product.name} ({formatCalories(product.caloriesPer100g)} / 100 g)
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="nutrition-grams">{t('nutrition.grams')}</Label>
              <Input
                id="nutrition-grams"
                inputMode="decimal"
                onChange={(event) => updateForm('amountGrams', event.target.value)}
                type="number"
                value={form.amountGrams}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="nutrition-meal">{t('nutrition.mealLabel')}</Label>
              <Input
                id="nutrition-meal"
                onChange={(event) => updateForm('mealLabel', event.target.value)}
                value={form.mealLabel}
              />
            </div>
            <div className="space-y-2 md:col-span-2 xl:col-span-3">
              <Label htmlFor="nutrition-entry-notes">{t('nutrition.entryNotes')}</Label>
              <Textarea
                id="nutrition-entry-notes"
                onChange={(event) => updateForm('notes', event.target.value)}
                value={form.notes}
              />
            </div>
          </div>
          <Button disabled={addEntryMutation.isPending} type="submit">
            {t('nutrition.addFood')}
          </Button>
        </form>
      </AdminSection>
    );
  }

  function renderEntryRows(entries: AtlasDailyNutritionEntry[]) {
    return entries.map((entry) => {
      const isEditing = editingEntryId === entry.id;
      const rowMessage = rowError?.entryId === entry.id ? rowError.message : null;
      const productName = entry.productNameSnapshot;

      return (
        <TableRow key={entry.id}>
          <TableCell>
            <div className="font-medium">{productName}</div>
            <div className="text-xs text-muted-foreground">
              {formatNumber(entry.caloriesPer100gSnapshot, 0)} kcal / 100 g
            </div>
            {rowMessage ? <div className="atlas-row-error">{rowMessage}</div> : null}
          </TableCell>
          <TableCell>
            {isEditing ? (
              <Input
                aria-label={`${t('nutrition.grams')} for ${productName}`}
                inputMode="decimal"
                onChange={(event) => updateEditForm('amountGrams', event.target.value)}
                type="number"
                value={editForm.amountGrams}
              />
            ) : (
              formatGrams(entry.amountGrams)
            )}
          </TableCell>
          <TableCell>
            {isEditing ? (
              <Input
                aria-label={`${t('nutrition.mealLabel')} for ${productName}`}
                onChange={(event) => updateEditForm('mealLabel', event.target.value)}
                value={editForm.mealLabel}
              />
            ) : (
              (entry.mealLabel ?? '—')
            )}
          </TableCell>
          <TableCell>
            {isEditing ? (
              <Input
                aria-label={`${t('nutrition.notes')} for ${productName}`}
                onChange={(event) => updateEditForm('notes', event.target.value)}
                value={editForm.notes}
              />
            ) : (
              (entry.notes ?? '—')
            )}
          </TableCell>
          <TableCell className="atlas-macro-cell">
            {formatNumber(entry.macros.calories, 0)}
          </TableCell>
          <TableCell className="atlas-macro-cell">{formatGrams(entry.macros.protein)}</TableCell>
          <TableCell className="atlas-macro-cell">{formatGrams(entry.macros.fat)}</TableCell>
          <TableCell className="atlas-macro-cell">{formatGrams(entry.macros.carbs)}</TableCell>
          <TableCell>
            <div className="atlas-table-actions">
              {isEditing ? (
                <>
                  <Button
                    aria-label={`${t('nutrition.saveEntry')} ${productName}`}
                    disabled={updateEntryMutation.isPending}
                    onClick={() => saveEdit(entry)}
                    size="sm"
                    type="button"
                  >
                    {t('nutrition.saveEntry')}
                  </Button>
                  <Button
                    aria-label={`${t('nutrition.cancelEntryEdit')} ${productName}`}
                    onClick={() => setEditingEntryId(null)}
                    size="sm"
                    type="button"
                    variant="outline"
                  >
                    {t('nutrition.cancelEntryEdit')}
                  </Button>
                </>
              ) : (
                <>
                  <Button
                    aria-label={`${t('nutrition.editProduct')} ${productName}`}
                    onClick={() => startEdit(entry)}
                    size="sm"
                    type="button"
                    variant="outline"
                  >
                    {t('nutrition.editProduct')}
                  </Button>
                  <Button
                    aria-label={`${t('nutrition.deleteEntry')} ${productName}`}
                    disabled={deleteEntryMutation.isPending}
                    onClick={() => deleteEntryMutation.mutate(entry)}
                    size="sm"
                    type="button"
                    variant="outline"
                  >
                    {t('nutrition.deleteEntry')}
                  </Button>
                </>
              )}
            </div>
          </TableCell>
        </TableRow>
      );
    });
  }

  function renderEntries(log: AtlasDailyNutritionLog) {
    if (log.entries.length === 0) {
      return (
        <AdminEmptyState
          title={t('nutrition.emptyDayTitle')}
          description={t('nutrition.emptyDayDescription')}
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
        <TableBody>{renderEntryRows(log.entries)}</TableBody>
      </Table>
    );
  }

  return (
    <AdminPageShell className="atlas-nutrition-log">
      <AdminPageHeader
        actions={
          <Button asChild variant="outline">
            <Link to="/atlas/nutrition/products">{t('nutrition.manageProducts')}</Link>
          </Button>
        }
        description={t('nutrition.dailyLogDescription')}
        title={t('nutrition.title')}
      />

      <AdminToolbar className="atlas-date-switcher">
        <div className="atlas-date-controls">
          <Button
            aria-label={t('nutrition.previousDay')}
            onClick={() => changeDate(addDays(selectedDate, -1))}
            type="button"
            variant="outline"
          >
            {t('nutrition.previousDay')}
          </Button>
          <div className="atlas-selected-date">{selectedDateLabel}</div>
          <Button
            aria-label={t('nutrition.nextDay')}
            onClick={() => changeDate(addDays(selectedDate, 1))}
            type="button"
            variant="outline"
          >
            {t('nutrition.nextDay')}
          </Button>
        </div>
        <Button onClick={() => changeDate(getTodayDateString())} type="button" variant="secondary">
          {t('nutrition.today')}
        </Button>
      </AdminToolbar>

      {renderSummaryCards()}

      {visibleSuccessMessage ? (
        <Alert>
          <AlertTitle>{visibleSuccessMessage}</AlertTitle>
        </Alert>
      ) : null}

      {isLoading ? (
        <Card aria-label={t('nutrition.loadingNutrition')} role="status">
          <CardHeader>
            <CardTitle>{t('nutrition.loadingNutrition')}</CardTitle>
            <CardDescription>{t('nutrition.dailyLogDescription')}</CardDescription>
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
          <AlertTitle>{t('nutrition.loadNutritionError')}</AlertTitle>
          <AlertDescription className="space-y-3">
            <span>{errorMessageFromUnknown(dailyLogQuery.error ?? productsQuery.error)}</span>
            <Button
              onClick={() => {
                if (dailyLogQuery.isError) {
                  dailyLogQuery.refetch();
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
        <>
          {renderEntryForm()}
          <AdminSection
            title={t('nutrition.entries')}
            description={`${t('nutrition.entries')} · ${selectedDateLabel}`}
          >
            <div className="mb-4 flex flex-wrap gap-2">
              <Badge variant="secondary">
                {t('nutrition.calories')}: {formatCalories(totals.calories)}
              </Badge>
              <Badge variant="outline">
                {t('nutrition.protein')}: {formatGrams(totals.protein)}
              </Badge>
              <Badge variant="outline">
                {t('nutrition.fat')}: {formatGrams(totals.fat)}
              </Badge>
              <Badge variant="outline">
                {t('nutrition.carbs')}: {formatGrams(totals.carbs)}
              </Badge>
            </div>
            {dailyLog ? renderEntries(dailyLog) : null}
          </AdminSection>
        </>
      ) : null}
    </AdminPageShell>
  );
}
