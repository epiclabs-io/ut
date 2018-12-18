// Copyright 2018 The ut/microtest Authors
// This file is part of ut/microtest library.
//
// This library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this library. If not, see <http://www.gnu.org/licenses/>.

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
