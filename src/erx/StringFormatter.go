package erx

import (
	"strings"
	"strconv"
	"fmt"
	"os"
)

type StringFormatter struct {
	indent string
}

func NewStringFormatter(indent string) *StringFormatter {
	formatter := new(StringFormatter)
	formatter.indent = indent
	return formatter
}

func (f *StringFormatter) Format(err Error) string {
	return f.formatLevel(err, 0)
}

func (f *StringFormatter) formatLevel(err Error, level int) string {
	result := ""
	result += strings.Repeat(f.indent, level)
	result += err.Message()
	result += "\n"
	result += strings.Repeat(f.indent, level)
	result += err.File() + ": " + strconv.Itoa(err.Line())
	result += "\n"
	level++
	if len(err.Variables())>0 {
		result += strings.Repeat(f.indent, level)
		result += "Scope variables:\n"
		for name, val := range err.Variables() {
			result += strings.Repeat(f.indent, level+1)
			result += name + "\t: "
			switch i := val.(type) {
				case string :
					result += i
				case fmt.Stringer :
					result += i.String()
				default :
					result += fmt.Sprint(i)
			}
			result += "\n"
		}
	}
	
	curErr := err.Errors().Front()
	if curErr!=nil {
		result += strings.Repeat(f.indent, level)
		result += "Scope errors:\n"
		for curErr!=nil {
			switch i := curErr.Value.(type) {
				case Error :
					result += f.formatLevel(i, level+1)
				case os.Error :
					result += strings.Repeat(f.indent, level+1)
					result += i.String()
				default :
					result += "???\n"
			}
			curErr = curErr.Next()
		}
	}
	return result
}
