package structures

//Cache is a generic struct to LRU cache different paritions data
type Cache struct {
	capacity      int
	content       map[string]interface{}
	cached        int
	requestIndex  int
	lastRequested map[string]int
}

//GetCache create and return a new cache object using a path and a capacity limit
func GetCache(capacity int) Cache {
	return Cache{capacity: capacity, content: make(map[string]interface{}), cached: 0, requestIndex: 1, lastRequested: make(map[string]int)}
}

//Get returns a pointer to a parition specific object with an index
func (cache *Cache) Get(key string) (interface{}, bool) {
	//cache hit
	if val, ok := cache.content[key]; ok {
		cache.addIndex(key)
		return val, true
	}

	//cache miss
	return nil, false
}

//Clear remove all objects from cache
func (cache *Cache) Clear() {
	cache.cached = 0
	cache.content = make(map[string]interface{})
	cache.requestIndex = 1
	cache.lastRequested = make(map[string]int)
}

//Set save a specific value into an object
func (cache *Cache) Set(key string, object interface{}) {
	cache.content[key] = object

	cache.addIndex(key)
	if cache.cached > cache.capacity {
		cache.removeLeastUsed()
	}
}

//Delete removes a cached object from the cache
func (cache *Cache) Delete(key string) {
	delete(cache.lastRequested, key)
	delete(cache.content, key)
	cache.cached = len(cache.lastRequested)
}

func (cache *Cache) addIndex(key string) {
	cache.lastRequested[key] = cache.requestIndex
	cache.requestIndex++
	cache.cached = len(cache.lastRequested)
}

func (cache *Cache) removeLeastUsed() {
	var minValue = 1000000000
	var minIndex = ""

	for k, v := range cache.lastRequested {
		if v < minValue {
			minValue = v
			minIndex = k
		}
	}
	delete(cache.lastRequested, minIndex)
	delete(cache.content, minIndex)
	cache.cached = len(cache.lastRequested)
}
