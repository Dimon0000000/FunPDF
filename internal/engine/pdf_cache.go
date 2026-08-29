package engine

import (
	"sync"
	"time"
)

type PDFTextCacheItem struct {
	Text     string    `json:"text"`
	ExpireAt time.Time `json:"expire_at"`
}

type PDFTextCache struct {
	mu    sync.RWMutex
	items map[string]PDFTextCacheItem
	ttl   time.Duration
}

var PDFText = NewPDFTextCache(time.Hour * 2)

func NewPDFTextCache(ttl time.Duration) *PDFTextCache {
	return &PDFTextCache{
		items: make(map[string]PDFTextCacheItem),
		ttl:   ttl,
	}
}

func (c *PDFTextCache) Get(fileID string) (string, bool) {
	c.mu.RLock()
	item, ok := c.items[fileID]
	c.mu.RUnlock()

	if !ok {
		return "", false
	}
	if !item.ExpireAt.IsZero() && time.Now().After(item.ExpireAt) {
		c.Delete(fileID)
		return "", false
	}

	return item.Text, true
}

func (c *PDFTextCache) Set(fileID, text string) {
	c.mu.Lock()
	c.items[fileID] = PDFTextCacheItem{
		Text:     text,
		ExpireAt: time.Now().Add(c.ttl),
	}
	defer c.mu.Unlock()
}

func (c *PDFTextCache) Delete(fileID string) bool {
	c.mu.Lock()
	_, ok := c.items[fileID]
	delete(c.items, fileID)
	c.mu.Unlock()
	return ok
}
