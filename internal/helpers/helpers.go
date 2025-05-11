package helpers

import (
	"fmt"
	"math/rand/v2"
	"net"
	"os/exec"
	"reflect"
	"runtime"
)

func GenerateCsrfToken(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = chars[rand.IntN(len(chars))]
	}
	return string(res)[:length]
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
