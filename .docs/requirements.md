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

- GET /metrics (show system statistics: total uptime, how many targets are actively being monitored, and the memory footprint)
    - {
        uptime_seconds: int
        monitored_targets: int
        memory_alloc_bytes: int
        memory_sys_bytes: int
        total_requests: int
        requests_per_minute: float
    }
```
