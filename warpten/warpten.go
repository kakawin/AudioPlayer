// 播放器的入口文件， 会fork出一个进程启动服务器守护程序
package main

import (
	gc "code.google.com/p/goncurses"
	"fmt"
	"log"
	"os"
	"os/exec"
	"warpten/client"
)

func main() {
	// 启动服务器守护程序
	cmd := exec.Command("warpten-daemon", "-d")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// 启动gui
	stdscr, err := gc.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer gc.End()

	gc.StartColor()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_CYAN, gc.C_BLACK)

	var cli *client.WarptenCli = client.NewWarptenCli("unix", "/tmp/warpten.sock")
	version := cli.CmdVersion()

	y, x := stdscr.MaxYX()
	header := "Warpten Player"
	stdscr.ColorOn(2)
	stdscr.MovePrint(0, (x/2)-(len(header)/2), header+" v"+version)
	stdscr.ColorOff(2)

	menuwin, _ := gc.NewWindow(y-2, x, 1, 0)
	menuwin.Keypad(true)
	stdscr.Refresh()
	menuwin.Refresh()

	for {
		gc.Update()
		if ch := menuwin.GetChar(); ch == 'q' {
			return
		}
	}

	// 关闭服务器守护程序
	cmd.Process.Signal(os.Interrupt)
}
