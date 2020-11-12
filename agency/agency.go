package agency

type Listener interface {
	AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string)
	AddMessage(connectionID, id, message string, sentByMe bool)
	UpdateMessage(connectionID, id, delivered bool)
}

type Agency interface {
	Init(l Listener)
	Invite() (string, string, error)
	Connect(invitation string) (string, error)
	SendMessage(connectionID, message string) (string, error)
}
