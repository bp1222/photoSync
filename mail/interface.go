package mail

type Sender interface {
	Send(from, frameId, image string) error
}
