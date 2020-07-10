package cats

import (
	"context"
)

type Cat string

func NewCat(name string) *Cat {
	cat := Cat(name)
	return &cat
}

//go:generate mockery -name=GetCatByID
type GetCatByID interface {
	ByID(id int) (*Cat, error)
}

type CatsAPI struct {
	Datastore GetCatByID
}

func (c CatsAPI) Fetch(ctx context.Context) (chan *Cat, chan error) {
	outCh := make(chan *Cat)
	errCh := make(chan error)
	go func() {
		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				close(outCh)
				return
			default:
			}

			cat, err := c.Datastore.ByID(i)

			if err != nil {
				errCh <- err
				continue
			}
			if cat == nil {
				close(outCh)
				return
			}
			outCh <- cat
		}
	}()
	return outCh, errCh
}
