# maskthis.com - Global URL Shortener

Fast, privacy-respecting URL shortener running on Bunny Magic Containers.

## Setup Status

- ✅ Project structure created
- ✅ Go 1.21.6 installed
- ✅ Go module initialized
- ✅ Database migration ready
- ⏳ Bunny Database (in progress)
- ⏳ Backend implementation (next)

## Project Structure

```
maskthis/
├── cmd/api/              # Main application
├── internal/
│   ├── database/         # Database connection
│   ├── shortener/        # URL shortening logic
│   └── analytics/        # Click tracking
├── web/static/           # Frontend files
├── migrations/           # SQL migrations
├── .env.template         # Environment template
└── README.md            # This file
```

## Next Steps

1. Create Bunny Database at https://bunny.net
2. Copy `.env.template` to `.env`
3. Fill in database credentials
4. Run migration: Apply `migrations/001_initial.sql` to your database
5. Build the backend (we'll do this together)
6. Test locally
7. Deploy to Bunny Magic Containers

## Database Schema

Simple 2-table design:
- **urls** - Short codes mapped to long URLs
- **clicks** - Analytics for each redirect

## Tech Stack

- Backend: Go 1.21
- Database: Bunny Database (libSQL)
- Frontend: HTML + Tailwind CSS
- Deployment: Bunny Magic Containers

## Documentation

Full documentation in funlab-docs:
- Design: `docs/maskthis-url-shortener-design.md`
- Quick Start: `docs/maskthis-url-shortener-quickstart.md`
