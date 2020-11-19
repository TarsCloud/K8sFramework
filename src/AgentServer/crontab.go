package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// 定时器，启动的时候执行一次，以后每天晚上12点执行
func StartCronTabTask() {
	// 清理log和core
	dirs := []string{"/usr/local/app/tars/app_log/", "/usr/local/app/tars/app_log/"}
	patterns := []string{"*.log", "core.*"}
	walkDirRemove(dirs, patterns)
	// .....其他
}

func walkDirRemove(dirs []string, patterns []string) {
	if len(dirs) != len(patterns) {
		fmt.Printf("dirs size: %d != patterns size: %d\n", len(dirs), len(patterns))
	}

	rmFunc := func() {
		printExistFlag := true
		for {
			now := time.Now()

			for i := 0; i < len(dirs); i++ {
				exist, err := pathExists(dirs[i])
				if err != nil {
					fmt.Printf("dir: %s exists err: %s\n", dirs[i], err)
					continue
				}
				if !exist {
					if printExistFlag {
						fmt.Printf("dir: %s do not existed.\n", dirs[i])
					}
					continue
				}

				str := fmt.Sprintf("find %s -mtime +5 -name \"%s\" | xargs rm -f", dirs[i], patterns[i])
				cmd := exec.Command("/bin/sh", "-c", str)

				var out bytes.Buffer
				cmd.Stdout = &out

				err = cmd.Run()
				if err != nil {
					fmt.Printf("error to execute \"%s\", error: %s\n", str, err)
					continue
				}
				fmt.Printf("succ. to execute \"%s\", msg: %s\n", str, out.String())
			}

			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

			printExistFlag = false
			fmt.Printf("crond finish now: %s, next execute time: %s\n", now, next)

			t := time.NewTimer(next.Sub(now))
			<- t.C
		}
	}

	go rmFunc()
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}