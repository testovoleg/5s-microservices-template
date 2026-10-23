package utils

import (
	"sync"
	"time"
)

// TTLCache — кэш процесса «ключ → значение с TTL» для редко меняющихся справочников
// (настройки компании, каталоги, правила), которые дорого читать из БД/Admin API на
// каждый запрос. Переполнение (max записей) сбрасывает карту целиком — грубая, но дешёвая
// защита от роста памяти на редких ключах. Загрузка при промахе идёт без singleflight:
// два одновременных промаха по одному ключу прочитают источник дважды и запишут одно и
// то же. Инвалидация — Delete/Reset после записи в источник.
type TTLCache[V any] struct {
	mu    sync.Mutex
	ttl   time.Duration
	max   int
	items map[string]ttlEntry[V]
}

type ttlEntry[V any] struct {
	loaded time.Time
	value  V
}

func NewTTLCache[V any](ttl time.Duration, max int) *TTLCache[V] {
	return &TTLCache[V]{ttl: ttl, max: max, items: map[string]ttlEntry[V]{}}
}

// Get — значение по ключу, если оно есть и не протухло.
func (c *TTLCache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok || time.Since(e.loaded) >= c.ttl {
		var zero V
		return zero, false
	}
	return e.value, true
}

func (c *TTLCache[V]) Put(key string, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.max > 0 && len(c.items) >= c.max {
		c.items = map[string]ttlEntry[V]{}
	}
	c.items[key] = ttlEntry[V]{loaded: time.Now(), value: v}
}

// GetOrLoad — значение из кэша либо от load с записью в кэш; ошибка load не кэшируется.
func (c *TTLCache[V]) GetOrLoad(key string, load func() (V, error)) (V, error) {
	if v, ok := c.Get(key); ok {
		return v, nil
	}
	v, err := load()
	if err != nil {
		return v, err
	}
	c.Put(key, v)
	return v, nil
}

func (c *TTLCache[V]) Delete(keys ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *TTLCache[V]) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]ttlEntry[V]{}
}
