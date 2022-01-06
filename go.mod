module github.com/hoffie/multilog_exporter

go 1.17

require (
	github.com/fstab/grok_exporter v0.2.9-0.20200921195934-c2f8f34a8f6b
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.3.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	golang.org/x/exp v0.0.0-20200917184745-18d7dbdd5567 // indirect
	golang.org/x/sys v0.0.0-20200918174421-af09f7315aff // indirect
)

replace github.com/fstab/grok_exporter => github.com/hoffie/grok_exporter v0.2.9-0.20220106073126-d462cb52eaa4
