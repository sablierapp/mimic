# Mimic

Mimic is a configurable web-server with a configurable behavior.

## Usage

```bash
mimic
2024/10/21 15:35:15 Application is starting... Should start in 2 seconds.
2024/10/21 15:35:17 Starting up on port 80 (started in 2 seconds)
```

```bash
 curl -v http://localhost:8080
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.5.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Mon, 21 Oct 2024 20:04:59 GMT
< Content-Length: 17
<
* Connection #0 to host localhost left intact
Mimic says hello!
```

## Configuration

```bash
mimic --help
Usage of mimic:
  -exit-code int
        The exit code of the application.
  -healthy
        If the application should be healthy. (default true)
  -healthy-after duration
        The duration after which the application will serve 200 to the /health endpoint. (default 10s)
  -port string
        Server listening port (default "80")
  -running
        If the application should be running. If set to false, the application will exit. (default true)
  -running-after duration
        The duration after which the application will serve content. (default 2s)
```
