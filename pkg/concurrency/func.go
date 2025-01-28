package concurrency

import (
	"errors"
	"sync"
)

func ForEach[Item any, Response any](items []Item, concurrency int, function func(item Item) (Response, error)) ([]Response, error) {
	var (
		mutex     = sync.Mutex{}
		wg        = sync.WaitGroup{}
		errs      = make([]error, 0)
		ch        = make(chan struct{}, concurrency)
		responses = make([]Response, 0, 0)
	)

	wg.Add(len(items))
	for _, item := range items {
		go func(item Item, wg *sync.WaitGroup) {
			ch <- struct{}{}
			defer wg.Done()

			response, err := function(item)

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				errs = append(errs, err)
				<-ch
				return
			}

			responses = append(responses, response)
			<-ch
		}(item, &wg)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return responses, nil
}
