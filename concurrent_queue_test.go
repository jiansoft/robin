package robin

import (
	"reflect"
	"testing"
	"time"
)

func TestConcurrentQueue_Enqueue(t *testing.T) {
	concurrentQueue := NewConcurrentQueue()
	type args struct {
		item interface{}
	}
	tests := []struct {
		name   string
		fields *ConcurrentQueue
		args   []args
	}{
		{"Test_ConcurrentQueue_Enqueue_1", concurrentQueue, []args{{item: "a"}, {item: "b"}, {item: "c"}, {item: "d"}, {item: "e"}}},
		{"Test_ConcurrentQueue_Enqueue_2", concurrentQueue, []args{{item: "1"}, {item: "2"}, {item: "3"}, {item: "4"}, {item: "5"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for ttt := range tt.args {
				tt.fields.Enqueue(ttt)
			}
			t.Logf("%v count:%v", tt.name, tt.fields.Count())
		})
	}
}

func TestConcurrentQueue_TryPeek(t *testing.T) {
	concurrentQueue := NewConcurrentQueue()
	type args struct {
		item interface{}
	}
	tests := []struct {
		name   string
		fields *ConcurrentQueue
		args   args
		want   interface{}
		want1  bool
		want2  int
	}{
		{"Test_ConcurrentQueue_TryPeek_1", concurrentQueue, args{item: "a"}, "a", true, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := tt.fields.TryPeek()
			if got2 {
				t.Errorf("ConcurrentQueue.TryPeek() got = %v, want %v", got1, nil)
			}
			tt.fields.Enqueue(tt.args.item)

			got, got1 := tt.fields.TryPeek()

			if tt.want2 != tt.fields.Count() {
				t.Errorf("ConcurrentQueue.Count() got = %v, want %v", tt.fields.Count(), tt.want2)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConcurrentQueue.TryPeek() got = %v, want %v", got, tt.want)
			}

			if got1 != tt.want1 {
				t.Errorf("ConcurrentQueue.TryPeek() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestConcurrentQueue_TryDequeue(t *testing.T) {
	concurrentQueue := NewConcurrentQueue()
	type args struct {
		item interface{}
	}
	tests := []struct {
		name   string
		fields *ConcurrentQueue
		args   args
		want   interface{}
		want1  bool
		want2  int
	}{
		{"Test_ConcurrentQueue_TryDequeue_1", concurrentQueue, args{item: "a"}, "a", true, 0},
		{"Test_ConcurrentQueue_TryDequeue_2", concurrentQueue, args{item: "b"}, "b", true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Enqueue(tt.args.item)
			got, got1 := tt.fields.TryDequeue()

			if tt.want2 != tt.fields.Count() {
				t.Errorf("ConcurrentQueue.Count() got = %v, want %v", tt.fields.Count(), tt.want2)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConcurrentQueue.TryDequeue() got = %v, want %v", got, tt.want)
			}

			if got1 != tt.want1 {
				t.Errorf("ConcurrentQueue.TryDequeue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestConcurrentQueue_Random(t *testing.T) {
	concurrentQueue := NewConcurrentQueue()
	type args struct {
		item interface{}
	}
	tests := []struct {
		name   string
		fields *ConcurrentQueue
		args   args
	}{
		{"Test_ConcurrentQueue_Random_1", concurrentQueue, args{item: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				RightNow().Do(func() {
					tt.fields.Enqueue(tt.args.item)
				})
				RightNow().Do(func() {
					_, _ = tt.fields.TryDequeue()
				})
				RightNow().Do(func() {
					_ = tt.fields.Count()
				})
			}
			timeout := time.NewTimer(time.Duration(1000) * time.Millisecond)

			select {
			case <-timeout.C:
				t.Logf("concurrentQueue now has count:%v", tt.fields.Count())
			}
		})
	}
}
