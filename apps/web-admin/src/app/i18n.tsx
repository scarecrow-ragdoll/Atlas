// FILE: apps/web-admin/src/app/i18n.tsx
// VERSION: 1.2.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide lightweight EN/RU translation state for Atlas web-admin pages.
//   SCOPE: Owns local language persistence, document language sync, translation lookup, and provider/hook exports; excludes server-side locale negotiation and settings-page UI.
//   DEPENDS: react.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   I18nProvider - Provides selected language and translation lookup to child routes.
//   useI18n - Reads current language helpers from the provider.
//   LANGUAGE_STORAGE_KEY - LocalStorage key for the selected Atlas language.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.2.1 - Added localized weekly template entry fallback and aria label fragments.
// END_CHANGE_SUMMARY

import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';

export const LANGUAGE_STORAGE_KEY = 'atlas-language';

export type Language = 'en' | 'ru';

type TranslationKey =
  | 'nutrition.actions'
  | 'nutrition.active'
  | 'nutrition.addFood'
  | 'nutrition.addPlannedEntry'
  | 'nutrition.applyResult'
  | 'nutrition.applyToWeek'
  | 'nutrition.archiveProduct'
  | 'nutrition.archived'
  | 'nutrition.cancelEdit'
  | 'nutrition.cancelEntryEdit'
  | 'nutrition.carbs'
  | 'nutrition.carbsPer100g'
  | 'nutrition.calories'
  | 'nutrition.caloriesPer100g'
  | 'nutrition.chooseProduct'
  | 'nutrition.createProduct'
  | 'nutrition.createdCount'
  | 'nutrition.dailyLogDescription'
  | 'nutrition.deleteEntry'
  | 'nutrition.editProduct'
  | 'nutrition.emptyDayDescription'
  | 'nutrition.emptyDayTitle'
  | 'nutrition.emptyProductLibrary'
  | 'nutrition.emptyProductLibraryDescription'
  | 'nutrition.emptyWeeklyPlanDescription'
  | 'nutrition.emptyWeeklyPlanTitle'
  | 'nutrition.entries'
  | 'nutrition.entryCount'
  | 'nutrition.entryDeleted'
  | 'nutrition.entryUpdated'
  | 'nutrition.entryAdded'
  | 'nutrition.entryNotes'
  | 'nutrition.fat'
  | 'nutrition.fatPer100g'
  | 'nutrition.foodLog'
  | 'nutrition.forEntry'
  | 'nutrition.grams'
  | 'nutrition.gramsPositive'
  | 'nutrition.includeArchived'
  | 'nutrition.loadNutritionError'
  | 'nutrition.loadProductsError'
  | 'nutrition.loadWeeklyPlanError'
  | 'nutrition.loadingNutrition'
  | 'nutrition.loadingProducts'
  | 'nutrition.loadingWeeklyPlan'
  | 'nutrition.macroValidationError'
  | 'nutrition.manageProducts'
  | 'nutrition.mealLabel'
  | 'nutrition.nameValidationError'
  | 'nutrition.nextDay'
  | 'nutrition.nextWeek'
  | 'nutrition.notes'
  | 'nutrition.previousDay'
  | 'nutrition.previousWeek'
  | 'nutrition.product'
  | 'nutrition.plannedEntry'
  | 'nutrition.productArchived'
  | 'nutrition.productCreated'
  | 'nutrition.productLibrary'
  | 'nutrition.productLibraryDescription'
  | 'nutrition.productName'
  | 'nutrition.products'
  | 'nutrition.productRestored'
  | 'nutrition.productUpdated'
  | 'nutrition.protein'
  | 'nutrition.proteinPer100g'
  | 'nutrition.restoreProduct'
  | 'nutrition.retry'
  | 'nutrition.saveBeforeApply'
  | 'nutrition.saveProduct'
  | 'nutrition.saveEntry'
  | 'nutrition.saveTemplate'
  | 'nutrition.searchProducts'
  | 'nutrition.selectProduct'
  | 'nutrition.skippedCount'
  | 'nutrition.status'
  | 'nutrition.templateApplied'
  | 'nutrition.templateNotes'
  | 'nutrition.templateSaved'
  | 'nutrition.templateTitle'
  | 'nutrition.title'
  | 'nutrition.today'
  | 'nutrition.updateProduct'
  | 'nutrition.unknownProduct'
  | 'nutrition.weeklyPlan'
  | 'nutrition.weeklyPlanDescription'
  | 'nutrition.weeklyTotals'
  | 'nutrition.weekOf';

