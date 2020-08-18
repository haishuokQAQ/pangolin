package conf

type Formatter string

const (
	JSONFormater    Formatter = "json"
	ConsoleFormater Formatter = "console"
)

func isValidFormatter(f Formatter) bool {
	return (f == JSONFormater) ||
		(f == ConsoleFormater)
}
