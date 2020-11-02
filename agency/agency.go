package agency

type Listener interface {
	AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string, initiatedByUs bool)
}

type Agency interface {
	Init(l Listener)
	Invite() (string, error)
	Connect(invitation string) (string, error)
}
