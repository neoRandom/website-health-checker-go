```
API:

- GET / (Dashboard with everything)

- GET /sites?=<number>&offset=<number>
- GET /sites/check?url=<target url>

- GET /sites/list (Current targets)
- POST /sites/list?url=<target url> (Add new target)
- PUT /sites/list?id=<target id>?url=<new url> (Update URL)
- DELETE /sites/list?id=<target id> (Removes target)

- GET /metrics (show system statistics: total uptime, how many targets are actively being monitored, and the memory footprint)
```