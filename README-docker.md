# Build a local docker image
 ```
docker build . -t multilog_exporter:local
 ```

# Run multilog_exporter from local image
 ```
docker run -it -p 9144:9144 -v $(pwd)/doc/example.yaml:/mlex.yaml multilog_exporter:local
 ```