const translations: Record<Language, Record<TranslationKey, string>> = {
  en: {
    'nutrition.actions': 'Actions',
    'nutrition.active': 'Active',
    'nutrition.addFood': 'Add food',
    'nutrition.addPlannedEntry': 'Add planned entry',
    'nutrition.applyResult': 'Apply result',
    'nutrition.applyToWeek': 'Apply to Week',
    'nutrition.archiveProduct': 'Archive',
    'nutrition.archived': 'Archived',
    'nutrition.cancelEdit': 'Cancel edit',
    'nutrition.cancelEntryEdit': 'Cancel',
    'nutrition.carbs': 'Carbs',
    'nutrition.carbsPer100g': 'Carbs per 100g',
    'nutrition.calories': 'Calories',
    'nutrition.caloriesPer100g': 'Calories per 100g',
    'nutrition.chooseProduct': 'Choose a product',
    'nutrition.createProduct': 'Create product',
    'nutrition.createdCount': 'created',
    'nutrition.dailyLogDescription':
      'Log products and grams eaten on the selected day. Totals are calculated from product snapshots.',
    'nutrition.deleteEntry': 'Delete',
    'nutrition.editProduct': 'Edit',
    'nutrition.emptyDayDescription':
      'Add a product and grams to start calculating daily nutrition totals.',
    'nutrition.emptyDayTitle': 'No food entries yet',
    'nutrition.emptyProductLibrary': 'No products yet',
    'nutrition.emptyProductLibraryDescription': 'Create your first product to use it in food logs.',
    'nutrition.emptyWeeklyPlanDescription':
      'Save a title or add planned entries to create the weekly template for this week.',
    'nutrition.emptyWeeklyPlanTitle': 'No weekly plan yet',
    'nutrition.entries': 'Food entries',
    'nutrition.entryCount': 'entries',
    'nutrition.entryDeleted': 'Food entry deleted',
    'nutrition.entryUpdated': 'Food entry updated',
    'nutrition.entryAdded': 'Food entry added',
    'nutrition.entryNotes': 'Entry notes',
    'nutrition.fat': 'Fat',
    'nutrition.fatPer100g': 'Fat per 100g',
    'nutrition.foodLog': 'Food log',
    'nutrition.forEntry': 'for entry',
    'nutrition.grams': 'Grams',
    'nutrition.gramsPositive': 'Grams must be greater than 0',
    'nutrition.includeArchived': 'Include archived',
    'nutrition.loadNutritionError': 'Failed to load nutrition data',
    'nutrition.loadProductsError': 'Failed to load products',
    'nutrition.loadWeeklyPlanError': 'Failed to load weekly plan',
    'nutrition.loadingNutrition': 'Loading nutrition data',
    'nutrition.loadingProducts': 'Loading products',
    'nutrition.loadingWeeklyPlan': 'Loading weekly plan',
    'nutrition.macroValidationError': 'Macro values must be zero or greater',
    'nutrition.manageProducts': 'Manage products',
    'nutrition.mealLabel': 'Meal label',
    'nutrition.nameValidationError': 'Product name is required',
    'nutrition.nextDay': 'Next day',
    'nutrition.nextWeek': 'Next week',
    'nutrition.notes': 'Notes',
    'nutrition.previousDay': 'Previous day',
    'nutrition.previousWeek': 'Previous week',
    'nutrition.plannedEntry': 'Planned entry',
    'nutrition.product': 'Product',
    'nutrition.productArchived': 'Product archived',
    'nutrition.productCreated': 'Product created',
    'nutrition.productLibrary': 'Product Library',
    'nutrition.productLibraryDescription':
      'Create private foods with macros per 100g for daily logs and weekly plans.',
    'nutrition.productName': 'Product name',
    'nutrition.products': 'Products',
    'nutrition.productRestored': 'Product restored',
    'nutrition.productUpdated': 'Product updated',
    'nutrition.protein': 'Protein',
    'nutrition.proteinPer100g': 'Protein per 100g',
    'nutrition.restoreProduct': 'Restore',
    'nutrition.retry': 'Retry',
    'nutrition.saveBeforeApply': 'Save the weekly template before applying it to days',
    'nutrition.saveEntry': 'Save',
    'nutrition.saveProduct': 'Save product',
    'nutrition.saveTemplate': 'Save Template',
    'nutrition.searchProducts': 'Search products',
    'nutrition.selectProduct': 'Select product',
    'nutrition.skippedCount': 'skipped',
    'nutrition.status': 'Status',
    'nutrition.templateApplied': 'Template applied to week',
    'nutrition.templateNotes': 'Template notes',
    'nutrition.templateSaved': 'Template saved',
    'nutrition.templateTitle': 'Template title',
    'nutrition.title': 'Nutrition',
    'nutrition.today': 'Today',
    'nutrition.updateProduct': 'Update product',
    'nutrition.unknownProduct': 'Unknown product',
    'nutrition.weeklyPlan': 'Weekly Plan',
    'nutrition.weeklyPlanDescription':
      'Plan reusable meals for a week. Saving edits the template only; applying seeds empty factual days.',
    'nutrition.weeklyTotals': 'Weekly totals',
    'nutrition.weekOf': 'Week of',
  },
  ru: {
    'nutrition.actions': 'Действия',
    'nutrition.active': 'Активные',
    'nutrition.addFood': 'Добавить продукт',
    'nutrition.addPlannedEntry': 'Добавить плановую запись',
    'nutrition.applyResult': 'Результат применения',
    'nutrition.applyToWeek': 'Применить к неделе',
    'nutrition.archiveProduct': 'Архивировать',
    'nutrition.archived': 'Архивные',
    'nutrition.cancelEdit': 'Отменить редактирование',
    'nutrition.cancelEntryEdit': 'Отмена',
    'nutrition.carbs': 'Углеводы',
    'nutrition.carbsPer100g': 'Углеводы на 100 г',
    'nutrition.calories': 'Калории',
    'nutrition.caloriesPer100g': 'Калории на 100 г',
    'nutrition.chooseProduct': 'Выберите продукт',
    'nutrition.createProduct': 'Создать продукт',
    'nutrition.createdCount': 'создано',
    'nutrition.dailyLogDescription':
      'Записывайте продукты и граммы за выбранный день. Итоги считаются по снимкам КБЖУ продуктов.',
    'nutrition.deleteEntry': 'Удалить',
    'nutrition.editProduct': 'Редактировать',
    'nutrition.emptyDayDescription':
      'Добавьте продукт и граммы, чтобы начать считать дневные итоги питания.',
    'nutrition.emptyDayTitle': 'Записей питания пока нет',
    'nutrition.emptyProductLibrary': 'Продуктов пока нет',
    'nutrition.emptyProductLibraryDescription':
      'Создайте первый продукт, чтобы использовать его в дневнике питания.',
    'nutrition.emptyWeeklyPlanDescription':
      'Сохраните название или добавьте плановые записи, чтобы создать недельный шаблон.',
    'nutrition.emptyWeeklyPlanTitle': 'Недельного плана пока нет',
    'nutrition.entries': 'Записи питания',
    'nutrition.entryCount': 'записей',
    'nutrition.entryDeleted': 'Запись питания удалена',
    'nutrition.entryUpdated': 'Запись питания обновлена',
    'nutrition.entryAdded': 'Запись питания добавлена',
    'nutrition.entryNotes': 'Заметки к записи',
    'nutrition.fat': 'Жиры',
    'nutrition.fatPer100g': 'Жиры на 100 г',
    'nutrition.foodLog': 'Дневник питания',
    'nutrition.forEntry': 'для записи',
    'nutrition.grams': 'Граммы',
    'nutrition.gramsPositive': 'Граммы должны быть больше 0',
    'nutrition.includeArchived': 'Включая архивные',
    'nutrition.loadNutritionError': 'Не удалось загрузить данные питания',
    'nutrition.loadProductsError': 'Не удалось загрузить продукты',
    'nutrition.loadWeeklyPlanError': 'Не удалось загрузить недельный план',
    'nutrition.loadingNutrition': 'Загрузка данных питания',
    'nutrition.loadingProducts': 'Загрузка продуктов',
    'nutrition.loadingWeeklyPlan': 'Загрузка недельного плана',
    'nutrition.macroValidationError': 'Значения КБЖУ должны быть не меньше нуля',
    'nutrition.manageProducts': 'Управлять продуктами',
    'nutrition.mealLabel': 'Прием пищи',
    'nutrition.nameValidationError': 'Название продукта обязательно',
    'nutrition.nextDay': 'Следующий день',
    'nutrition.nextWeek': 'Следующая неделя',
    'nutrition.notes': 'Заметки',
    'nutrition.previousDay': 'Предыдущий день',
    'nutrition.previousWeek': 'Предыдущая неделя',
    'nutrition.plannedEntry': 'Плановая запись',
    'nutrition.product': 'Продукт',
    'nutrition.productArchived': 'Продукт архивирован',
    'nutrition.productCreated': 'Продукт создан',
    'nutrition.productLibrary': 'Библиотека продуктов',
    'nutrition.productLibraryDescription':
      'Создавайте личные продукты с КБЖУ на 100 г для дневника и недельных планов.',
    'nutrition.productName': 'Название продукта',
    'nutrition.products': 'Продукты',
    'nutrition.productRestored': 'Продукт восстановлен',
    'nutrition.productUpdated': 'Продукт обновлен',
    'nutrition.protein': 'Белки',
    'nutrition.proteinPer100g': 'Белки на 100 г',
    'nutrition.restoreProduct': 'Восстановить',
    'nutrition.retry': 'Повторить',
    'nutrition.saveBeforeApply': 'Сначала сохраните недельный шаблон',
    'nutrition.saveEntry': 'Сохранить',
    'nutrition.saveProduct': 'Сохранить продукт',
    'nutrition.saveTemplate': 'Сохранить шаблон',
    'nutrition.searchProducts': 'Поиск продуктов',
    'nutrition.selectProduct': 'Выберите продукт',
    'nutrition.skippedCount': 'пропущено',
    'nutrition.status': 'Статус',
    'nutrition.templateApplied': 'Шаблон применен к неделе',
    'nutrition.templateNotes': 'Заметки шаблона',
    'nutrition.templateSaved': 'Шаблон сохранен',
    'nutrition.templateTitle': 'Название шаблона',
    'nutrition.title': 'Питание',
    'nutrition.today': 'Сегодня',
    'nutrition.updateProduct': 'Обновить продукт',
    'nutrition.unknownProduct': 'Неизвестный продукт',
    'nutrition.weeklyPlan': 'Недельный план',
    'nutrition.weeklyPlanDescription':
      'Планируйте повторяемые приемы пищи на неделю. Сохранение меняет только шаблон; применение заполняет пустые фактические дни.',
    'nutrition.weeklyTotals': 'Итоги недели',
    'nutrition.weekOf': 'Неделя от',
  },
};

