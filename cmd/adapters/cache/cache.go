package cache

type Cache[T comparable, R any] struct {
	data map[T]R
}

func NewCache[T comparable, R any]() *Cache[T, R] {
	return &Cache[T, R]{
		data: make(map[T]R),
	}
}

func (c *Cache[T, R]) GetCache(key T) (R, bool) {
	data, ok := c.data[key]
	return data, ok
}

func (c *Cache[T, R]) SetCache(key T, value R) {
	c.data[key] = value
}
