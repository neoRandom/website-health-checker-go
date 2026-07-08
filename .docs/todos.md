## Business Rules

- [ ] Update models
    - [ ] Site
        - [ ] Split "Url" into "Host" and "Endpoint"
        - [ ] Add "Description"

- [ ] Add setting information: Check Method, Pooling/Interval, Timeout

## Code Quality

- [ ] Improve error handling.
    - [ ] Add error messages and runtime checks

- [ ] Add context (ctx) to every part that I'll benefit from having it
    - [ ] Database operations

## Issues

- [ ] Fix the SQLite concurrency issue: Error worker saving site 'https://github.com': database is locked (5) (SQLITE_BUSY)

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

#### Challenges:

- [ ] Correct usage of concurrency for probing

- [ ] Avoid data leaking and hanging calls

- [ ] Minimal footprint (memory and network)

- [ ] Resilience against edge cases
    - [ ] Network going down
    - [ ] Anti-bot systems

- [ ] Apply quality principles
    - [ ] Dependency Inversion
    - [ ] Single Responsibility
    - [ ] Separate business logic from external dependencies

