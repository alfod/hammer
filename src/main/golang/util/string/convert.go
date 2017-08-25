package string

import "bytes"

//convert string from under_line format to camel format
func toCamel(under_line string, include_head bool) string {
	var buffer bytes.Buffer
	for i := 0; i < len(under_line); i++ {
		if i == 0 && include_head {
			buffer.WriteByte(under_line[i] - 32)
			continue
		}
		if under_line[i] == '_' {
			i += 1
			buffer.WriteByte(under_line[i] - 32)
			continue
		}
		buffer.WriteByte(under_line[i])
	}

	//fmt.Println(buffer.String())
	return buffer.String()
}

func ToLowerCamel(under_line string) string {
	return toCamel(under_line, false)
}

func ToUpperCamel(under_line string) string {
	return toCamel(under_line, true)
}

