// FILE: apps/web-admin/src/app/i18n.tsx
// VERSION: 1.0.0
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

import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';

export const LANGUAGE_STORAGE_KEY = 'atlas-language';

export type Language = 'en' | 'ru';

type TranslationKey =
  | 'nutrition.actions'
  | 'nutrition.active'
  | 'nutrition.archiveProduct'
  | 'nutrition.archived'
  | 'nutrition.cancelEdit'
  | 'nutrition.carbsPer100g'
  | 'nutrition.caloriesPer100g'
  | 'nutrition.createProduct'
  | 'nutrition.editProduct'
  | 'nutrition.emptyProductLibrary'
  | 'nutrition.emptyProductLibraryDescription'
  | 'nutrition.fatPer100g'
  | 'nutrition.foodLog'
  | 'nutrition.includeArchived'
  | 'nutrition.loadProductsError'
  | 'nutrition.loadingProducts'
  | 'nutrition.macroValidationError'
  | 'nutrition.nameValidationError'
  | 'nutrition.notes'
  | 'nutrition.productArchived'
  | 'nutrition.productCreated'
  | 'nutrition.productLibrary'
  | 'nutrition.productLibraryDescription'
  | 'nutrition.productName'
  | 'nutrition.products'
  | 'nutrition.productRestored'
  | 'nutrition.productUpdated'
  | 'nutrition.proteinPer100g'
  | 'nutrition.restoreProduct'
  | 'nutrition.retry'
  | 'nutrition.saveProduct'
  | 'nutrition.searchProducts'
  | 'nutrition.status'
  | 'nutrition.updateProduct';

const translations: Record<Language, Record<TranslationKey, string>> = {
  en: {
    'nutrition.actions': 'Actions',
    'nutrition.active': 'Active',
    'nutrition.archiveProduct': 'Archive',
    'nutrition.archived': 'Archived',
    'nutrition.cancelEdit': 'Cancel edit',
    'nutrition.carbsPer100g': 'Carbs per 100g',
    'nutrition.caloriesPer100g': 'Calories per 100g',
    'nutrition.createProduct': 'Create product',
    'nutrition.editProduct': 'Edit',
    'nutrition.emptyProductLibrary': 'No products yet',
    'nutrition.emptyProductLibraryDescription': 'Create your first product to use it in food logs.',
    'nutrition.fatPer100g': 'Fat per 100g',
    'nutrition.foodLog': 'Food log',
    'nutrition.includeArchived': 'Include archived',
    'nutrition.loadProductsError': 'Failed to load products',
    'nutrition.loadingProducts': 'Loading products',
    'nutrition.macroValidationError': 'Macro values must be zero or greater',
    'nutrition.nameValidationError': 'Product name is required',
    'nutrition.notes': 'Notes',
    'nutrition.productArchived': 'Product archived',
    'nutrition.productCreated': 'Product created',
    'nutrition.productLibrary': 'Product Library',
    'nutrition.productLibraryDescription':
      'Create private foods with macros per 100g for daily logs and weekly plans.',
    'nutrition.productName': 'Product name',
    'nutrition.products': 'Products',
    'nutrition.productRestored': 'Product restored',
    'nutrition.productUpdated': 'Product updated',
    'nutrition.proteinPer100g': 'Protein per 100g',
    'nutrition.restoreProduct': 'Restore',
    'nutrition.retry': 'Retry',
    'nutrition.saveProduct': 'Save product',
    'nutrition.searchProducts': 'Search products',
    'nutrition.status': 'Status',
    'nutrition.updateProduct': 'Update product',
  },
  ru: {
    'nutrition.actions': 'Действия',
    'nutrition.active': 'Активные',
    'nutrition.archiveProduct': 'Архивировать',
    'nutrition.archived': 'Архивные',
    'nutrition.cancelEdit': 'Отменить редактирование',
    'nutrition.carbsPer100g': 'Углеводы на 100 г',
    'nutrition.caloriesPer100g': 'Калории на 100 г',
    'nutrition.createProduct': 'Создать продукт',
    'nutrition.editProduct': 'Редактировать',
    'nutrition.emptyProductLibrary': 'Продуктов пока нет',
    'nutrition.emptyProductLibraryDescription':
      'Создайте первый продукт, чтобы использовать его в дневнике питания.',
    'nutrition.fatPer100g': 'Жиры на 100 г',
    'nutrition.foodLog': 'Дневник питания',
    'nutrition.includeArchived': 'Включая архивные',
    'nutrition.loadProductsError': 'Не удалось загрузить продукты',
    'nutrition.loadingProducts': 'Загрузка продуктов',
    'nutrition.macroValidationError': 'Значения КБЖУ должны быть не меньше нуля',
    'nutrition.nameValidationError': 'Название продукта обязательно',
    'nutrition.notes': 'Заметки',
    'nutrition.productArchived': 'Продукт архивирован',
    'nutrition.productCreated': 'Продукт создан',
    'nutrition.productLibrary': 'Библиотека продуктов',
    'nutrition.productLibraryDescription':
      'Создавайте личные продукты с КБЖУ на 100 г для дневника и недельных планов.',
    'nutrition.productName': 'Название продукта',
    'nutrition.products': 'Продукты',
    'nutrition.productRestored': 'Продукт восстановлен',
    'nutrition.productUpdated': 'Продукт обновлен',
    'nutrition.proteinPer100g': 'Белки на 100 г',
    'nutrition.restoreProduct': 'Восстановить',
    'nutrition.retry': 'Повторить',
    'nutrition.saveProduct': 'Сохранить продукт',
    'nutrition.searchProducts': 'Поиск продуктов',
    'nutrition.status': 'Статус',
    'nutrition.updateProduct': 'Обновить продукт',
  },
};

type I18nContextValue = {
  language: Language;
  setLanguage: (language: Language) => void;
  t: (key: TranslationKey) => string;
};

const I18nContext = createContext<I18nContextValue | null>(null);

function readStoredLanguage(): Language {
  if (typeof window === 'undefined' || !window.localStorage) {
    return 'en';
  }

  const storedLanguage = window.localStorage.getItem(LANGUAGE_STORAGE_KEY);
  return storedLanguage === 'ru' ? 'ru' : 'en';
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
    window.localStorage?.setItem(LANGUAGE_STORAGE_KEY, language);
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
