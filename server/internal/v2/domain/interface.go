package domain

type Subscription interface {
	Err() <-chan error
	Unsubscribe()
}
