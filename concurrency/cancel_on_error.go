package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
)

func main() {
	// สร้าง context ที่สามารถยกเลิกได้ (cancelable context)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ปิด context เมื่อ main function จบ
	var wg sync.WaitGroup

	mock := map[int]string{
		1: "1 val",
		2: "2 val",
		3: "3 val",
		4: "4 val",
	}

	for index, v := range mock {
		// ป้องกันปัญหา variable capture ใน goroutine
		index, v := index, v
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := doSomeRoutine(ctx, index, len(mock), v); err != nil {
				log.Println("cancel go routine", err)
				cancel() // cancel ทุก goroutine ที่ใช้ ctx นี้
			}

		}()
	}

	wg.Wait()
	log.Println("all routines finished")
}

func doSomeRoutine(ctx context.Context, k, length int, v string) error {
	select {
	case <-ctx.Done():
		// ถ้า context ถูกยกเลิกแล้ว ให้ return error ทันที
		return fmt.Errorf("routine %d canceled: %w", k, ctx.Err())
	default:
		log.Printf("routine %d working on: %s", k, v)
	}
	// สุ่ม error เพื่อจำลองว่ามีบาง routine ล้มเหลว
	if rand.IntN(length) == 0 {
		return fmt.Errorf("routine %d failed", k)
	}
	log.Printf("routine %d finished successfully", k)
	return nil
}
