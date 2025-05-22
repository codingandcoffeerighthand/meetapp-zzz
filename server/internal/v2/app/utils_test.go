package app

import (
	"fmt"
	"testing"
	"time"
)

var cnt int

func pull(roomID string) error {
	cnt++
	fmt.Printf("  ➤ roomID: %s, exec #%d at %s\n", roomID, cnt, time.Now().Format("15:04:05.000"))
	return nil
}

func TestHandlePUllRoom(t *testing.T) {
	f := func() error {
		return pull("123")
	}
	deb := HandlePullRoom(f, 1000*time.Millisecond)

	ch1 := <-deb() // trả channel ch1
	fmt.Println("ch1", ch1)
	ch2 := <-deb() // trả channel ch1
	fmt.Println("ch2", ch2)
	// time.Sleep(100 * time.Millisecond)
	// ch2 := <-deb() // ch1 sẽ bị huỷ, chỉ ch2 có cơ hội
	// fmt.Println("ch2", ch2)
	// time.Sleep(100 * time.Millisecond)
	// ch3 := <-deb() // chỉ ch3 thực sự chạy
	// fmt.Println("ch3", ch3)
	// ch1, ch2 không nhận gì (deadlock nếu bạn <- mà không timeout)
	// ch3 sẽ nhận kết quả
}
