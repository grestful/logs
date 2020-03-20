package logs

import (
	"fmt"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	SetConsole(ConsoleConfig{
		Enable: true,
		Level:  "DEBUG",
	})
	fmt.Println(Global["default"].LogWriter)
	SetDefaultLog(Global["stdout"])
	Info("%s %s %s", "1", " 2222", "  333333   !!!")

	time.Sleep(time.Second)
}