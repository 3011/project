package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func downloadFile() {
	downloadUrl := "https://github.com/3011/FileLibBot-go/releases/latest/download/FileLibBot-go"
	resp, err := http.Get(downloadUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, err := os.Create("FileLibBot-go")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

func getPid() string {
	f, err := os.Open("test.txt")
	if err != nil {
		return ""
	}
	defer f.Close()

	br := bufio.NewReader(f)
	pid, _, _ := br.ReadLine()

	return string(pid)
}

func writePid(pid int) {
	f, _ := os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	f.WriteString(strconv.Itoa(pid))
}

func main() {
	pid := getPid()

	if pid != "" {
		fmt.Printf("kill Pid: %v\n", pid)
		cmd := exec.Command("kill", pid)
		cmd.Start()
		cmd.Wait()
	}

	downloadFile()

	exec.Command("chmod", "u+x", "./FileLibBot-go").Start()

	cmd := exec.Command("./FileLibBot-go")
	cmd.Env = os.Environ()
	cmd.Start()

	writePid(cmd.Process.Pid)
	fmt.Printf("start Pid: %v\n", cmd.Process.Pid)
}
