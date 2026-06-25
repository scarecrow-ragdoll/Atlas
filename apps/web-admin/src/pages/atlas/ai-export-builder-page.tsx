// FILE: apps/web-admin/src/pages/atlas/ai-export-builder-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the Atlas AI Export route as a real local prompt and ZIP builder.
//   SCOPE: Lets admins select date range and export sections, calls only local guarded AI export REST endpoints, and displays progress, errors, prompt preview, and safe download links; excludes external AI API calls, backend changes, new shell/sidebar/topbar, and reference/mock UI.
//   DEPENDS: @tanstack/react-query, react, apps/web-admin/src/app/i18n.tsx, apps/web-admin/src/pages/atlas/ai-export-api.ts, apps/web-admin/src/styles/atlas.css, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - API-backed `/atlas/ai-export` route content for local AI export generation.
// END_MODULE_MAP

import { useMutation } from '@tanstack/react-query';
import { type FormEvent, useMemo, useState } from 'react';
import { useI18n } from '../../app/i18n';
import '../../styles/atlas.css';
import {
  generateAtlasAiExport,
  type GenerateAtlasAiExportInput,
  type GenerateAtlasAiExportResult,
} from './ai-export-api';
import {
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
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
  Checkbox,
  Input,
  Label,
  Textarea,
} from '@shared/ui';

type AiExportBuilderPageProps = {
  initialStartDate?: string;
  initialEndDate?: string;
};

type ExportSectionId = 'nutrition' | 'cardio' | 'measurements' | 'photos';

type ExportSection = {
  id: ExportSectionId;
  label: string;
  description: string;
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
};

