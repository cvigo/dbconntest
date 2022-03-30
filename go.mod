module dbconntest

go 1.17

require (
	github.com/godror/godror v0.32.0
	github.com/ibmdb/go_ibm_db v0.3.5
	github.com/montanaflynn/stats v0.6.6
	github.com/spf13/cobra v1.4.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/godror/knownpb v0.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/ibmdb/go_ibm_db v0.3.5 => scm.live.es.nextgen.igrupobbva/connectors/go_ibm_db v0.3.5