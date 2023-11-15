package adapters

import "sync"

type SenderCache struct {
	lock  sync.Mutex
	cache map[string]string
}

func (c *SenderCache) GetSenderName(senderID string, fetchGeSanderName func(sendId string) (string, error)) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if name, ok := c.cache[senderID]; ok {
		return name, nil
	}

	name, err := fetchGeSanderName(senderID)
	if err != nil {
		return "", err
	}

	c.cache[senderID] = name
	return name, nil
}

func NewSenderCache() *SenderCache {
	return &SenderCache{
		cache: make(map[string]string),
	}
}
