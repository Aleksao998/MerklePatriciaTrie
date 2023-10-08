package mockStorage

type (
	hasDelegate    func(key []byte) (bool, error)
	getDelegate    func(key []byte) ([]byte, error)
	putDelegate    func(key []byte, value []byte) error
	deleteDelegate func(key []byte) error
)

type MockStorage struct {
	HasFn    hasDelegate
	GetFn    getDelegate
	PutFn    putDelegate
	DeleteFn deleteDelegate
}

func (m *MockStorage) Has(key []byte) (bool, error) {
	if m.HasFn != nil {
		return m.HasFn(key)
	}
	return false, nil
}

func (m *MockStorage) Get(key []byte) ([]byte, error) {
	if m.GetFn != nil {
		return m.GetFn(key)
	}
	return nil, nil
}

func (m *MockStorage) Put(key []byte, value []byte) error {
	if m.PutFn != nil {
		return m.PutFn(key, value)
	}
	return nil
}

func (m *MockStorage) Delete(key []byte) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(key)
	}
	return nil
}
