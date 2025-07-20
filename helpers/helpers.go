package helpers

import (
	"fmt"
	"math/rand/v2"
	"net"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	MinTitleUpdateIntervalMinutes = 1
	MaxTitleUpdateIntervalMinutes = 1440

	VarNamePlaceholderPrefix = "$$"
	VariablePlaceholderValue = "-"
)

func GenerateCsrfToken(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, length)
	for i := range length {
		res[i] = chars[rand.IntN(len(chars))]
	}
	return string(res)[:length]
}

func GenerateVarPlaceholderString(varName string) string {
	return fmt.Sprintf("%v%v", VarNamePlaceholderPrefix, varName)
}

func GetVarNameFromPlaceholderString(placeholderString string) string {
	return strings.Replace(placeholderString, VarNamePlaceholderPrefix, "", 1)
}

func GetStringReplacerFromMap(m map[string]string, allowEmptyReplacements bool, allowSubstrings bool) (*strings.Replacer, error) {
	replacementList := []string{}
	replaceFromKeys := make([]string, 0, len(m))

	for replaceFrom, replaceTo := range m {
		if replaceFrom == "" {
			return nil, fmt.Errorf("mapping of %q to %q is not valid - replaceFrom must not be empty", replaceFrom, replaceTo)
		}
		if !allowEmptyReplacements && replaceTo == "" {
			return nil, fmt.Errorf("mapping of %q to %q is not valid - they must both be populated", replaceFrom, replaceTo)
		}
		replacementList = append(replacementList, replaceFrom, replaceTo)
		replaceFromKeys = append(replaceFromKeys, replaceFrom)
	}
	// check if any replaceFrom is a substring of any other replaceFrom
	if !allowSubstrings {
		for _, k := range replaceFromKeys {
			for _, j := range replaceFromKeys {
				if k == j {
					continue
				}
				if strings.Contains(k, j) {
					return nil, fmt.Errorf("invalid replacement map: key %q contains key %q as a substring", k, j)
				}
			}
		}
	}
	return strings.NewReplacer(replacementList...), nil
}

// convert seconds to HH:MM::SS
func GetTimeStringFromSeconds(totalSeconds int) string {
	duration := time.Duration(totalSeconds) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func PortInUse(hostAndPort string) bool {
	listener, err := net.Listen("tcp", hostAndPort)
	if err != nil {
		return true
	}
	listener.Close()
	return false
}

func OpenUrlInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("powershell", "-Command", fmt.Sprintf("Start-Process '%s'", url)).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform - cannot open browser")
	}
	return err
}

func NumFieldsInStruct(val any) (int, error) {
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Struct {
		return t.NumField(), nil
	}
	return -1, fmt.Errorf("%v is not a struct", val)
}

// only works for homogenous structs
func GenerateMapFromHomogenousStruct[ParentType any, FieldValueType any](strct ParentType) map[string]FieldValueType {
	fields := reflect.TypeOf(strct)
	vals := reflect.ValueOf(strct)
	res := map[string]FieldValueType{}
	for idx := range vals.NumField() {
		fieldName := fields.Field(idx).Name
		value := vals.Field(idx).Interface().(FieldValueType)
		res[fieldName] = value
	}
	return res
}

// Returns unique variable names
func ExtractVariableNamesFromText(text string) []string {
	fmtStr := `%s(\w+)`
	escapedPrefix := regexp.QuoteMeta(VarNamePlaceholderPrefix)
	r := regexp.MustCompile(fmt.Sprintf(fmtStr, escapedPrefix))
	matches := r.FindAllStringSubmatch(text, -1)
	variableNamesMap := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			variableNamesMap[match[1]] = struct{}{}
		}
	}
	variableNamesSlice := make([]string, 0, len(variableNamesMap))
	for v := range variableNamesMap {
		variableNamesSlice = append(variableNamesSlice, v)
	}
	return variableNamesSlice
}
