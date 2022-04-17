package twitch_api_wrapper

import (
	"strings"
)

type Command struct {
	Tags    map[string]string
	Prefix  string
	Command string
	Args    []string
	Suffix  string
}

func (c *Command) String() string {
	return "Prefix: " + c.Prefix + " Command: " + c.Command + " Args: " + strings.Join(c.Args, " ") + " Suffix: " + c.Suffix
}

func ReadString(reader *strings.Reader, until byte) string {
	result := strings.Builder{}
	char, ok := reader.ReadByte()

	for ok == nil {
		if char == until {
			break
		}
		result.WriteByte(char)
		char, ok = reader.ReadByte()
	}

	return result.String()
}

func ReadTags(reader *strings.Reader) map[string]string {
	result := make(map[string]string)

	for {
		tag := ReadString(reader, '=')
		if tag == "" {
			break
		}
		value := ReadString(reader, ';')
		result[tag] = value
	}

	return result
}

func ParsePacket(packet string) *Command {
	reader := strings.NewReader(packet)
	command := Command{}
	args := make([]string, 15)
	arg := 0

	char, ok := reader.ReadByte()
	for ok == nil {
		if char == ':' && command.Prefix == "" && command.Command == "" {
			command.Prefix = ReadString(reader, ' ')
		} else if char == '@' && command.Tags == nil {
			command.Tags = ReadTags(strings.NewReader(ReadString(reader, ' ')))
		} else if command.Command == "" {
			_, err := reader.Seek(-1, 1)
			if err != nil {
				continue
			}

			command.Command = ReadString(reader, ' ')
		} else if char == ':' {
			command.Suffix = ReadString(reader, '\r')
		} else {
			_, err := reader.Seek(-1, 1)
			if err != nil {
				continue
			}

			args[arg] = ReadString(reader, ' ')
			arg++
		}
		char, ok = reader.ReadByte()
	}

	command.Args = args
	return &command
}
