package aprocs

import (
	"sync"

	ss "secsplit"
)

func Process(proc Proc, iter ss.ChunkIterator) (err error) {
	errors := make(chan error)
	chunks := make(chan *ss.Chunk)
	done := make(chan struct{})

	go func() {
		defer close(errors)
		wg := sync.WaitGroup{}
		defer wg.Wait()
		for c := range chunks {
			ch := proc.Process(c)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for res := range ch {
					errors <- res.Err
				}
			}()
		}
	}()

	go func() {
		defer close(chunks)
		for iter.Next() {
			c := iter.Chunk()
			select {
			case chunks <- c:
			case <-done:
				return
			}
		}
	}()

	collect := func() (err error) {
		defer close(done)
		for err = range errors {
			if err != nil {
				return
			}
		}
		return
	}

	err = collect()
	for range errors {
	}
	if err != nil {
		return
	}

	return iter.Err()
}
