// +build !findy

package agency

type Mock struct{}

var Instance Agency = &Mock{}

func (m *Mock) Init() {}

func (m *Mock) Invite() (invitation string, err error) {
	return "", nil
}

func (m *Mock) Connect() (string, error) {
	return "", nil
}
