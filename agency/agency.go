package agency

type Agency interface {
	Init()
	Invite() (string, error)
	Connect() (string, error)
}
