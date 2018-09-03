# Redis Service

## redis launch in Docker and Test

```bash
docker run --rm -p 6379:6379 --name test-redis -d redis
telnet localhost 6379
Trying ::1...
Connected to localhost.
Escape character is '^]'.
MONITOR
+OK
QUIT
+OK
Connection closed by foreign host.
```
