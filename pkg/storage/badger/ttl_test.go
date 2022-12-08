package badger

import (
	"testing"
)

func TestTTL(t *testing.T) {
	//dir := "./temp-ttl"
	//if _, err := os.Stat(dir); os.IsNotExist(err) {
	//	err := os.Mkdir(dir, 0777)
	//	require.Equal(t, nil, err)
	//}
	//s := NewStore(DefaultOptions(dir).WithGCInternal(time.Second * 10))
	//for i := 0; true; i++ {
	//	res, err := os.ReadFile("./trace_119091.gz")
	//	if err != nil {
	//		panic(err)
	//	}
	//	s, err := s.SaveProfile(fmt.Sprintf("%d", i), res, 1*time.Minute)
	//	if err != nil {
	//		fmt.Println(s, err)
	//	}
	//	if i%1000 == 0 {
	//		time.Sleep(time.Second * 10)
	//	}
	//}
}
