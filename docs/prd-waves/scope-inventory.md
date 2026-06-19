# Scope Inventory

## Capability Groups

1. Foundation (infrastructure)
2. Exercise management
3. Workout tracking
4. Cardio/body tracking
5. Nutrition tracking
6. Progress visualization
7. AI export
8. AI review storage
9. Backup/restore

## User Journey Groups

- Exercise CRUD journey
- Workout logging journey
- Check-in journey
- Nutrition template journey
- AI export journey
- Backup journey

## Data Lifecycle Groups

- Settings (CRUD)
- UserProfile (CRUD)
- Exercises (CRUD + media)
- Workouts (create by date, update)
- Cardio (CRUD)
- Body (CRUD)
- Nutrition (template + override)
- AI (export + review)
- Backup (full cycle)

## Integration And Operations Groups

- GraphQL API
- PostgreSQL persistence
- Redis sessions
- File storage for media
- ZIP generation

## Client Experience Groups

- 11 pages defined in frontend-pages/
- Shared loading/error/empty states
- Date-oriented navigation
- PIN auth flow

## Security Compliance Groups

- PIN hashing (not plaintext)
- Session-based access
- Media file protection
- No sensitive logging

## Explicit Deferrals

- Apple Health
- Telegram bot
- OpenAI API
- Barcode scanner
- Multi-user mode