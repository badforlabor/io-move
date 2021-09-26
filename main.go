/**
 * Auth :   liubo
 * Date :   2021/9/26 12:04
 * Comment:
 */

package main

import (
	"flag"
	"fmt"
	"time"
)

var src = flag.String("src", "", "src=./xxx")
var dst = flag.String("dst", "", "dst=./xxx")
var action = flag.String("action", "move", "move, del")

func main() {
	flag.Parse()

	if len(*action) == 0 {
		return
	}

	switch *action {
	case "move":
		doMove()
	}

}

func doMove() {
	fmt.Println("move args=", *src, *dst)
	if len(*src) == 0 || len(*dst) == 0 {
		return
	}

	var now = time.Now()
	var count = 0
	var e = MovePath(*src, *dst, func(s string, e error) {
		count++
		if time.Now().Sub(now) > time.Second {
			fmt.Println("processed-count", count)
			now = time.Now()
		}
		if e != nil {
			fmt.Println("processed error:", s, e.Error())
		}
	})
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println("done...")
}
