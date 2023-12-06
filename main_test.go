package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pschlump/filelib"
)

var help_txt = `NAME:
   go-pandoc - A server for pandoc command

USAGE:
   go-pandoc [global options] command [command options] [arguments...]

COMMANDS:
   run      run pandoc service
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
`

func Test_main(t *testing.T) {
	// os.Args = []string{"./go-pandoc", "run", "-c", "app.conf"}
	os.Args = []string{"./go-pandoc", "help"}
	os.Mkdir("./out", 0755)
	out, err := filelib.Fopen("./out/help.txt", "w")
	if err != nil {
		t.Fatalf("Unable to open './out/help.txt', error:%s\n", err)
	}
	stdout := os.Stdout
	os.Stdout = out
	main()
	out.Close()
	os.Stdout = stdout

	buf1, err := ioutil.ReadFile("./out/help.txt")
	if err != nil {
		t.Fatalf("Unable to open './out/help.txt', error:%s\n", err)
	}
	if string(buf1) != help_txt {
		t.Errorf("Output did not match, -->>%s<<-- vs -->>%s<--\n", buf1, help_txt)
	}
}
