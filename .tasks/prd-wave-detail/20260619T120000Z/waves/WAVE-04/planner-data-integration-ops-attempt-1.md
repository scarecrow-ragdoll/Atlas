# WAVE-04 Planner: Data / API / Integration / Ops

## Data Lifecycle

### Tables and Relationships

1. **cardio_entry**
   - id (UUID, PK), user_id (UUID, FK → users), daily_log_id (UUID, FK → daily_log), cardio_type (VARCHAR, NOT NULL), duration_minutes (INT, NOT NULL), avg_pulse (INT, nullable), heart_rate_zone (VARCHAR, nullable), notes (TEXT, nullable), created_at, updated_at
   - FK: daily_log_id → daily_log(id) ON DELETE CASCADE
   - Indexes: idx_cardio_entry_daily_log (daily_log_id), idx_cardio_entry_user_date (user_id via daily_log)

2. **body_weight_entry**
   - id (UUID, PK), user_id (UUID, FK → users), date (DATE, NOT NULL), weight (REAL, NOT NULL), source (VARCHAR, NOT NULL), notes (TEXT, nullable), created_at, updated_at
   - Unique: (user_id, date) — one weight per user per date? **Recommend:** allow multiple entries per date (source varies). Unique constraint on (user_id, date, source) or no unique.
   - Indexes: idx_body_weight_user_date (user_id, date DESC)

3. **body_check_in**
   - id (UUID, PK), user_id (UUID, FK → users), date (DATE, NOT NULL UNIQUE), weight (REAL, nullable), body_fat_percentage (REAL, nullable), notes (TEXT, nullable), created_at, updated_at
   - Unique: date (one check-in per date)
   - Indexes: idx_body_check_in_date (date DESC)

4. **body_measurement**
   - id (UUID, PK), check_in_id (UUID, FK → body_check_in), measurement_type (VARCHAR, NOT NULL), side (VARCHAR, nullable), value (REAL, NOT NULL), created_at, updated_at
   - FK: check_in_id → body_check_in(id) ON DELETE CASCADE
   - Unique: (check_in_id, measurement_type, side) — one measurement per type+side per check-in
   - Indexes: idx_body_measurement_checkin (check_in_id)

5. **progress_photo**
   - id (UUID, PK), check_in_id (UUID, FK → body_check_in), file_path (VARCHAR, NOT NULL), original_file_name (VARCHAR, NOT NULL), mime_type (VARCHAR, NOT NULL), size_bytes (BIGINT, NOT NULL), angle (VARCHAR, nullable), label (VARCHAR, nullable), notes (TEXT, nullable), created_at, updated_at
   - FK: check_in_id → body_check_in(id) ON DELETE CASCADE
   - Indexes: idx_progress_photo_checkin (check_in_id)

6. **week_flag**
   - id (UUID, PK), user_id (UUID, FK → users), week_start_date (DATE, NOT NULL), flag_type (VARCHAR, NOT NULL), notes (TEXT, nullable), created_at, updated_at
   - Unique: (week_start_date, flag_type) — one flag of each type per week
   - Indexes: idx_week_flag_week (week_start_date)

### Cascade Rules
- Delete body_check_in → cascade delete body_measurement records, progress_photo records, and physical photo files
- Delete cardio_entry → no cascade needed (owned by daily_log)
- Delete daily_log → cascade delete cardio_entries

## API Design

### GraphQL Operations

