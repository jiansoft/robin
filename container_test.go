package robin

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	type fields struct {
		c *container
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"TestContainer", fields{&container{}}},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loop := 100
			for i := 0; i < loop; i++ {
				wg.Add(1)
				RightNow().Do(func() {
					for i := 0; i < loop; i++ {
						pending1 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
						pending2 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
						pending3 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
						tt.fields.c.Add(pending1)
						tt.fields.c.Add(pending2)
						tt.fields.c.Add(pending3)
						tt.fields.c.Remove(pending2)
						tt.fields.c.Remove(pending3)
						tt.fields.c.Remove(pending3)
					}
					wg.Done()
				})

				pending1 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
				pending2 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
				pending3 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
				tt.fields.c.Add(pending1)
				tt.fields.c.Add(pending2)
				tt.fields.c.Add(pending3)
				tt.fields.c.Remove(pending1)
				tt.fields.c.Remove(pending2)
				tt.fields.c.Remove(pending3)

				wg.Wait()

				ss := tt.fields.c.Items()
				assert.Equal(t, len(ss), tt.fields.c.Count(), "they should be equal %+v", ss)
				assert.Equal(t, len(ss), loop, "they should be equal %+v", ss)
				assert.Equal(t, tt.fields.c.Count(), loop, "they should be equal %+v", ss)
				tt.fields.c.Dispose()
			}
		})
	}
}
