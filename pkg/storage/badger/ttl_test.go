package badger

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTL(t *testing.T) {
	dir := "./temp-ttl"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		require.Equal(t, nil, err)
	}
	s := NewStore(DefaultOptions(dir))

	//for i := 0; true; i++ {
	//	a, s, err := s.GetProfile(fmt.Sprintf("%d", i))
	//	fmt.Sprintln(a, s, err)
	//}

	for i := 0; true; i++ {
		res, err := os.ReadFile("./trace_119091.gz")
		if err != nil {
			panic(err)
		}
		s, err := s.SaveProfile(fmt.Sprintf("%d", i), res, 10*time.Second)
		if err != nil {
			fmt.Println(s, err)
		}
	}
}
