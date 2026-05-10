package or

import (
	"testing"
	"time"
)

func TestOrWithoutChannelsReturnsClosedChannel(t *testing.T) {
	done := or()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("expected returned channel to be closed")
	}
}

func TestOrWithOneChannelReturnsAfterInputClosed(t *testing.T) {
	ch := make(chan interface{})
	done := or(ch)

	select {
	case <-done:
		t.Fatal("done closed before input channel")
	default:
	}

	close(ch)

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("expected done to close after input channel")
	}
}

func TestOrClosesWhenAnyChannelCloses(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	done := or(ch1, ch2, ch3)

	select {
	case <-done:
		t.Fatal("done closed before any input channel")
	default:
	}

	close(ch2)

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("expected done to close after one input channel closes")
	}
}

func TestOrDoesNotWaitForAllChannels(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	done := or(ch1, ch2, ch3)

	close(ch1)

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("expected done to close after first closed channel")
	}
}

func TestOrWorksWithManyChannels(t *testing.T) {
	const channelsCount = 100

	channels := make([]<-chan interface{}, 0, channelsCount)
	closers := make([]chan interface{}, 0, channelsCount)

	for i := 0; i < channelsCount; i++ {
		ch := make(chan interface{})
		closers = append(closers, ch)
		channels = append(channels, ch)
	}

	done := or(channels...)

	close(closers[73])

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected done to close when one of many channels closes")
	}
}

func TestOrExampleTiming(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		ch := make(chan interface{})

		go func() {
			defer close(ch)
			time.Sleep(after)
		}()

		return ch
	}

	start := time.Now()

	<-or(
		sig(200*time.Millisecond),
		sig(100*time.Millisecond),
		sig(10*time.Millisecond),
		sig(150*time.Millisecond),
	)

	elapsed := time.Since(start)

	if elapsed > 80*time.Millisecond {
		t.Fatalf("or waited too long: %v", elapsed)
	}
}