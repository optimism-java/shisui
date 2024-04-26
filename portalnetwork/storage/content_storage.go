package storage

import "fmt"

var ErrContentNotFound = fmt.Errorf("content not found")

type ContentType byte

type ContentKey struct {
	selector ContentType
	data     []byte
}

func NewContentKey(selector ContentType, data []byte) *ContentKey {
	return &ContentKey{
		selector: selector,
		data:     data,
	}
}

func (c *ContentKey) Encode() []byte {
	res := make([]byte, 0, len(c.data)+1)
	res = append(res, byte(c.selector))
	res = append(res, c.data...)
	return res
}

type ContentStorage interface {
	Get(contentKey []byte, contentId []byte) ([]byte, error)

	Put(contentKey []byte, contentId []byte, content []byte) error
}

type MockStorage struct {
	Db map[string][]byte
}

func (m *MockStorage) Get(contentKey []byte, contentId []byte) ([]byte, error) {
	if content, ok := m.Db[string(contentId)]; ok {
		return content, nil
	}
	return nil, ErrContentNotFound
}

func (m *MockStorage) Put(contentKey []byte, contentId []byte, content []byte) error {
	m.Db[string(contentId)] = content
	return nil
}