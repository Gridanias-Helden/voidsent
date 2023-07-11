package higherorlower_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/gridanias-helden/voidsent/pkg/services"
	"github.com/gridanias-helden/voidsent/pkg/services/games/higherorlower"
)

var broker = services.NewBroker()

type Test struct {
	value   int
	tries   int
	history []int
	index   int
	t       *testing.T

	wg    *sync.WaitGroup
	games int
}

func NewTest(t *testing.T, wg *sync.WaitGroup) services.Service {
	return &Test{
		value:   49,
		tries:   7,
		history: []int{50, 25, 37, 43, 46, 48, 49},
		t:       t,
		wg:      wg,
	}
}

func (g *Test) Send(from string, to string, topic string, body any) {
	switch vs := body.(type) {
	case map[string]int:
		if v, ok := vs["value"]; ok {
			if v > g.value {
				broker.Send(to, from, topic, "lower")
			} else if v < g.value {
				broker.Send(to, from, topic, "higher")
			} else {
				g.games++
				if g.games < 1000 {
					g.value = 49
					g.tries = 7
					g.history = []int{50, 25, 37, 43, 46, 48, 49}
					broker.Send(to, from, topic, "start")
				} else {
					broker.Send(to, from, topic, "exit")
				}
			}
		}

	case string:
		broker.Send(to, "count", "", "done")
		g.wg.Done()

	default:
		g.t.Errorf("Unexpected type %t", vs)
	}
}

type Count struct {
	c int
	m sync.RWMutex
}

func (c *Count) Send(_ string, _ string, _ string, _ any) {
	c.m.Lock()
	defer c.m.Unlock()

	c.c++
}

func TestHigherOrLowerGameplay(t *testing.T) {
	const count = 1000
	wg := &sync.WaitGroup{}
	wg.Add(count)
	counter := &Count{c: 0}
	for index := 0; index < count; index++ {
		broker.AddService("count", counter)
		broker.AddService("hol"+strconv.Itoa(index), higherorlower.New(broker))
		broker.AddService("test"+strconv.Itoa(index), NewTest(t, wg))
	}

	for index := 0; index < count; index++ {
		broker.Send("test"+strconv.Itoa(index), "hol"+strconv.Itoa(index), "", "start")
	}

	wg.Wait()
}
