# Bunny Magic Containers Deployment Guide

## Environment Variables for Bunny Container

Copy these into Bunny's container configuration:

```
PORT=8080
ENVIRONMENT=production
BASE_URL=https://maskthis.com

DATABASE_URL=libsql://01KG95ZJRZMR3YWYGJ746KMBRX-maskthis.lite.bunnydb.net/
DATABASE_AUTH_TOKEN=eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSJ9.eyJwIjp7InJvIjpudWxsLCJydyI6eyJucyI6WyJtYXNrdGhpcyJdLCJ0YWdzIjpudWxsfSwicm9hIjpudWxsLCJyd2EiOm51bGwsImRkbCI6bnVsbH0sImlhdCI6MTc2OTgzNDc5N30.6Qoxrjh6sXglv1__lLbRMMXu4SGGrmYYD9dCkCUEw148OAjNixE34ByixxEjivlL29QjqgQZgZpon0HPuvXmAA

SHORT_CODE_LENGTH=6
ANON_HOURLY_LIMIT=10
ANON_DAILY_LIMIT=100
```

## Container Settings

- **Name:** maskthis
- **Port:** 8080
- **Regions:** All (for global distribution)
- **Auto-scaling:** Enabled
- **Health Check:** /api/health

## Post-Deployment

After container is deployed, you'll get a URL like:
```
https://maskthis-xxxxx.b-cdn.net
```

Test it:
```bash
curl https://maskthis-xxxxx.b-cdn.net/api/health
```

Then configure DNS to point maskthis.com to this URL.
