## Business Rules

- [x] Update models
    - [x] Site
        - [x] Add "Description"

- [x] Add config fields: Pooling/Interval, Timeout

## Code Quality

- [x] Improve response handling.
    - [x] Add error messages and enhance runtime checks
    - [x] Set Headers in every endpoint

- [x] Enhance background context support to improve graceful shutdowns
    - [x] Database operations
    - [x] Services
    - [x] Schedules

## Issues

- [x] Fix the SQLite concurrency issue: database is locked (5) (SQLITE_BUSY)

## Testing

- [ ] Add unit tests, integration and stress tests

- [ ] Improve Bruno's endpoint testing. Broader automatic testing.

- [ ] Test the before and after of the project, and improve performance if needed

- [ ] Implement automatic testing on push using GitHub Actions

## Resilience

- [ ] Implement deploy resilience features before public deploy
    - [ ] Rate Limiter
    - [ ] Request Size Limit (downscaled; probably already handled by Fly.io)
    - [ ] Authentication for CRUD operations
    - [ ] Request Caching (downscaled; probably already handled by Fly.io)
    - [ ] Query and Request Timeouts
    - [ ] Idempotency
    - [ ] Query size limits / Pagination
    - [ ] Input validation
    - [ ] External persistent structured logging (to know what happened)
    - [ ] Panic recovery (to avoid downtime)
    - [ ] Retries (see SQLite concurrency issue)
    - [ ] Security Headers
    - [ ] Warnings/Alerts on high usage (e.g.: egress)
    - [ ] Database Indexes (e.g.: results)

## Features / Roadmap

- [x] Health-check multiple websites

- [ ] Real-time dashboard

- [x] User-modifiable list of URLs

#### Technical details:

- [x] Add diagnosing methods
    - [x] Monitoring via Prometheus/Grafana
    - [x] Profiling via `net/http/pprof`

- [ ] Add ADRs and Code Documentation

