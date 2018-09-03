# URL Shortener

Service description

```bash
                           - path /                --> ???
                         /
user --> url-shortener < - - path /short&url={url} --> receive a short url
                         \
                           - path /r/{key}         --> 301 to original url
```