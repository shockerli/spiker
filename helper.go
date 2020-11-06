package spiker

import (
	"encoding/json"
	"errors"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// ParseNumber parses the number in the string
func ParseNumber(value string) (float64, error) {
	var index int
	var number []rune
	r, size := utf8.DecodeRuneInString(value[index:])
	if size > 0 && (r == '-' || unicode.IsDigit(r)) {
		number = append(number, r)
		index++
		r, size = utf8.DecodeRuneInString(value[index:])
		for size > 0 && unicode.IsDigit(r) {
			number = append(number, r)
			index++
			r, size = utf8.DecodeRuneInString(value[index:])
		}
	}

	if size > 0 && r == '.' && len(number) > 0 && string(number) != "-" {
		number = append(number, '.')
		index++
		r, size = utf8.DecodeRuneInString(value[index:])
		for size > 0 && unicode.IsDigit(r) {
			number = append(number, r)
			index++
			r, size = utf8.DecodeRuneInString(value[index:])
		}
	}

	if len(number) < 1 || string(number) == "-" {
		return 0, errors.New("NOT A VALID NUMERIC STRING")
	}

	num, err := strconv.ParseFloat(string(number), 64)

	return num, err
}

// IsNumber is a number string
func IsNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

// IsTrue a variable is true
func IsTrue(value interface{}) bool {
	switch value := value.(type) {
	case string:
		return len(value) > 0

	case int:
		return value != 0

	case float64:
		return value != 0

	case bool:
		return value

	case ValueList:
		return len(value) > 0

	case ValueMap:
		return len(value) > 0
	}

	return false
}

// Interface2String convert interface{} to string
func Interface2String(inter interface{}) string {
	switch inter := inter.(type) {
	case string:
		return inter
	case int:
		return strconv.Itoa(inter)
	case float64:
		return strconv.FormatFloat(inter, 'f', -1, 64)
	case bool:
		if inter {
			return "1"
		}
		return ""

	// just for compare
	case ValueList, ValueMap:
		js, _ := json.Marshal(inter)
		return string(js)
	}

	return ""
}

// Interface2Float64 convert interface{} to float64
func Interface2Float64(inter interface{}) float64 {
	switch inter := inter.(type) {
	case string:
		num, _ := ParseNumber(inter)
		return num
	case int:
		return float64(inter)
	case float64:
		return inter
	case bool:
		if inter {
			return 1
		}
		return 0
	}

	return 0
}
