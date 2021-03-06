package logs

import (
	"fmt"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	var bt = []byte{91,50,48,50,48,45,48,51,45,50,48,84,49,52,58,51,52,58,51,49,46,55,49,55,43,48,56,48,48,93,91,68,69,66,85,71,93,91,98,111,115,115,95,80,111,100,99,97,115,116,93,32,58,105,112,58,32,58,58,49,44,32,109,101,116,104,111,100,58,32,
		71,69,84,44,32,112,97,116,104,58,32,47,44,32,99,111,100,101,58,32,52,48,52,44,32,97,103,101,110,116,58,32,77,111,122,105,108,108,97,47,53,46,48,32,40,87,105,110,100,111,119,115,32,78,84,32,49,48,46,48,59,32,87,105,110,54,52,59,32,120,54,52,41,32,65,112,112,108,101,87,101,98,75,105,116,47,53,51,55,46,51,54,32,40,
		75,72,84,77,76,44,32,108,105,107,101,32,71,101,99,107,111,41,32,67,104,114,111,109,101,47,55,57,46,48,46,51,57,52,53,46,49,49,55,32,83,97,102,97,114,105,47,53,51,55,46,
		51,54,44,32,101,114,114,111,114,32,10}
	fmt.Println(string(bt[:]))
	SetFile(FileConfig{
		Enable:   true,
		Category: "file",
		Level:    "DEBUG",
		Filename:  "D:\\work\\cast.log",
		Rotate:   true,
		Daily:    true,
		Sanitize: false,
	})
	//writer := GetLogger("file")
	SetDefaultLog(Global["file"])
	Info("%s %s %s", "1", " 2222", "  333333   !!!")

	time.Sleep(time.Second)
}