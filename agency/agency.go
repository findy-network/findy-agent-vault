package agency

type Listener interface {
	AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string)
	AddMessage(connectionId, id, message string, sentByMe bool)
	UpdateMessage(connectionId, id, delivered bool)
}

type Agency interface {
	Init(l Listener)
	Invite() (string, string, error)
	Connect(invitation string) (string, error)
	SendMessage(connectionId, message string) (string, error)
}
