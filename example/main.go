package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jiansoft/robin/v2"
)

// timingTolerance 計時驗證的容許誤差（考慮 OS 排程與 timer 精度）
// timingTolerance allowed deviation for timing verification (accounts for OS scheduling & timer precision)
const timingTolerance = 150 * time.Millisecond

// gameEvent 用於驗證 TypedChannel 結構體泛型的測試型別
// gameEvent is a struct used to verify TypedChannel with non-primitive types
type gameEvent struct {
	Source string
	Target string
	Damage int
	Time   time.Time
}

var (
	passCount atomic.Int32
	failCount atomic.Int32
)

// verify 記錄並輸出單一驗證結果
// verify logs and prints a single verification result
func verify(name string, ok bool, format string, args ...any) {
	detail := fmt.Sprintf(format, args...)
	if ok {
		passCount.Add(1)
		log.Printf("  [PASS] %s — %s", name, detail)
	} else {
		failCount.Add(1)
		log.Printf("  [FAIL] %s — %s", name, detail)
	}
}

// nearDuration 檢查實際耗時是否在 expected ± tol 範圍內
// nearDuration checks if actual duration is within expected ± tol
func nearDuration(actual, expected, tol time.Duration) bool {
	diff := actual - expected
	if diff < 0 {
		diff = -diff
	}
	return diff <= tol
}

// waitTimeout 等待 WaitGroup 完成或逾時；回傳 true 表示在時限內完成
// waitTimeout waits for WaitGroup completion or timeout; returns true if completed in time
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

// main is the program entry point for the verification example.
// main 是此驗證範例程式的進入點。
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("========================================")
	log.Println(" Robin Library — Verification Examples")
	log.Println("========================================")

	// 同步驗證（即時完成）
	// Synchronous verification (completes instantly)
	verifyAbs()
	verifyConcurrentQueue()
	verifyConcurrentStack()
	verifyConcurrentBag()

	// 非同步驗證（需等待排程執行）
	// Asynchronous verification (waits for scheduled execution)
	verifyFiber()
	verifyCron()
	verifyCronScheduler()
	verifyChannel()
	verifyTypedChannel()

	// 展示長間隔排程 API（僅展示用法，無法在短時間內驗證）
	// Showcase long-interval scheduling APIs (usage demo only, cannot verify in short time)
	showSchedulingExamples()

	// 結果摘要 / Results summary
	log.Println("========================================")
	p, f := passCount.Load(), failCount.Load()
	if f == 0 {
		log.Printf(" ALL PASSED: %d/%d", p, p+f)
	} else {
		log.Printf(" RESULTS: %d passed, %d FAILED", p, f)
	}
	log.Println("========================================")
}

// ================================================================
// Abs 驗證 / Abs Verification
// ================================================================

// verifyAbs runs a grouped verification flow and reports pass/fail results.
// verifyAbs 會執行一組驗證流程並輸出通過/失敗結果。
func verifyAbs() {
	log.Println("\n--- Abs — 泛型絕對值 / Generic absolute value ---")

	// Abs[T] 支援 int, int8~64, float32, float64
	// Abs[T] supports int, int8~64, float32, float64
	verify("Abs(-42)", robin.Abs(-42) == 42,
		"got %d, want 42", robin.Abs(-42))
	verify("Abs(42)", robin.Abs(42) == 42,
		"got %d, want 42", robin.Abs(42))
	verify("Abs(-3.14)", robin.Abs(-3.14) == 3.14,
		"got %.2f, want 3.14", robin.Abs(-3.14))
	verify("Abs(0)", robin.Abs(0) == 0,
		"got %d, want 0", robin.Abs(0))
}

// ================================================================
// ConcurrentQueue 驗證 / ConcurrentQueue Verification
// ================================================================

// verifyConcurrentQueue runs a grouped verification flow and reports pass/fail results.
// verifyConcurrentQueue 會執行一組驗證流程並輸出通過/失敗結果。
func verifyConcurrentQueue() {
	log.Println("\n--- ConcurrentQueue[T] — 執行緒安全 FIFO 佇列 / Thread-safe FIFO queue ---")

	// NewConcurrentQueue 建立泛型佇列（環形緩衝區實作）
	// NewConcurrentQueue creates a generic queue (ring buffer implementation)
	q := robin.NewConcurrentQueue[string]()

	// 空佇列操作安全，TryPeek/TryDequeue 回傳 (zero, false)
	// Empty queue is safe; TryPeek/TryDequeue return (zero, false)
	_, ok := q.TryPeek()
	verify("Queue.TryPeek(empty)", !ok, "ok=%v, want false", ok)
	_, ok = q.TryDequeue()
	verify("Queue.TryDequeue(empty)", !ok, "ok=%v, want false", ok)
	verify("Queue.Len(empty)", q.Len() == 0, "len=%d, want 0", q.Len())

	// Enqueue 加入元素 + Len 驗證數量
	// Enqueue adds elements + Len verifies count
	q.Enqueue("A")
	q.Enqueue("B")
	q.Enqueue("C")
	verify("Queue.Enqueue+Len", q.Len() == 3, "len=%d, want 3", q.Len())

	// TryPeek 查看前端元素但不移除
	// TryPeek peeks at front without removing
	val, ok := q.TryPeek()
	verify("Queue.TryPeek", ok && val == "A", "val=%q ok=%v, want A/true", val, ok)
	verify("Queue.Len(after peek)", q.Len() == 3, "len=%d, want 3 (unchanged)", q.Len())

	// TryDequeue 取出並移除前端元素（FIFO 順序）
	// TryDequeue removes front element (FIFO order)
	val, _ = q.TryDequeue()
	verify("Queue.TryDequeue(1st)", val == "A", "val=%q, want A", val)
	val, _ = q.TryDequeue()
	verify("Queue.TryDequeue(2nd)", val == "B", "val=%q, want B", val)

	// ToArray 複製剩餘元素到新 slice
	// ToArray copies remaining elements to a new slice
	arr := q.ToArray()
	verify("Queue.ToArray", len(arr) == 1 && arr[0] == "C", "arr=%v, want [C]", arr)

	// Clear 清空佇列並重設容量
	// Clear removes all elements and resets capacity
	q.Clear()
	verify("Queue.Clear", q.Len() == 0, "len=%d, want 0", q.Len())
}

