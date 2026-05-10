package or

// or объединяет done-каналы в один.
// Возвращаемый канал закрывается, когда закрывается любой входной канал.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		done := make(chan interface{})
		close(done)
		return done
	case 1:
		return channels[0]
	}

	done := make(chan interface{})

	go func() {
		defer close(done)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], done)...):
			}
		}
	}()

	return done
}