#### CardioEntry
```graphql
type CardioEntry {
  id: ID!
  dailyLogId: ID!
  cardioType: CardioType!
  durationMinutes: Int!
  avgPulse: Int
  heartRateZone: HeartRateZone
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

enum CardioType { WALKING RUNNING BIKE ELLIPTICAL TREADMILL OTHER }
enum HeartRateZone { ZONE_1 ZONE_2 ZONE_3 ZONE_4 ZONE_5 UNKNOWN }

input CreateCardioEntryInput {
  date: Date!
  cardioType: CardioType!
  durationMinutes: Int!
  avgPulse: Int
  heartRateZone: HeartRateZone
  notes: String
}

input UpdateCardioEntryInput {
  id: ID!
  cardioType: CardioType
  durationMinutes: Int
  avgPulse: Int
  heartRateZone: HeartRateZone
  notes: String
}

extend type Mutation {
  createCardioEntry(input: CreateCardioEntryInput!): CardioEntryResult!
  updateCardioEntry(input: UpdateCardioEntryInput!): CardioEntryResult!
  deleteCardioEntry(id: ID!): DeleteResult!
}

extend type Query {
  cardioEntries(date: Date!): [CardioEntry!]!
  cardioEntry(id: ID!): CardioEntry
}
```

#### BodyWeightEntry
```graphql
type BodyWeightEntry {
  id: ID!
  date: Date!
  weight: Float!
  source: WeightSource!
  notes: String
  createdAt: Time!
}

enum WeightSource { SCALE MANUAL UNKNOWN }

input CreateBodyWeightEntryInput {
  date: Date!
  weight: Float!
  source: WeightSource!
  notes: String
}

input UpdateBodyWeightEntryInput {
  id: ID!
  weight: Float
  source: WeightSource
  notes: String
}

extend type Mutation {
  createBodyWeightEntry(input: CreateBodyWeightEntryInput!): BodyWeightEntryResult!
  updateBodyWeightEntry(input: UpdateBodyWeightEntryInput!): BodyWeightEntryResult!
  deleteBodyWeightEntry(id: ID!): DeleteResult!
}

extend type Query {
  bodyWeightEntries(startDate: Date, endDate: Date): [BodyWeightEntry!]!
  bodyWeightEntry(id: ID!): BodyWeightEntry
  latestBodyWeight: BodyWeightEntry
}
```

#### BodyCheckIn (with nested types)
```graphql
type BodyCheckIn {
  id: ID!
  date: Date!
  weight: Float
  bodyFatPercentage: Float
  notes: String
  measurements: [BodyMeasurement!]!
  photos: [ProgressPhoto!]!
  createdAt: Time!
  updatedAt: Time!
}

type BodyMeasurement {
  id: ID!
  checkInId: ID!
  measurementType: MeasurementType!
  side: MeasurementSide
  value: Float!
  createdAt: Time!
}

type ProgressPhoto {
  id: ID!
  checkInId: ID!
  filePath: String!
  originalFileName: String!
  mimeType: String!
  sizeBytes: Int!
  angle: PhotoAngle
  label: String
  notes: String
  createdAt: Time!
  url: String!  # resolved download URL
}

enum MeasurementType { NECK SHOULDERS FOREARMS BICEPS CHEST WAIST ABDOMEN HIPS THIGH CALF }
enum MeasurementSide { LEFT RIGHT }
enum PhotoAngle { FRONT SIDE BACK CUSTOM }

input CreateBodyMeasurementInput {
  checkInId: ID!
  measurementType: MeasurementType!
  side: MeasurementSide
  value: Float!
}

input UpdateBodyMeasurementInput {
  id: ID!
  value: Float
  side: MeasurementSide
}

input CreateBodyCheckInInput {
  date: Date!
  weight: Float
  bodyFatPercentage: Float
  notes: String
}

input UpdateBodyCheckInInput {
  id: ID!
  weight: Float
  bodyFatPercentage: Float
  notes: String
}

extend type Mutation {
  createBodyCheckIn(input: CreateBodyCheckInInput!): BodyCheckInResult!
  updateBodyCheckIn(input: UpdateBodyCheckInInput!): BodyCheckInResult!
  deleteBodyCheckIn(id: ID!): DeleteResult!
  createBodyMeasurement(input: CreateBodyMeasurementInput!): BodyMeasurementResult!
  updateBodyMeasurement(input: UpdateBodyMeasurementInput!): BodyMeasurementResult!
  deleteBodyMeasurement(id: ID!): DeleteResult!
  deleteProgressPhoto(id: ID!): DeleteResult!
}

extend type Query {
  bodyCheckIns(startDate: Date, endDate: Date): [BodyCheckIn!]!
  bodyCheckIn(id: ID!): BodyCheckIn
  progressPhotos(checkInId: ID!): [ProgressPhoto!]!
}
```