// ================================================================
// ConcurrentStack 驗證 / ConcurrentStack Verification
// ================================================================

// verifyConcurrentStack runs a grouped verification flow and reports pass/fail results.
// verifyConcurrentStack 會執行一組驗證流程並輸出通過/失敗結果。
func verifyConcurrentStack() {
	log.Println("\n--- ConcurrentStack[T] — 執行緒安全 LIFO 堆疊 / Thread-safe LIFO stack ---")

	s := robin.NewConcurrentStack[int]()

	// 空堆疊操作安全
	// Empty stack operations are safe
	_, ok := s.TryPeek()
	verify("Stack.TryPeek(empty)", !ok, "ok=%v, want false", ok)
	_, ok = s.TryPop()
	verify("Stack.TryPop(empty)", !ok, "ok=%v, want false", ok)

	// Push 推入 + Len 驗證
	// Push elements + Len verification
	s.Push(10)
	s.Push(20)
	s.Push(30)
	verify("Stack.Push+Len", s.Len() == 3, "len=%d, want 3", s.Len())

	// TryPeek 查看頂端元素但不移除
	// TryPeek peeks at top without removing
	val, ok := s.TryPeek()
	verify("Stack.TryPeek", ok && val == 30, "val=%d ok=%v, want 30/true", val, ok)

	// TryPop 彈出頂端元素（LIFO 順序）
	// TryPop pops top element (LIFO order)
	val, _ = s.TryPop()
	verify("Stack.TryPop(1st)", val == 30, "val=%d, want 30", val)
	val, _ = s.TryPop()
	verify("Stack.TryPop(2nd)", val == 20, "val=%d, want 20", val)

	// ToArray 複製元素（LIFO 順序：頂端在前）
	// ToArray copies elements (LIFO order: top first)
	arr := s.ToArray()
	verify("Stack.ToArray", len(arr) == 1 && arr[0] == 10, "arr=%v, want [10]", arr)

	// Clear 清空堆疊
	// Clear removes all elements
	s.Clear()
	verify("Stack.Clear", s.Len() == 0, "len=%d, want 0", s.Len())
}

// ================================================================
// ConcurrentBag 驗證 / ConcurrentBag Verification
// ================================================================

// verifyConcurrentBag runs a grouped verification flow and reports pass/fail results.
// verifyConcurrentBag 會執行一組驗證流程並輸出通過/失敗結果。
func verifyConcurrentBag() {
	log.Println("\n--- ConcurrentBag[T] — 執行緒安全無序集合 / Thread-safe unordered collection ---")

	b := robin.NewConcurrentBag[float64]()

	// 空集合操作安全
	// Empty bag operations are safe
	_, ok := b.TryTake()
	verify("Bag.TryTake(empty)", !ok, "ok=%v, want false", ok)
	verify("Bag.Len(empty)", b.Len() == 0, "len=%d, want 0", b.Len())

	// Add 加入元素 + Len 驗證
	// Add elements + Len verification
	b.Add(1.1)
	b.Add(2.2)
	b.Add(3.3)
	verify("Bag.Add+Len", b.Len() == 3, "len=%d, want 3", b.Len())

	// TryTake 取出一個元素（從尾部取出，無序集合不保證順序）
	// TryTake removes and returns an element (from tail, unordered collection)
	val, ok := b.TryTake()
	verify("Bag.TryTake", ok && val == 3.3, "val=%.1f ok=%v, want 3.3/true", val, ok)
	verify("Bag.Len(after take)", b.Len() == 2, "len=%d, want 2", b.Len())

	// ToArray 複製元素到新 slice
	// ToArray copies elements to a new slice
	arr := b.ToArray()
	verify("Bag.ToArray", len(arr) == 2, "len=%d, want 2", len(arr))

	// Clear 清空集合
	// Clear removes all elements
	b.Clear()
	verify("Bag.Clear", b.Len() == 0, "len=%d, want 0", b.Len())
}

// ================================================================
// Fiber 驗證 / Fiber Verification
// ================================================================

