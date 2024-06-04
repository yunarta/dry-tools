module github.com/yunarta/dry-tools/dry-nexus-upload

go 1.21

require (
	github.com/spf13/cobra v1.8.0
	github.com/yunarta/dry-tools/dry-config v1.0.0
//github.com/yunarta/dry-tools/dry-config v1.0.0
)

replace github.com/yunarta/dry-tools/dry-config v1.0.0 => ../dry-config

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