#### WeekFlag
```graphql
type WeekFlag {
  id: ID!
  weekStartDate: Date!
  flagType: WeekFlagType!
  notes: String
  createdAt: Time!
}

enum WeekFlagType {
  POOR_SLEEP HIGH_STRESS ILLNESS INJURY_PAIN CYCLE
  CALORIE_DEFICIT SURPLUS MAINTENANCE MISSED_WORKOUTS TRAVEL
}

input CreateWeekFlagInput {
  weekStartDate: Date!
  flagType: WeekFlagType!
  notes: String
}

extend type Mutation {
  createWeekFlag(input: CreateWeekFlagInput!): WeekFlagResult!
  deleteWeekFlag(id: ID!): DeleteResult!
}

extend type Query {
  weekFlags(weekStartDate: Date): [WeekFlag!]!
}
```

#### Union Results (consistent with WAVE-02 pattern)
```graphql
union CardioEntryResult = CardioEntry | ValidationError | AuthError
union BodyWeightEntryResult = BodyWeightEntry | ValidationError | AuthError
union BodyCheckInResult = BodyCheckIn | ValidationError | AuthError
union BodyMeasurementResult = BodyMeasurement | ValidationError | AuthError
union WeekFlagResult = WeekFlag | ValidationError | AuthError
```

### REST Endpoints
- `POST /api/v1/progress-photos/upload` — multipart upload, fields: checkInId, file, angle, label, notes
- `GET /api/v1/progress-photos/{id}` — download file with correct content type
- `DELETE /api/v1/progress-photos/{id}` — delete photo and physical file

### Media Storage Pattern
- Path: `<WAVE-01 BasePath>/progress-photos/<checkin_id>/<uuid>.<ext>`
- File validation: server-side MIME detection (JPEG/PNG/WEBP), per-type size limits (25MB)
- Memory-safe upload: `r.ParseMultipartForm(maxBytes)` per TDEC-008

## Operations

### Log Markers
- `[CardioEntry][create|update|delete|get|list]`
- `[BodyWeightEntry][create|update|delete|get|list|latest]`
- `[BodyCheckIn][create|update|delete|get|list]`
- `[BodyMeasurement][create|update|delete]`
- `[ProgressPhoto][upload|download|delete]`
- `[WeekFlag][create|delete|list]`

### Error Codes (REST)
- `FILE_TOO_LARGE`, `INVALID_FILE_TYPE`, `NOT_FOUND`, `INTERNAL_ERROR`, `UNAUTHORIZED`

### Error Format
```json
{ "error": { "code": "ERROR_CODE", "message": "Human readable" } }
```

### Migration Strategy
- 6 goose migrations (00082-00087), sequential
- Down migrations available for rollback
- Migrations run at startup via existing goose mechanism

### DailyLog Auto-Creation
CardioEntry requires dailyLogId. For the `createCardioEntry(date: ...)` mutation:
1. Look up existing DailyLog for user + date
2. If not found, auto-create DailyLog with that date (see WAVE-03 for daily_log table definition; if WAVE-03 not yet deployed, WAVE-04 must create or ensure the daily_log table exists)
3. Create CardioEntry with the dailyLogId

### Rollout / Rollback
- Rollout: merge PR, CI builds and runs tests, Dokploy update. New tables created via goose migrations.
- Rollback: revert PR, deploy previous image, run goose down migrations.
- Compatibility: all new operations are additive. No existing API changes.