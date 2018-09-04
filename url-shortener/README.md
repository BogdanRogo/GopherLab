# URL Shortener

Service description

```bash
                           - path /                --> index.html
                         /
user --> url-shortener < - - path /short&url={url} --> receive a short url
                         \
                           - path /r/{key}         --> 301 to original url
```

## Examples

```bash
# Method GET a URL shorted
curl http://localhost:8081/short?url=https://medium.com/metrosystemsro/gitops-with-weave-flux-40997e929254
{"url":"localhost:8081/r/2600343750"}

# use a short url
curl http://localhost:8081/r/2600343750
<a href="https://medium.com/metrosystemsro/gitops-with-weave-flux-40997e929254">Moved Permanently</a>.

```