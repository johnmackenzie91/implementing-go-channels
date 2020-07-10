package cats_test

import (
	"no_vcs/me/channels/internal/cats"
	"no_vcs/me/channels/internal/cats/mocks"
	"testing"
	"time"

	"context"

	"github.com/stretchr/testify/assert"
)

func TestFetch_ChannelClosedWhenNoMoreRecordsToSend(t *testing.T) {
	// init the mock data source
	getCatMock := mocks.GetCatByID{}
	getCatMock.On("ByID", 1).Return(cats.NewCat("Flossy"), nil)
	getCatMock.On("ByID", 2).Return(cats.NewCat("Mildred"), nil)
	getCatMock.On("ByID", 3).Return(cats.NewCat("Meowth"), nil)
	// cat with id of 4 does not exist, signalling no more records left
	getCatMock.On("ByID", 4).Return(nil, nil)

	ctx := context.TODO()

	outCh, errCh := (cats.CatsAPI{
		Datastore: &getCatMock,
	}).Fetch(ctx)

	// check for errors, this test does not expect any
	go func() {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}()

	output := []*cats.Cat{}

	var closeLoop bool
	for {
		select {
		case <-time.After(2 * time.Second):
			t.Fatal("assumed that test is stuck in infinite loop")
		case v, ok := <-outCh:
			if !ok {
				closeLoop = true
				break
			}
			output = append(output, v)
		}
		if closeLoop {
			break
		}
	}

	expectedOutput := []*cats.Cat{
		cats.NewCat("Flossy"),
		cats.NewCat("Mildred"),
		cats.NewCat("Meowth"),
	}
	assert.Equal(t, expectedOutput, output)

	getCatMock.AssertExpectations(t)
}

func TestFetch_ChannelClosedWhenContextIsCancelled(t *testing.T) {
	// init the mock data source
	// we do not expect any calls to this method, the context should be closed already be closed
	getCatMock := mocks.GetCatByID{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	outCh, errCh := (cats.CatsAPI{
		Datastore: &getCatMock,
	}).Fetch(ctx)

	// check for errors, this test does not expect any
	go func() {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}()

	output := []*cats.Cat{}

	var closeLoop bool
	for {
		select {
		case <-time.After(2 * time.Second):
			t.Fatal("assumed that test is stuck in infinite loop")
		case v, ok := <-outCh:
			if !ok {
				closeLoop = true
				break
			}
			output = append(output, v)
		}
		if closeLoop {
			break
		}
	}

	assert.Empty(t, output)
	getCatMock.AssertExpectations(t)
}
