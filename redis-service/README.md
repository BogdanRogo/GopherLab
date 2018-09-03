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


# Endpoints
## GET /ping
Is used to check connection to redis server. Does not accept payload

## POST /set-key
Is used to set a value to a specific key
### Payload example: 
```json
{
	"key": "key2",
	"value": "This is a values",
    "ttl": 500 // value in seconds
}
```

### Response example:
Success: 
```json
{
    "message": "Success",
    "status": 200,
    "ttl": 488 // time to live
}
```

## DELETE /del-keys?keys=k1&keys=k2&keys=k3
Is used to delete an array of givven keys

### Response example
```json
{
    "message": "Successfully deleted keys: [k1 k2]",
    "status": 200
}
```