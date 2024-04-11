package utility

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFixedLengthUUID(t *testing.T) {
	InitSnowflakeNode(100)
	cases := make([]uint8, 0)
	var i uint8
	for i = 0; i < 10; i++ {
		cases = append(cases, 8+i)
	}

	for _, length := range cases {
		uuid := GenerateFixedLengthUUID(length)
		assert.Equal(t, length, uint8(len(uuid)))
		fmt.Printf("GenerateFixedLengthUUID length=%d, uuid=%s\n", length, uuid)
	}
}

func TestGenerateFixedLengthUUIDUnique(t *testing.T) {
	var m sync.Map
	var g sync.WaitGroup
	goroutineCount := 1000
	generateTimes := 100
	InitSnowflakeNode(100)

	for i := 0; i < goroutineCount; i++ {
		g.Add(1)
		go func() {
			for i := 0; i < generateTimes; i++ {
				time.Sleep(time.Millisecond * 1)
				uuid := GenerateFixedLengthUUID(14)
				_, ok := m.Load(uuid)
				if ok {
					panic(fmt.Errorf("uuid %s is found\n", uuid))
				} else {
					m.Store(uuid, uuid)
				}
			}
			g.Done()
		}()
	}
	g.Wait()
}
func TestSnowflake(t *testing.T) {
	InitSnowflakeNode(100)
	id := GenerateSnowflakeID()
	idTime := TimeInSnowflakeID(id)
	fmt.Println(id, idTime, id.Base36(), id.Base64())
}

func TestGenUUID(t *testing.T) {
	id := GenUUID()
	fmt.Println(id)
}
