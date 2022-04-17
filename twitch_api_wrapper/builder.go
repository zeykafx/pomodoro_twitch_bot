package twitch_api_wrapper

import "strings"

func (c *Command) Build() string {
	builder := strings.Builder{}
	if c.Prefix != "" {
		builder.WriteString(":" + c.Prefix + " ")
	}
	builder.WriteString(c.Command)
	for _, arg := range c.Args {
		if arg != "" {
			builder.WriteString(" " + arg)
		}
	}
	if c.Suffix != "" {
		builder.WriteString(" :" + c.Suffix)
	}
	return builder.String()
}