// verifyFiber runs a grouped verification flow and reports pass/fail results.
// verifyFiber 會執行一組驗證流程並輸出通過/失敗結果。
func verifyFiber() {
	log.Println("\n--- Fiber — 纖程任務執行 / Fiber task execution ---")

	// NewGoroutineMulti — 多 goroutine 纖程：每個任務批次由新 goroutine 執行
	// NewGoroutineMulti — Multi-goroutine fiber: each task batch runs in a new goroutine

	// Enqueue — 驗證任務被執行且參數正確傳遞
	// Enqueue — verify task executes with correct parameters
	{
		gm := robin.NewGoroutineMulti()
		defer gm.Dispose()

		var wg sync.WaitGroup
		var got atomic.Value
		wg.Add(1)
		gm.Enqueue(func(msg string) {
			got.Store(msg)
			wg.Done()
		}, "hello")
		completed := waitTimeout(&wg, 2*time.Second)
		verify("GoroutineMulti.Enqueue",
			completed && got.Load() == "hello",
			"completed=%v got=%v, want hello", completed, got.Load())
	}

	// Schedule — 驗證延遲 N 毫秒後執行一次
	// Schedule — verify execution after N milliseconds delay
	{
		gm := robin.NewGoroutineMulti()
		defer gm.Dispose()

		var wg sync.WaitGroup
		var elapsed atomic.Int64
		wg.Add(1)
		start := time.Now()
		gm.Schedule(200, func() {
			elapsed.Store(int64(time.Since(start)))
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		actual := time.Duration(elapsed.Load())
		ok := completed && nearDuration(actual, 200*time.Millisecond, timingTolerance)
		verify("GoroutineMulti.Schedule(200ms)", ok,
			"delay=%v, want ~200ms", actual.Round(time.Millisecond))
	}

	// ScheduleOnInterval — 驗證重複執行次數與間隔
	// ScheduleOnInterval — verify repeated execution count and interval
	{
		gm := robin.NewGoroutineMulti()
		defer gm.Dispose()

		var count atomic.Int32
		var firstAt, thirdAt atomic.Int64
		var wg sync.WaitGroup
		wg.Add(3)
		start := time.Now()
		d := gm.ScheduleOnInterval(0, 100, func() {
			n := count.Add(1)
			switch n {
			case 1:
				firstAt.Store(int64(time.Since(start)))
				wg.Done()
			case 2:
				wg.Done()
			case 3:
				thirdAt.Store(int64(time.Since(start)))
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 3*time.Second)
		d.Dispose()
		first := time.Duration(firstAt.Load())
		third := time.Duration(thirdAt.Load())
		interval := third - first
		verify("GoroutineMulti.ScheduleOnInterval",
			completed && count.Load() >= 3 && nearDuration(interval, 200*time.Millisecond, timingTolerance),
			"count=%d, 1st→3rd interval=%v, want ~200ms", count.Load(), interval.Round(time.Millisecond))
	}

	// Disposable.Dispose — 驗證取消排程後停止執行
	// Disposable.Dispose — verify task stops after disposal
	{
		gm := robin.NewGoroutineMulti()
		defer gm.Dispose()

		var count atomic.Int32
		d := gm.ScheduleOnInterval(0, 50, func() {
			count.Add(1)
		})
		time.Sleep(200 * time.Millisecond)
		d.Dispose()
		countAtDispose := count.Load()
		time.Sleep(200 * time.Millisecond)
		countAfter := count.Load()
		// 允許最多 1 次額外執行（已在途的任務）
		// Allow at most 1 extra execution (in-flight task)
		ok := countAfter <= countAtDispose+1
		verify("Disposable.Dispose(stops)", ok,
			"at_dispose=%d, after_wait=%d", countAtDispose, countAfter)
	}

	// NewGoroutineSingle — 單 goroutine 纖程：驗證任務按入列順序串行執行
	// NewGoroutineSingle — Single-goroutine fiber: verify tasks execute serially in enqueue order
	{
		gs := robin.NewGoroutineSingle()
		defer gs.Dispose()

		var order []int
		var mu sync.Mutex
		var wg sync.WaitGroup
		wg.Add(3)
		for i := 1; i <= 3; i++ {
			n := i
			gs.Enqueue(func() {
				mu.Lock()
				order = append(order, n)
				mu.Unlock()
				wg.Done()
			})
		}
		completed := waitTimeout(&wg, 2*time.Second)
		ok := completed && len(order) == 3 && order[0] == 1 && order[1] == 2 && order[2] == 3
		verify("GoroutineSingle.order", ok,
			"order=%v, want [1 2 3]", order)
	}

	// GoroutineSingle.Schedule — 驗證延遲 N 毫秒後串行執行
	// GoroutineSingle.Schedule — verify serial execution after N milliseconds delay
	{
		gs := robin.NewGoroutineSingle()
		defer gs.Dispose()

		var wg sync.WaitGroup
		var elapsed atomic.Int64
		wg.Add(1)
		start := time.Now()
		gs.Schedule(200, func() {
			elapsed.Store(int64(time.Since(start)))
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		actual := time.Duration(elapsed.Load())
		ok := completed && nearDuration(actual, 200*time.Millisecond, timingTolerance)
		verify("GoroutineSingle.Schedule(200ms)", ok,
			"delay=%v, want ~200ms", actual.Round(time.Millisecond))
	}

	// GoroutineSingle.ScheduleOnInterval — 驗證串行重複執行
	// GoroutineSingle.ScheduleOnInterval — verify serial repeated execution
	{
		gs := robin.NewGoroutineSingle()
		defer gs.Dispose()

		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(3)
		d := gs.ScheduleOnInterval(0, 100, func() {
			if count.Add(1) <= 3 {
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 3*time.Second)
		d.Dispose()
		got := count.Load()
		verify("GoroutineSingle.ScheduleOnInterval",
			completed && got >= 3,
			"count=%d, want >=3", got)
	}

	// Fiber.Dispose — 驗證 Dispose 後 Enqueue 不會 panic
	// Fiber.Dispose — verify Enqueue after Dispose doesn't panic
	{
		gm := robin.NewGoroutineMulti()
		gm.Dispose()
		noPanic := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					noPanic = false
				}
			}()
			gm.Enqueue(func() {})
		}()
		verify("Fiber.Dispose+Enqueue", noPanic, "no panic=%v", noPanic)
	}
}

// ================================================================
// Cron 驗證 / Cron Verification (global functions)
// ================================================================

// verifyCron runs a grouped verification flow and reports pass/fail results.
// verifyCron 會執行一組驗證流程並輸出通過/失敗結果。
func verifyCron() {
	log.Println("\n--- Cron — 排程工作驗證 / Job scheduling verification ---")

	// RightNow — 立即執行，驗證延遲極短
	// RightNow — immediate execution, verify minimal delay
	{
		var wg sync.WaitGroup
		var elapsed atomic.Int64
		wg.Add(1)
		start := time.Now()
		robin.RightNow().Do(func() {
			elapsed.Store(int64(time.Since(start)))
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		actual := time.Duration(elapsed.Load())
		verify("RightNow", completed && actual < 200*time.Millisecond,
			"delay=%v, want <200ms", actual.Round(time.Millisecond))
	}

	// Delay(N) — 驗證延遲 N 毫秒後執行
	// Delay(N) — verify execution after N milliseconds
	{
		var wg sync.WaitGroup
		var elapsed atomic.Int64
		wg.Add(1)
		start := time.Now()
		robin.Delay(300).Do(func() {
			elapsed.Store(int64(time.Since(start)))
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		actual := time.Duration(elapsed.Load())
		ok := completed && nearDuration(actual, 300*time.Millisecond, timingTolerance)
		verify("Delay(300ms)", ok,
			"delay=%v, want ~300ms", actual.Round(time.Millisecond))
	}

	// Every(N).Milliseconds().Times(3) — 驗證恰好執行 3 次
	// Every(N).Milliseconds().Times(3) — verify exactly 3 executions
	{
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(3)
		robin.Every(100).Milliseconds().Times(3).Do(func() {
			if count.Add(1) <= 3 {
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 3*time.Second)
		// 等待額外時間確認不會多執行
		// Wait extra to confirm no over-execution
		time.Sleep(300 * time.Millisecond)
		got := count.Load()
		verify("Every(100ms).Times(3)", completed && got == 3,
			"count=%d, want exactly 3", got)
	}

	// Every(1).Seconds().Times(2) — 驗證執行間隔 ~1s
	// Every(1).Seconds().Times(2) — verify ~1s execution interval
	{
		var execTimes [2]atomic.Int64
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(2)
		robin.Every(1).Seconds().Times(2).Do(func() {
			n := count.Add(1)
			if n <= 2 {
				execTimes[n-1].Store(time.Now().UnixNano())
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 5*time.Second)
		if completed {
			t1 := time.Unix(0, execTimes[0].Load())
			t2 := time.Unix(0, execTimes[1].Load())
			interval := t2.Sub(t1)
			ok := nearDuration(interval, 1*time.Second, timingTolerance)
			verify("Every(1s).Times(2)", ok,
				"interval=%v, want ~1s", interval.Round(time.Millisecond))
		} else {
			verify("Every(1s).Times(2)", false, "timeout")
		}
	}

	// Until(future) — 驗證在指定時間點執行
	// Until(future) — verify execution at specified time point
	{
		var wg sync.WaitGroup
		var execNano atomic.Int64
		wg.Add(1)
		target := time.Now().Add(500 * time.Millisecond)
		robin.Until(target).Do(func() {
			execNano.Store(time.Now().UnixNano())
			wg.Done()
		})
		completed := waitTimeout(&wg, 3*time.Second)
		if completed {
			execAt := time.Unix(0, execNano.Load())
			diff := execAt.Sub(target)
			if diff < 0 {
				diff = -diff
			}
			verify("Until(future)", diff < timingTolerance,
				"target=%s, actual=%s, diff=%v",
				target.Format("15:04:05.000"), execAt.Format("15:04:05.000"),
				diff.Round(time.Millisecond))
		} else {
			verify("Until(future)", false, "timeout")
		}
	}

	// Until(past) — 過去時間不執行，Dispose 不 panic
	// Until(past) — past time should not execute, Dispose should not panic
	{
		var executed atomic.Bool
		d := robin.Until(time.Now().Add(-1 * time.Second)).Do(func() {
			executed.Store(true)
		})
		time.Sleep(300 * time.Millisecond)
		noPanic := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					noPanic = false
				}
			}()
			d.Dispose()
		}()
		verify("Until(past)", !executed.Load() && noPanic,
			"executed=%v, dispose_panic=%v", executed.Load(), !noPanic)
	}

	// Delay(N).Times(3) — 延遲後重複執行，驗證次數與間隔
	// Delay(N).Times(3) — repeat after delay, verify count and interval
	{
		var count atomic.Int32
		var firstNano, lastNano atomic.Int64
		var wg sync.WaitGroup
		wg.Add(3)
		robin.Delay(200).Times(3).Do(func() {
			now := time.Now().UnixNano()
			n := count.Add(1)
			if n == 1 {
				firstNano.Store(now)
			}
			lastNano.Store(now)
			if n <= 3 {
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 5*time.Second)
		time.Sleep(300 * time.Millisecond)
		got := count.Load()
		first := time.Unix(0, firstNano.Load())
		last := time.Unix(0, lastNano.Load())
		span := last.Sub(first)
		// 3 次執行，間隔 200ms，跨度應 ~400ms
		// 3 executions at 200ms interval, total span ~400ms
		ok := completed && got == 3 && nearDuration(span, 400*time.Millisecond, timingTolerance)
		verify("Delay(200ms).Times(3)", ok,
			"count=%d, span=%v, want 3 / ~400ms", got, span.Round(time.Millisecond))
	}

	// AfterExecuteTask — 驗證任務完成後才計算下次時間
	// AfterExecuteTask — verify next time calculated after task completes
	{
		var execTimes [2]atomic.Int64
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(2)
		robin.Every(200).Milliseconds().AfterExecuteTask().Times(2).Do(func() {
			n := count.Add(1)
			if n <= 2 {
				execTimes[n-1].Store(time.Now().UnixNano())
				// 模擬耗時任務 100ms / simulate 100ms task duration
				time.Sleep(100 * time.Millisecond)
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 5*time.Second)
		if completed {
			t1 := time.Unix(0, execTimes[0].Load())
			t2 := time.Unix(0, execTimes[1].Load())
			interval := t2.Sub(t1)
			// AfterExecuteTask: interval = 任務耗時(100ms) + 排程間隔(200ms) ≈ 300ms
			// AfterExecuteTask: interval = task_duration(100ms) + schedule_interval(200ms) ≈ 300ms
			ok := nearDuration(interval, 300*time.Millisecond, timingTolerance)
			verify("AfterExecuteTask", ok,
				"interval=%v, want ~300ms (200ms+100ms task)", interval.Round(time.Millisecond))
		} else {
			verify("AfterExecuteTask", false, "timeout")
		}
	}

	// BeforeExecuteTask — 驗證執行前即計算下次時間（預設行為）
	// BeforeExecuteTask — verify next time calculated before execution (default)
	{
		var execTimes [2]atomic.Int64
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(2)
		robin.Every(300).Milliseconds().BeforeExecuteTask().Times(2).Do(func() {
			n := count.Add(1)
			if n <= 2 {
				execTimes[n-1].Store(time.Now().UnixNano())
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 5*time.Second)
		if completed {
			t1 := time.Unix(0, execTimes[0].Load())
			t2 := time.Unix(0, execTimes[1].Load())
			interval := t2.Sub(t1)
			// BeforeExecuteTask: interval ≈ 排程間隔(300ms)，不含任務耗時
			// BeforeExecuteTask: interval ≈ schedule_interval(300ms), excludes task duration
			ok := nearDuration(interval, 300*time.Millisecond, timingTolerance)
			verify("BeforeExecuteTask", ok,
				"interval=%v, want ~300ms", interval.Round(time.Millisecond))
		} else {
			verify("BeforeExecuteTask", false, "timeout")
		}
	}

	// Job.Dispose — 驗證取消後停止執行
	// Job.Dispose — verify execution stops after cancellation
	{
		var count atomic.Int32
		job := robin.Every(50).Milliseconds().Times(100).Do(func() {
			count.Add(1)
		})
		time.Sleep(200 * time.Millisecond)
		job.Dispose()
		countAtDispose := count.Load()
		time.Sleep(200 * time.Millisecond)
		countAfter := count.Load()
		ok := countAfter <= countAtDispose+1
		verify("Job.Dispose", ok,
			"at_dispose=%d, after_wait=%d (should stop)", countAtDispose, countAfter)
	}
}

// ================================================================
// CronScheduler 驗證 / CronScheduler Verification (instance-based)
// ================================================================

// verifyCronScheduler runs a grouped verification flow and reports pass/fail results.
// verifyCronScheduler 會執行一組驗證流程並輸出通過/失敗結果。
func verifyCronScheduler() {
	log.Println("\n--- CronScheduler — 實例化排程器 / Instance-based scheduler ---")

	// NewCronScheduler — 擁有獨立 Fiber，與全域排程器互不影響
	// NewCronScheduler — owns its own Fiber, independent from global scheduler
	cs := robin.NewCronScheduler()

	// CronScheduler.RightNow — 立即執行
	// CronScheduler.RightNow — immediate execution
	{
		var wg sync.WaitGroup
		var executed atomic.Bool
		wg.Add(1)
		cs.RightNow().Do(func() {
			executed.Store(true)
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		verify("CronScheduler.RightNow", completed && executed.Load(),
			"executed=%v", executed.Load())
	}

	// CronScheduler.Delay — 延遲執行
	// CronScheduler.Delay — delayed execution
	{
		var wg sync.WaitGroup
		var elapsed atomic.Int64
		wg.Add(1)
		start := time.Now()
		cs.Delay(300).Do(func() {
			elapsed.Store(int64(time.Since(start)))
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		actual := time.Duration(elapsed.Load())
		ok := completed && nearDuration(actual, 300*time.Millisecond, timingTolerance)
		verify("CronScheduler.Delay(300ms)", ok,
			"delay=%v, want ~300ms", actual.Round(time.Millisecond))
	}

	// CronScheduler.Every + Times — 驗證執行次數
	// CronScheduler.Every + Times — verify execution count
	{
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(2)
		cs.Every(150).Milliseconds().Times(2).Do(func() {
			if count.Add(1) <= 2 {
				wg.Done()
			}
		})
		completed := waitTimeout(&wg, 3*time.Second)
		time.Sleep(300 * time.Millisecond)
		got := count.Load()
		verify("CronScheduler.Every(150ms).Times(2)", completed && got == 2,
			"count=%d, want 2", got)
	}

	// CronScheduler.Until — 指定時間執行
	// CronScheduler.Until — execute at specified time
	{
		var wg sync.WaitGroup
		var executed atomic.Bool
		wg.Add(1)
		cs.Until(time.Now().Add(300 * time.Millisecond)).Do(func() {
			executed.Store(true)
			wg.Done()
		})
		completed := waitTimeout(&wg, 2*time.Second)
		verify("CronScheduler.Until", completed && executed.Load(),
			"executed=%v", executed.Load())
	}

	// CronScheduler.Everyday — 每天排程（驗證 Job 建立成功）
	// CronScheduler.Everyday — daily schedule (verify job created successfully)
	{
		noPanic := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					noPanic = false
				}
			}()
			cs.Everyday().At(12, 0, 0).Do(func() {})
		}()
		verify("CronScheduler.Everyday", noPanic, "created without panic")
	}

	// CronScheduler.Dispose — 停止所有排程
	// CronScheduler.Dispose — stops all scheduled jobs
	{
		cs2 := robin.NewCronScheduler()
		var count atomic.Int32
		cs2.Every(50).Milliseconds().Times(100).Do(func() {
			count.Add(1)
		})
		time.Sleep(200 * time.Millisecond)
		cs2.Dispose()
		countAtDispose := count.Load()
		time.Sleep(200 * time.Millisecond)
		countAfter := count.Load()
		ok := countAfter <= countAtDispose+1
		verify("CronScheduler.Dispose", ok,
			"at_dispose=%d, after_wait=%d (should stop)", countAtDispose, countAfter)
	}

	cs.Dispose()
}

// ================================================================
// Channel 驗證 / Channel Verification (Pub/Sub)
// ================================================================

// verifyChannel runs a grouped verification flow and reports pass/fail results.
// verifyChannel 會執行一組驗證流程並輸出通過/失敗結果。
func verifyChannel() {
	log.Println("\n--- Channel — 發布/訂閱 / Pub/Sub ---")

	// Subscribe + Publish — 驗證所有訂閱者都收到訊息
	// Subscribe + Publish — verify all subscribers receive the message
	{
		ch := robin.NewChannel()
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(3)
		ch.Subscribe(func(msg string) { count.Add(1); wg.Done() })
		ch.Subscribe(func(msg string) { count.Add(1); wg.Done() })
		ch.Subscribe(func(msg string) { count.Add(1); wg.Done() })
		ch.Publish("test")
		completed := waitTimeout(&wg, 2*time.Second)
		verify("Channel.Publish(3 subs)",
			completed && count.Load() == 3,
			"received=%d, want 3", count.Load())
	}

	// Publish — 驗證訊息內容正確傳遞
	// Publish — verify message content is correctly delivered
	{
		ch := robin.NewChannel()
		var got atomic.Value
		var wg sync.WaitGroup
		wg.Add(1)
		ch.Subscribe(func(msg string) {
			got.Store(msg)
			wg.Done()
		})
		ch.Publish("hello robin")
		completed := waitTimeout(&wg, 2*time.Second)
		verify("Channel.Publish(content)",
			completed && got.Load() == "hello robin",
			"got=%v, want 'hello robin'", got.Load())
	}

	// Count — 驗證訂閱者計數
	// Count — verify subscriber count
	{
		ch := robin.NewChannel()
		ch.Subscribe(func(string) {})
		ch.Subscribe(func(string) {})
		verify("Channel.Count", ch.Count() == 2,
			"count=%d, want 2", ch.Count())
	}

	// Subscriber.Unsubscribe — 透過訂閱者物件取消訂閱
	// Subscriber.Unsubscribe — unsubscribe via subscriber object
	{
		ch := robin.NewChannel()
		ch.Subscribe(func(string) {})
		sub := ch.Subscribe(func(string) {})
		sub.Unsubscribe()
		verify("Subscriber.Unsubscribe", ch.Count() == 1,
			"count=%d, want 1", ch.Count())
	}

	// Channel.Unsubscribe — 透過頻道方法取消訂閱
	// Channel.Unsubscribe — unsubscribe via channel method
	{
		ch := robin.NewChannel()
		sub := ch.Subscribe(func(string) {})
		ch.Unsubscribe(sub)
		verify("Channel.Unsubscribe", ch.Count() == 0,
			"count=%d, want 0", ch.Count())
	}

	// 取消訂閱後不再收到訊息
	// After unsubscribe, subscriber no longer receives messages
	{
		ch := robin.NewChannel()
		var sub1Count, sub2Count atomic.Int32
		ch.Subscribe(func(msg string) { sub1Count.Add(1) })
		sub2 := ch.Subscribe(func(msg string) { sub2Count.Add(1) })

		ch.Publish("first")
		time.Sleep(200 * time.Millisecond)

		sub2.Unsubscribe()
		ch.Publish("second")
		time.Sleep(200 * time.Millisecond)

		ok := sub1Count.Load() == 2 && sub2Count.Load() == 1
		verify("Unsubscribe(no more msg)", ok,
			"sub1=%d (want 2), sub2=%d (want 1)", sub1Count.Load(), sub2Count.Load())
	}

	// Clear — 移除所有訂閱者
	// Clear — remove all subscribers
	{
		ch := robin.NewChannel()
		ch.Subscribe(func(string) {})
		ch.Subscribe(func(string) {})
		ch.Clear()
		verify("Channel.Clear", ch.Count() == 0,
			"count=%d, want 0", ch.Count())
	}

	// Publish 到無訂閱者的頻道不 panic
	// Publish to empty channel doesn't panic
	{
		ch := robin.NewChannel()
		noPanic := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					noPanic = false
				}
			}()
			ch.Publish("test")
		}()
		time.Sleep(100 * time.Millisecond)
		verify("Channel.Publish(empty)", noPanic, "no panic=%v", noPanic)
	}
}

// ================================================================
// TypedChannel 驗證 / TypedChannel Verification (Generic Pub/Sub)
// ================================================================

// verifyTypedChannel runs a grouped verification flow and reports pass/fail results.
// verifyTypedChannel 會執行一組驗證流程並輸出通過/失敗結果。
func verifyTypedChannel() {
	log.Println("\n--- TypedChannel — 泛型發布/訂閱（無反射）/ Generic Pub/Sub (no reflection) ---")

	// ---- string 型別 / string type ----

	// Subscribe + Publish — 驗證所有訂閱者都收到字串訊息
	// Subscribe + Publish — verify all subscribers receive a string message
	{
		ch := robin.NewTypedChannel[string]()
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(3)
		ch.Subscribe(func(msg string) { count.Add(1); wg.Done() })
		ch.Subscribe(func(msg string) { count.Add(1); wg.Done() })
		ch.Subscribe(func(msg string) { count.Add(1); wg.Done() })
		ch.Publish("test")
		completed := waitTimeout(&wg, 2*time.Second)
		verify("Typed[string].Pub(3 subs)",
			completed && count.Load() == 3,
			"received=%d, want 3", count.Load())
	}

	// Publish — 驗證字串內容正確傳遞
	// Publish — verify string content is correctly delivered
	{
		ch := robin.NewTypedChannel[string]()
		var got atomic.Value
		var wg sync.WaitGroup
		wg.Add(1)
		ch.Subscribe(func(msg string) {
			got.Store(msg)
			wg.Done()
		})
		ch.Publish("hello typed")
		completed := waitTimeout(&wg, 2*time.Second)
		verify("Typed[string].Pub(content)",
			completed && got.Load() == "hello typed",
			"got=%v, want 'hello typed'", got.Load())
	}

	// ---- struct 型別 / struct type ----

	// Subscribe + Publish — 驗證所有訂閱者都收到結構體訊息
	// Subscribe + Publish — verify all subscribers receive a struct message
	{
		ch := robin.NewTypedChannel[gameEvent]()
		var count atomic.Int32
		var wg sync.WaitGroup
		wg.Add(3)
		ch.Subscribe(func(e gameEvent) { count.Add(1); wg.Done() })
		ch.Subscribe(func(e gameEvent) { count.Add(1); wg.Done() })
		ch.Subscribe(func(e gameEvent) { count.Add(1); wg.Done() })
		ch.Publish(gameEvent{Source: "boss", Damage: 999})
		completed := waitTimeout(&wg, 2*time.Second)
		verify("Typed[struct].Pub(3 subs)",
			completed && count.Load() == 3,
			"received=%d, want 3", count.Load())
	}

	// Publish — 驗證結構體欄位完整傳遞
	// Publish — verify struct fields are correctly delivered
	{
		ch := robin.NewTypedChannel[gameEvent]()
		var got atomic.Value
		var wg sync.WaitGroup
		wg.Add(1)
		sent := gameEvent{Source: "archer", Target: "dragon", Damage: 42, Time: time.Now()}
		ch.Subscribe(func(e gameEvent) {
			got.Store(e)
			wg.Done()
		})
		ch.Publish(sent)
		completed := waitTimeout(&wg, 2*time.Second)
		e, _ := got.Load().(gameEvent)
		verify("Typed[struct].Pub(fields)",
			completed && e.Source == sent.Source && e.Target == sent.Target && e.Damage == sent.Damage,
			"got=%+v, want Source=%s Target=%s Damage=%d", e, sent.Source, sent.Target, sent.Damage)
	}

	// ---- Count / Unsubscribe / Clear（混合 string 與 struct）----

	// Count — 驗證 string 訂閱者計數
	// Count — verify string subscriber count
	{
		ch := robin.NewTypedChannel[string]()
		ch.Subscribe(func(string) {})
		ch.Subscribe(func(string) {})
		verify("Typed[string].Count", ch.Count() == 2,
			"count=%d, want 2", ch.Count())
	}

	// Count — 驗證 struct 訂閱者計數
	// Count — verify struct subscriber count
	{
		ch := robin.NewTypedChannel[gameEvent]()
		ch.Subscribe(func(gameEvent) {})
		ch.Subscribe(func(gameEvent) {})
		ch.Subscribe(func(gameEvent) {})
		verify("Typed[struct].Count", ch.Count() == 3,
			"count=%d, want 3", ch.Count())
	}

	// Subscriber.Unsubscribe — 透過訂閱者物件取消訂閱（string）
	// Subscriber.Unsubscribe — unsubscribe via subscriber object (string)
	{
		ch := robin.NewTypedChannel[string]()
		ch.Subscribe(func(string) {})
		sub := ch.Subscribe(func(string) {})
		sub.Unsubscribe()
		verify("Typed[string].Sub.Unsub", ch.Count() == 1,
			"count=%d, want 1", ch.Count())
	}

	// Subscriber.Unsubscribe — 透過訂閱者物件取消訂閱（struct）
	// Subscriber.Unsubscribe — unsubscribe via subscriber object (struct)
	{
		ch := robin.NewTypedChannel[gameEvent]()
		ch.Subscribe(func(gameEvent) {})
		sub := ch.Subscribe(func(gameEvent) {})
		sub.Unsubscribe()
		verify("Typed[struct].Sub.Unsub", ch.Count() == 1,
			"count=%d, want 1", ch.Count())
	}

	// Channel.Unsubscribe — 透過頻道方法取消訂閱（struct）
	// Channel.Unsubscribe — unsubscribe via channel method (struct)
	{
		ch := robin.NewTypedChannel[gameEvent]()
		sub := ch.Subscribe(func(gameEvent) {})
		ch.Unsubscribe(sub)
		verify("Typed[struct].Ch.Unsub", ch.Count() == 0,
			"count=%d, want 0", ch.Count())
	}

	// 取消訂閱後不再收到結構體訊息
	// After unsubscribe, subscriber no longer receives struct messages
	{
		ch := robin.NewTypedChannel[gameEvent]()
		var sub1Count, sub2Count atomic.Int32
		ch.Subscribe(func(e gameEvent) { sub1Count.Add(1) })
		sub2 := ch.Subscribe(func(e gameEvent) { sub2Count.Add(1) })

		ch.Publish(gameEvent{Source: "mage", Damage: 100})
		time.Sleep(200 * time.Millisecond)

		sub2.Unsubscribe()
		ch.Publish(gameEvent{Source: "mage", Damage: 200})
		time.Sleep(200 * time.Millisecond)

		ok := sub1Count.Load() == 2 && sub2Count.Load() == 1
		verify("Typed[struct].Unsub(no msg)", ok,
			"sub1=%d (want 2), sub2=%d (want 1)", sub1Count.Load(), sub2Count.Load())
	}

	// Clear — 移除所有 struct 訂閱者
	// Clear — remove all struct subscribers
	{
		ch := robin.NewTypedChannel[gameEvent]()
		ch.Subscribe(func(gameEvent) {})
		ch.Subscribe(func(gameEvent) {})
		ch.Clear()
		verify("Typed[struct].Clear", ch.Count() == 0,
			"count=%d, want 0", ch.Count())
	}

	// ---- WithFiber（string + struct）----

	// NewTypedChannelWithFiber — 使用自訂 Fiber（string）
	// NewTypedChannelWithFiber — use custom Fiber (string)
	{
		f := robin.NewGoroutineSingle()
		defer f.Dispose()

		ch := robin.NewTypedChannelWithFiber[string](f)
		var got atomic.Value
		var wg sync.WaitGroup
		wg.Add(1)
		ch.Subscribe(func(msg string) {
			got.Store(msg)
			wg.Done()
		})
		ch.Publish("custom fiber")
		completed := waitTimeout(&wg, 2*time.Second)
		verify("WithFiber[string]",
			completed && got.Load() == "custom fiber",
			"got=%v, want 'custom fiber'", got.Load())
	}

	// NewTypedChannelWithFiber — 使用自訂 Fiber（struct）
	// NewTypedChannelWithFiber — use custom Fiber (struct)
	{
		f := robin.NewGoroutineSingle()
		defer f.Dispose()

		ch := robin.NewTypedChannelWithFiber[gameEvent](f)
		var got atomic.Value
		var wg sync.WaitGroup
		wg.Add(1)
		sent := gameEvent{Source: "healer", Target: "ally", Damage: -50}
		ch.Subscribe(func(e gameEvent) {
			got.Store(e)
			wg.Done()
		})
		ch.Publish(sent)
		completed := waitTimeout(&wg, 2*time.Second)
		e, _ := got.Load().(gameEvent)
		verify("WithFiber[struct]",
			completed && e.Source == sent.Source && e.Target == sent.Target && e.Damage == sent.Damage,
			"got=%+v, want Source=%s Target=%s Damage=%d", e, sent.Source, sent.Target, sent.Damage)
	}

	// ---- 邊界情境 / Edge cases ----

	// Publish 到無訂閱者的頻道不 panic（string）
	// Publish to empty channel doesn't panic (string)
	{
		ch := robin.NewTypedChannel[string]()
		noPanic := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					noPanic = false
				}
			}()
			ch.Publish("test")
		}()
		time.Sleep(100 * time.Millisecond)
		verify("Typed[string].Pub(empty)", noPanic, "no panic=%v", noPanic)
	}

	// Publish 到無訂閱者的頻道不 panic（struct）
	// Publish to empty channel doesn't panic (struct)
	{
		ch := robin.NewTypedChannel[gameEvent]()
		noPanic := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					noPanic = false
				}
			}()
			ch.Publish(gameEvent{Source: "ghost", Damage: 0})
		}()
		time.Sleep(100 * time.Millisecond)
		verify("Typed[struct].Pub(empty)", noPanic, "no panic=%v", noPanic)
	}
}

// ================================================================
// 長間隔排程展示（僅示範 API 用法，不驗證時序）
// Long-interval scheduling showcase (API usage demo only, no timing verification)
// ================================================================

// showSchedulingExamples is an internal helper function used in this file.
// showSchedulingExamples 是此檔案內部使用的輔助函式。
func showSchedulingExamples() {
	log.Println("\n--- Scheduling Examples — 排程 API 展示（僅建立，不等待執行）---")
	log.Println("  以下排程已建立但不驗證執行時序（需等待數小時/天/週才會觸發）")
	log.Println("  The following schedules are created but not verified (require hours/days/weeks to trigger)")

	// Every + Minutes — 每 N 分鐘執行
	// Every + Minutes — execute every N minutes
	robin.Every(30).Minutes().Do(func() {
		log.Println("[Showcase] Every 30 Minutes")
	})
	log.Println("  robin.Every(30).Minutes().Do(...)")

	// Every + Hours + At — 每 N 小時在指定 mm:ss 執行
	// Every + Hours + At — execute every N hours at specified mm:ss
	robin.Every(2).Hours().At(0, 30, 0).Do(func() {
		log.Println("[Showcase] Every 2 Hours at :30:00")
	})
	log.Println("  robin.Every(2).Hours().At(0, 30, 0).Do(...)")

	// Every + Days + At — 每 N 天在指定 HH:mm:ss 執行
	// Every + Days + At — execute every N days at specified HH:mm:ss
	robin.Every(1).Days().At(8, 0, 0).Do(func() {
		log.Println("[Showcase] Every 1 Day at 08:00:00")
	})
	log.Println("  robin.Every(1).Days().At(8, 0, 0).Do(...)")

	// Everyday + At — 每天在指定時間執行
	// Everyday + At — execute everyday at specified time
	robin.Everyday().At(23, 59, 59).Do(func() {
		log.Println("[Showcase] Everyday at 23:59:59")
	})
	log.Println("  robin.Everyday().At(23, 59, 59).Do(...)")

	// EveryMonday ~ EverySunday + At — 每週特定日在指定時間執行
	// EveryMonday ~ EverySunday + At — execute weekly on specific day
	robin.EveryMonday().At(9, 0, 0).Do(func() {})
	robin.EveryTuesday().At(9, 0, 0).Do(func() {})
	robin.EveryWednesday().At(9, 0, 0).Do(func() {})
	robin.EveryThursday().At(9, 0, 0).Do(func() {})
	robin.EveryFriday().At(9, 0, 0).Do(func() {})
	robin.EverySaturday().At(10, 0, 0).Do(func() {})
	robin.EverySunday().At(10, 0, 0).Do(func() {})
	log.Println("  robin.EveryMonday().At(9,0,0).Do(...)  ~ robin.EverySunday().At(10,0,0).Do(...)")

	// Between — 限制任務只在指定時間區間內執行
	// Between — restrict execution to a time range
	now := time.Now()
	f := time.Date(0, 1, 1, now.Hour(), 0, 0, 0, time.Local)
	t := time.Date(0, 1, 1, now.Hour(), 59, 59, 0, time.Local)
	robin.Every(30).Seconds().Between(f, t).Times(5).Do(func() {
		log.Println("[Showcase] Between range execution")
	})
	log.Printf("  robin.Every(30).Seconds().Between(%s, %s).Times(5).Do(...)",
		f.Format("15:04:05"), t.Format("15:04:05"))

	// CronScheduler 也支援所有相同的排程方法
	// CronScheduler also supports all the same scheduling methods
	cs := robin.NewCronScheduler()
	cs.EveryMonday().At(8, 0, 0).Do(func() {})
	cs.EveryTuesday().At(8, 0, 0).Do(func() {})
	cs.EveryWednesday().At(8, 0, 0).Do(func() {})
	cs.EveryThursday().At(8, 0, 0).Do(func() {})
	cs.EveryFriday().At(8, 0, 0).Do(func() {})
	cs.EverySaturday().At(9, 0, 0).Do(func() {})
	cs.EverySunday().At(9, 0, 0).Do(func() {})
	log.Println("  cs.EveryMonday().At(8,0,0).Do(...)  ~ cs.EverySunday().At(9,0,0).Do(...)")
}
