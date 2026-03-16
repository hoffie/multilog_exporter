# Build a local docker image
```
docker build . -t yopur-repo/multilog-exporter:your-tag
```

# Run multilog_exporter
```
docker run -it -p 9144:9144 \
  -v /var/log/:/logs/:ro \
  -v $(pwd)/doc/example.yaml:/mlex.yaml:ro \
  your-repo/multilog-exporter:your-tag
```
