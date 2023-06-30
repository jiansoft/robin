package robin

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func equal(t *testing.T, got, want any) {
	if !reflect.DeepEqual(got, want) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[37m%s:%d:\n got: %#v\nwant: %#v\033[39m\n ", filepath.Base(file), line, got, want)
		t.FailNow()
	}
}

func TestConcurrent(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name   string
		fields []any
		args   []args
	}{
		{"Test_Concurrent", []any{
			NewConcurrentQueue(),
			NewConcurrentStack(),
			NewConcurrentBag()},
			[]args{
				{item: "a"},
				{item: "b"},
				{item: "c"},
				{item: "d"},
				{item: "e"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, collection := range tt.fields {
				switch collection.(type) {
				case *ConcurrentQueue:
					c := collection.(*ConcurrentQueue)
					for _, item := range tt.args {
						c.Enqueue(item.item)
					}

					total := c.ToArray()
					equal(t, len(tt.args), len(total))
					for i := 0; i < len(total); i++ {
						equal(t, tt.args[i].item, total[i])
					}

					if v, ok := c.TryPeek(); ok {
						equal(t, tt.args[0].item, v.(string))
					}

					for i := 0; ; i++ {
						if v, ok := c.TryDequeue(); !ok {
							break
						} else {
							//t.Logf("first in-first out (FIFO) v:%s", v.(string))
							equal(t, tt.args[i].item, v.(string))
						}
					}

					for i := 0; i < 200; i++ {
						c.Enqueue(fmt.Sprintf("item%d", i))
					}

					c.Clear()
					// 移除大量元素，导致容器缩小
					for i := 0; i < 210; i++ {
						c.TryPeek()
						c.TryDequeue()
					}
					equal(t, 0, c.Len())

					break
				case *ConcurrentStack:
					c := collection.(*ConcurrentStack)
					for _, item := range tt.args {
						c.Push(item.item)
					}

					total := c.ToArray()
					equal(t, len(tt.args), len(total))

					for i := len(total); i > len(total); i-- {
						equal(t, tt.args[i].item, total[i])
					}

					if v, ok := c.TryPeek(); ok {
						equal(t, tt.args[len(tt.args)-1].item, v.(string))
					}

					for i := c.Len() - 1; ; i-- {
						if v, ok := c.TryPop(); !ok {
							break
						} else {
							//t.Logf("last in-first out (LIFO) v:%s", v.(string))
							equal(t, tt.args[i].item, v.(string))
						}
					}

					for i := 0; i < 200; i++ {
						c.Push(fmt.Sprintf("item%d", i))
					}

					c.Clear()
					// 移除大量元素，导致容器缩小
					for i := 0; i < 210; i++ {
						c.TryPeek()
						c.TryPop()
					}


					equal(t, 0, c.Len())
					break
				case *ConcurrentBag:
					c := collection.(*ConcurrentBag)
					for _, item := range tt.args {
						c.Add(item.item)
					}

					total := c.ToArray()
					assert.Equal(t, len(tt.args), len(total), "Expected len is eq")

					for i := 0; i < len(total); i++ {
						if v, ok := c.TryTake(); !ok {
							break
						} else {
							//t.Logf("unordered collection v:%s remain:%v", v.(string) , c.Len())
							equal(t, tt.args[i].item, v.(string))
							assert.Equal(t, tt.args[i].item, v.(string), "Expected items are eq")
						}
					}


					for i := 0; i < 200; i++ {
						c.Add(fmt.Sprintf("item%d", i))
					}

					// 移除大量元素，导致容器缩小
					for i := 0; i < 210; i++ {
						c.TryTake()
					}

					c.Clear()
					assert.Equal(t, 0,c.Len(), "Expected len is eq")
					break
				}
			}
		})
	}
}
