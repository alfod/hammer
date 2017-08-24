package string

import "bytes"

func UnderLineToCamel(string string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(string); i++ {
		if i == 0 {
			buffer.WriteByte(string[i] - 32)
			continue
		}
		if string[i] == '_' {
			i += 1
			buffer.WriteByte(string[i] - 32)
			continue
		}
		buffer.WriteByte(string[i])
	}

	//fmt.Println(buffer.String())
	return buffer.String()
}

/*

 */
func ToCamel(string string, exclude_head bool) string {
	var buffer bytes.Buffer
	for i, k := range string {
		if string[i] == '_' && i != 0 {
			i += 1
			buffer.WriteByte(byte(k - 32))
			continue
		}
		buffer.WriteByte(string[i])
	}

	//fmt.Println(buffer.String())
	return buffer.String()
}
