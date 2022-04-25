# Build a local docker image
docker build . -t multilog_exporter:local

# Run multilog_exporter from local image
docker run -it multilog_exporter:local
docker run -it -p 9144:9144 \
	-v $(pwd)/doc/:/config/ \
	multilog_exporter:local \
	./multilog_exporter \
	--config.file /config/example.yaml