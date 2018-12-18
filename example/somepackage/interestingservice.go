package somepackage

import (
	"fmt"
	"io/ioutil"
	"time"
)

type InterestingService struct {
	count    int
	ID       int
	interval time.Duration
	end      bool
}

func NewInterestingService(ID int, interval time.Duration, path string) *InterestingService {
	it := &InterestingService{ID: ID,
		interval: interval,
	}
	go func() {
		for !it.end {
			message := fmt.Sprintf("Interesting Service #%d running!. Count=%d\n", it.ID, it.count)
			fmt.Println(message)
			ioutil.WriteFile(path, []byte(message), 0666)
			time.Sleep(it.interval)
		}

	}()
	return it
}

func (it *InterestingService) Close() error {
	fmt.Printf("Terminating Interesting Service %d...\n", it.ID)
	it.end = true
	return nil
}
