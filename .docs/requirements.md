Deprecated

```
API:

- GET / (Dashboard with everything)
    - Overview: Kind of check, Pooling, Timeout
    - <To be added>

- GET /sites
    - HTMX integration
- GET /sites/check?url=<target url>
    - HTMX integration

- GET /sites/list (Current targets)
- POST /sites/list (Add new target)
    - {url: <target url>}
- PUT /sites/list (Update URL)
    - { id: <target id>, url: <new url>}
- DELETE /sites/list/{id} (Removes target)

- GET /metrics (export Prometheus-compatible system statistics)
```