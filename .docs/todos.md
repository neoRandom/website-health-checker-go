- Improve error handling.
Add error messages and runtime checks

- Fix the SQLite concurrency issue.
Error worker saving site 'https://github.com': database is locked (5) (SQLITE_BUSY)

- Add unit tests, integration and stress tests

- Test the before and after of the project
    - And improve performance if needed

- Implement relisience features before public deploy
    - Rate Limiter
    - Request Size Limit (downscaled; probably already handled by Fly.io)
    - Authentication for CRUD operations
    - Request Caching (downscaled; probably already handled by Fly.io)
    - Query and Request Timeouts
    - Idempotency
    - Query size limits / Pagination
    - Input validation
    - External persistent structured logging (to know what happened)
    - Panic recovery (to avoid downtime)
    - Retries (see SQLite concurrency issue)
    - Security Headers
    - Warnings/Alerts on high usage (e.g.: egress)
    - Database Indexes (e.g.: results)

- Add context (ctx) to every part that I'll benefit from having it
    - Database operations
