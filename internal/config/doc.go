// Package config provides configuration parsing and validation for logslice.
//
// It exposes two main entry points:
//
//   - [DefaultConfig] returns a Config with sensible defaults.
//   - [ParseFlags] parses a slice of CLI arguments (typically os.Args[1:])
//     into a validated Config, resolving timestamp strings via the parser
//     package.
//
// Config validation is separated from flag parsing so that configs built
// programmatically (e.g. in tests) can also be validated with [Config.Validate].
package config
