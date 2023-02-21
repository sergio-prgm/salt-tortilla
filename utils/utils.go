package utils

import (
	"fmt"
	"strings"
)

func PrintSlice(s []string) string {
	str := strings.Builder{}
	str.WriteString("[\n")
	for _, v := range s {
		vs := fmt.Sprintf("\t%s,\n", v)
		str.WriteString(vs)
	}
	str.WriteString("]")
	return str.String()
}
