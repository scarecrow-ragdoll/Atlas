# Scope

## In Scope

- Single-user self-hosted web application
- Optional PIN-based access control
- Dashboard with weekly summary
- Exercise library with user-created exercises
- Exercise media upload (images and video)
- Workout diary by date with calendar navigation
- Sets with weight, reps, optional RPE/RIR
- Exercise comments within workouts
- Working weight tracking and auto-population
- Cardio logging with type, duration, pulse/zone
- Body weight entries (standalone and within check-ins)
- Weekly body check-ins with measurements and photos
- Progress photos tied to check-ins
- Nutrition product catalog
- Weekly nutrition template with daily overrides
- Progress charts for training, body, nutrition
- AI prompt builder with persistent and one-time context
- AI export ZIP (manifest.json, data.json, summary.md, CSVs, photos)
- Week flags for AI context (sleep, stress, illness, etc.)
- AI review history with manual entry
- Full backup export (ZIP with manifest, data, media)
- Full backup import with dry-run validation
- Key scenario test coverage

## Out Of Scope

- User registration
- Multi-user mode
- Roles and permissions (beyond optional PIN)
- SaaS hosting
- Public user profiles
- Workout templates
- Training plan scheduling
- Quick repeat of previous workout
- Built-in exercise catalog (starter catalog)
- Recipes and meals
- Barcode scanner
- Apple Health integration
- Cloud backup
- Telegram bot
- OpenAI API integration (manual copy-paste to AI)
- Automatic in-app training plan generation
- Mobile companion app

## Non-Goals

- Competition with full-featured calorie tracking apps
- Real-time workout coaching
- Social features
- Public API for third-party integrations (MVP)
- Offline-first support (web only)
- Native mobile experience (web only)

## Dependencies

- Self-hosted deployment: Docker, Docker Compose
- Storage: PostgreSQL, Redis, file system volume for media
- Runtime: Bun 1.1+, Node 22+, Go 1.25
- Browser with modern JavaScript support

## Assumptions

- User has basic Docker deployment knowledge
- Single user per instance with multi-user-ready data model (all entities owned by default user via user_id, one default user created at bootstrap)
- AI analysis is performed externally (ChatGPT or similar)
- Media files are stored on local filesystem volume
- Backup/restore covers full data lifecycle (no incremental or partial backup)
- Cardio is a separate entity always attached to a DailyLog (daily aggregate entity per date)