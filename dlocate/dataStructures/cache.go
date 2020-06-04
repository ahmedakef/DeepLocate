package structures

import (
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

//Cache is a generic struct to LRU cache different paritions data
type Cache struct {
	capacity      int
	path          string
	content       map[string]interface{}
	cached        int
	requestIndex  int
	lastRequested map[string]int
}

//GetCache create and return a new cache object using a path and a capacity limit
func GetCache(capacity int, path string) Cache {
	return Cache{capacity: capacity, path: path, content: make(map[string]interface{}), cached: 0, requestIndex: 1, lastRequested: make(map[string]int)}
}

//Get returns a pointer to a parition specific object with an index
func (cache *Cache) Get(key string) (interface{}, error) {
	//cache hit
	if val, ok := cache.content[key]; ok {
		cache.addIndex(key)
		return &val, nil
	}

	//cache miss
	path := cache.path + key + ".gob"

	var object interface{}
	err := utils.ReadGob(path, object)
	if err != nil {
		log.Errorf("Error while reading object for partition %q: %v\n", key, err)
		return nil, err
	}

	cache.addIndex(key)
	if cache.cached > cache.capacity {
		cache.removeLeastUsed()
	}

	return &object, err
}

//Clear stores all changes to objects in files and remove them from cache
func (cache *Cache) Clear() {
	for cache.cached > 0 {
		cache.removeLeastUsed()
	}
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
	utils.SaveGob(cache.content[minIndex], cache.path+minIndex+".gob")
	delete(cache.lastRequested, minIndex)
	delete(cache.content, minIndex)
}
