package gb

// A Target is a placeholder for work which is completed asyncronusly.
type Target interface {
	// Result returns the result of the work as an error, or nil if the work
	// was performed successfully.
	// Implementers must observe these invariants
	// 1. There may be multiple concurrent callers to Result, or Result may
	//    be called many times in sequence, it must always return the same
	// 2. Result blocks until the work has been performed.
	Result() error
}

type target struct {
	c chan error
}

func newTarget(f func() error, deps ...Target) target {
	c := make(chan error, 1)
	go func() {
		for _, dep := range deps {
			if err := dep.Result(); err != nil {
				c <- err
				return
			}
		}
		c <- f()
	}()
	return target{c: c}
}

func (t *target) Result() error {
	err := <-t.c
	t.c <- err
	return err
}
