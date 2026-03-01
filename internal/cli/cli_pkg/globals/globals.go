package globals

type OutputFormat int

const (
	OutputKV OutputFormat = iota
	OutputJSON
	OutputJSONColorless
	OutputCollapsed
	OutputYAML
	OutputTOML
)

func (o OutputFormat) String() string {
	return [...]string{
		"kv",
		"json",
		"json-colorless",
		"collapsed",
		"yaml",
		"toml",
	}[o]
}

func ParseOutputFormat(s string) OutputFormat {
	switch s {
	case "json":
		return OutputJSON
	case "json-colorless":
		return OutputJSONColorless
	case "collapsed":
		return OutputCollapsed
	case "yaml":
		return OutputYAML
	case "toml":
		return OutputTOML
	default:
		return OutputKV
	}
}