function toDateString(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function addDays(date: Date, days: number) {
  const nextDate = new Date(date);
  nextDate.setDate(nextDate.getDate() + days);
  return nextDate;
}

function defaultStartDate() {
  return toDateString(addDays(new Date(), -27));
}

function defaultEndDate() {
  return toDateString(new Date());
}

function errorMessageFromUnknown(error: unknown): string {
  return error instanceof Error ? error.message : 'Request failed';
}

// START_CONTRACT: buildGenerateInput
//   PURPOSE: Convert controlled AI export form state into the backend REST request contract.
//   INPUTS: { dates, toggles, userComment - current page form state }
//   OUTPUTS: { GenerateAtlasAiExportInput - date range, section booleans, and nullable comment }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
// END_CONTRACT: buildGenerateInput
function buildGenerateInput({
  dateRangeEnd,
  dateRangeStart,
  includeCardio,
  includeMeasurements,
  includeNutrition,
  includePhotos,
  userComment,
}: {
  dateRangeStart: string;
  dateRangeEnd: string;
  includePhotos: boolean;
  includeNutrition: boolean;
  includeCardio: boolean;
  includeMeasurements: boolean;
  userComment: string;
}): GenerateAtlasAiExportInput {
  return {
    dateRangeStart,
    dateRangeEnd,
    includePhotos,
    includeNutrition,
    includeCardio,
    includeMeasurements,
    userComment: userComment.trim() || null,
  };
}

function SectionCheckbox({ checked, description, id, label, onCheckedChange }: ExportSection) {
  return (
    <div className="atlas-checkbox-row">
      <Checkbox
        aria-label={label}
        checked={checked}
        id={`ai-export-${id}`}
        onCheckedChange={(nextChecked) => onCheckedChange(nextChecked === true)}
      />
      <div className="grid gap-1">
        <Label htmlFor={`ai-export-${id}`}>{label}</Label>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}

function ReadyState({ result }: { result: GenerateAtlasAiExportResult }) {
  const { t } = useI18n();
  const exportResult = result.export;

  return (
    <AdminSection
      title={t('aiExport.readyTitle')}
      description={t('aiExport.localDownloadDescription')}
    >
      <div className="grid gap-4">
        <div className="flex flex-wrap gap-2">
          <Badge>{exportResult.dateRangeStart}</Badge>
          <Badge>{exportResult.dateRangeEnd}</Badge>
          {exportResult.includeNutrition ? <Badge>{t('aiExport.nutrition')}</Badge> : null}
          {exportResult.includeCardio ? <Badge>{t('aiExport.cardio')}</Badge> : null}
          {exportResult.includeMeasurements ? <Badge>{t('aiExport.measurements')}</Badge> : null}
          {exportResult.includePhotos ? <Badge>{t('aiExport.photos')}</Badge> : null}
        </div>
        <Button asChild>
          <a href={exportResult.downloadUrl}>{t('aiExport.downloadZip')}</a>
        </Button>
        <Card>
          <CardHeader>
            <CardTitle>
              <h3>{t('aiExport.promptPreview')}</h3>
            </CardTitle>
            <CardDescription>{t('aiExport.promptPreviewDescription')}</CardDescription>
          </CardHeader>
          <CardContent>
            <pre className="atlas-prompt-preview">
              {exportResult.generatedPrompt || t('aiExport.noPrompt')}
            </pre>
          </CardContent>
        </Card>
      </div>
    </AdminSection>
  );
}

// START_CONTRACT: AiExportBuilderPage
//   PURPOSE: Render local Atlas AI export generation controls and result states.
//   INPUTS: { initialStartDate?: string, initialEndDate?: string - deterministic test/default date overrides }
//   OUTPUTS: { JSX.Element - AI export builder with date controls, section toggles, privacy warning, progress, error, and ready states }
//   SIDE_EFFECTS: Calls generateAtlasAiExport through React Query mutation when the form is submitted.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
// END_CONTRACT: AiExportBuilderPage
export default function AiExportBuilderPage({
  initialEndDate,
  initialStartDate,
}: AiExportBuilderPageProps = {}) {
  const { t } = useI18n();
  const [dateRangeStart, setDateRangeStart] = useState(initialStartDate ?? defaultStartDate());
  const [dateRangeEnd, setDateRangeEnd] = useState(initialEndDate ?? defaultEndDate());
  const [includeNutrition, setIncludeNutrition] = useState(true);
  const [includeCardio, setIncludeCardio] = useState(true);
  const [includeMeasurements, setIncludeMeasurements] = useState(true);
  const [includePhotos, setIncludePhotos] = useState(false);
  const [userComment, setUserComment] = useState('');

  const generateMutation = useMutation({
    mutationFn: (input: GenerateAtlasAiExportInput) => generateAtlasAiExport(input),
  });

  const sectionOptions = useMemo<ExportSection[]>(
    () => [
      {
        id: 'nutrition',
        label: t('aiExport.nutrition'),
        description: t('aiExport.nutritionDescription'),
        checked: includeNutrition,
        onCheckedChange: setIncludeNutrition,
      },
      {
        id: 'cardio',
        label: t('aiExport.cardio'),
        description: t('aiExport.cardioDescription'),
        checked: includeCardio,
        onCheckedChange: setIncludeCardio,
      },
      {
        id: 'measurements',
        label: t('aiExport.measurements'),
        description: t('aiExport.measurementsDescription'),
        checked: includeMeasurements,
        onCheckedChange: setIncludeMeasurements,
      },
      {
        id: 'photos',
        label: t('aiExport.photos'),
        description: t('aiExport.photosDescription'),
        checked: includePhotos,
        onCheckedChange: setIncludePhotos,
      },
    ],
    [includeCardio, includeMeasurements, includeNutrition, includePhotos, t],
  );

  function currentInput() {
    return buildGenerateInput({
      dateRangeStart,
      dateRangeEnd,
      includePhotos,
      includeNutrition,
      includeCardio,
      includeMeasurements,
      userComment,
    });
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    generateMutation.mutate(currentInput());
  }

  function retryGenerate() {
    generateMutation.mutate(currentInput());
  }

  const errorMessage = generateMutation.isError
    ? errorMessageFromUnknown(generateMutation.error)
    : null;

  return (
    <AdminPageShell className="atlas-ai-export">
      <AdminPageHeader title={t('aiExport.title')} description={t('aiExport.description')} />

      <Alert>
        <AlertTitle>{t('aiExport.privacyTitle')}</AlertTitle>
        <AlertDescription>
          {t('aiExport.privacyDescription')} {t('aiExport.photosExcluded')}
        </AlertDescription>
      </Alert>

      <form className="grid gap-6" onSubmit={handleSubmit}>
        <AdminSection
          title={t('aiExport.dateRange')}
          description={t('aiExport.dateRangeDescription')}
        >
          <div className="atlas-form-grid">
            <div className="grid gap-2">
              <Label htmlFor="ai-export-start-date">{t('aiExport.startDate')}</Label>
              <Input
                id="ai-export-start-date"
                onChange={(event) => setDateRangeStart(event.target.value)}
                required
                type="date"
                value={dateRangeStart}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="ai-export-end-date">{t('aiExport.endDate')}</Label>
              <Input
                id="ai-export-end-date"
                onChange={(event) => setDateRangeEnd(event.target.value)}
                required
                type="date"
                value={dateRangeEnd}
              />
            </div>
          </div>
        </AdminSection>

        <AdminSection
          title={t('aiExport.includedSections')}
          description={t('aiExport.includedSectionsDescription')}
        >
          <div className="atlas-checkbox-grid">
            {sectionOptions.map((section) => (
              <SectionCheckbox key={section.id} {...section} />
            ))}
          </div>
        </AdminSection>

        <AdminSection
          title={t('aiExport.contextNotes')}
          description={t('aiExport.contextNotesDescription')}
        >
          <div className="grid gap-2">
            <Label htmlFor="ai-export-user-comment">{t('aiExport.userComment')}</Label>
            <Textarea
              id="ai-export-user-comment"
              onChange={(event) => setUserComment(event.target.value)}
              value={userComment}
            />
          </div>
        </AdminSection>

        <div className="flex flex-wrap items-center gap-3">
          <Button disabled={generateMutation.isPending} type="submit">
            {generateMutation.isPending ? t('aiExport.generating') : t('aiExport.generate')}
          </Button>
          <p className="text-sm text-muted-foreground">{t('aiExport.generateDescription')}</p>
        </div>
      </form>

      {generateMutation.isPending ? (
        <Alert>
          <AlertTitle>{t('aiExport.generating')}</AlertTitle>
          <AlertDescription>{t('aiExport.progressDescription')}</AlertDescription>
        </Alert>
      ) : null}

      {errorMessage ? (
        <Alert variant="destructive">
          <AlertTitle>{t('aiExport.errorTitle')}</AlertTitle>
          <AlertDescription>
            <span className="block">{errorMessage}</span>
            <Button className="mt-3" onClick={retryGenerate} type="button" variant="outline">
              {t('aiExport.retry')}
            </Button>
          </AlertDescription>
        </Alert>
      ) : null}

      {generateMutation.data && !generateMutation.isPending ? (
        <ReadyState result={generateMutation.data} />
      ) : null}
    </AdminPageShell>
  );
}
