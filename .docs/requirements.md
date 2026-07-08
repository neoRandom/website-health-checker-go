```
API:

- GET / (Dashboard with everything)
    - Header Status: Total sites, current OK, current DOWN (4xx or 5xx), current WATCH (StatusCode != ExpectedStatusCode)
    - Overview: Method, Pooling, Timeout
    - Incident Log: host, method, endpoint, response time MS, checked at, status code
    - Site list: host, description, status, last response time MS, last check, is secure, expected status code
    - Site status: 
        - site information: host, endpoint, description, expected status code
        - last 100 results
        - a horizontal list of vertical bars, with the height of the bars representing the response time, and the color representing the OK(green), DOWN (red), or WATCH (yellow).

- GET /sites (CRUD for "/sites/list")

- GET /sites/list (Current targets)
- POST /sites/list (Add new target)
    - {url: <target url>}
- PUT /sites/list (Update URL)
    - { id: <target id>, url: <new url>}
- DELETE /sites/list/{id} (Removes target)

- GET /metrics (export Prometheus-compatible system statistics)
```
