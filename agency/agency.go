package agency

type Listener interface {
	AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string)
}

type Agency interface {
	Init(l Listener)
	Invite() (string, string, error)
	Connect(invitation string) (string, error)
}