type I18nContextValue = {
  language: Language;
  setLanguage: (language: Language) => void;
  t: (key: TranslationKey) => string;
};

const I18nContext = createContext<I18nContextValue | null>(null);

function readStoredLanguage(): Language {
  if (typeof window === 'undefined') {
    return 'en';
  }

  try {
    const storedLanguage = window.localStorage?.getItem(LANGUAGE_STORAGE_KEY);
    return storedLanguage === 'ru' ? 'ru' : 'en';
  } catch {
    return 'en';
  }
}

function persistLanguage(language: Language) {
  try {
    window.localStorage?.setItem(LANGUAGE_STORAGE_KEY, language);
  } catch {
    // Storage may be unavailable in private or restricted browser contexts.
  }
}

// START_CONTRACT: I18nProvider
//   PURPOSE: Provide selected Atlas language state and translation lookup.
//   INPUTS: { children: ReactNode - subtree using translations, initialLanguage?: Language - test/default override }
//   OUTPUTS: { JSX.Element - i18n context provider }
//   SIDE_EFFECTS: Persists selected language to localStorage and mirrors it to document.documentElement.lang.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: I18nProvider
export function I18nProvider({
  children,
  initialLanguage,
}: {
  children: ReactNode;
  initialLanguage?: Language;
}) {
  const [language, setLanguage] = useState<Language>(() => initialLanguage ?? readStoredLanguage());

  useEffect(() => {
    document.documentElement.lang = language;
    persistLanguage(language);
  }, [language]);

  const value = useMemo<I18nContextValue>(
    () => ({
      language,
      setLanguage,
      t: (key) => translations[language][key],
    }),
    [language],
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

// START_CONTRACT: useI18n
//   PURPOSE: Expose current Atlas translation helpers to page components.
//   INPUTS: none.
//   OUTPUTS: { I18nContextValue - current language, setter, and lookup function }
//   SIDE_EFFECTS: Throws when used outside I18nProvider.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: useI18n
export function useI18n() {
  const context = useContext(I18nContext);

  if (!context) {
    throw new Error('useI18n must be used inside I18nProvider');
  }

  return context;
}
