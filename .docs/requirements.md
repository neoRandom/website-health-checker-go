```
API:

- GET / (Dashboard with everything)

- GET /sites?=<number>&offset=<number>
- GET /sites/check?url=<target url>

- GET /sites/list (Current targets)
- POST /sites/list (Add new target)
    - {url: <target url>}
- PUT /sites/list (Update URL)
    - { id: <target id>, url: <new url>}
- DELETE /sites/list/{id} (Removes target)

- GET /metrics (export Prometheus-compatible system statistics)
```
