package concurrency

import (
	"errors"
	"sync"
)

func ForEach[Item any](items []Item, concurrency int, function func(item Item) error) error {
	var (
		mutex = sync.Mutex{}
		wg    = sync.WaitGroup{}
		errs  = make([]error, 0)
		ch    = make(chan struct{}, concurrency)
	)

	wg.Add(len(items))
	for _, item := range items {
		go func(item Item, wg *sync.WaitGroup) {
			ch <- struct{}{}
			defer wg.Done()

			err := function(item)
			if err == nil {
				<-ch
				return
			}

			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
			<-ch
		}(item, &wg)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
