package robin

// Enqueue0 enqueues a zero-argument function on the fiber with compile-time type safety.
func Enqueue0(f Fiber, fn func()) {
	f.enqueueTask(task{fn: fn})
}

// Enqueue1 enqueues a one-argument function on the fiber with compile-time type safety.
func Enqueue1[T1 any](f Fiber, fn func(T1), a1 T1) {
	f.Enqueue(fn, a1)
}

// Enqueue2 enqueues a two-argument function on the fiber with compile-time type safety.
func Enqueue2[T1, T2 any](f Fiber, fn func(T1, T2), a1 T1, a2 T2) {
	f.Enqueue(fn, a1, a2)
}

// Enqueue3 enqueues a three-argument function on the fiber with compile-time type safety.
func Enqueue3[T1, T2, T3 any](f Fiber, fn func(T1, T2, T3), a1 T1, a2 T2, a3 T3) {
	f.Enqueue(fn, a1, a2, a3)
}
