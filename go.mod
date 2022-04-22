module dbconntest

go 1.17

require (
	github.com/godror/godror v0.33.0
	github.com/ibmdb/go_ibm_db v0.3.0
	github.com/montanaflynn/stats v0.6.6
	github.com/spf13/cobra v1.4.0
	go.uber.org/zap v1.21.0
)

require github.com/go-logr/logr v1.2.3 // indirect

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/godror/knownpb v0.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/ivanpirog/coloredcobra v1.0.0
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20220330033206-e17cdc41300f // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/godror/godror => github.com/cvigo/godror v0.33.1-0.20220421115748-171301d5ea0a

replace github.com/ibmdb/go_ibm_db => scm.live.es.nextgen.igrupobbva/connectors/go_ibm_db v0.3.6
