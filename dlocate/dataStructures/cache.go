package structures

import (
	"strconv"

	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

//Cache is a generic struct to LRU cache different paritions data
type Cache struct {
	capacity      int
	path          string
	content       map[int]interface{}
	cached        int
	requestIndex  int
	lastRequested map[int]int
}

//GetCache create and return a new cache object using a path and a capacity limit
func GetCache(capacity int, path string) Cache {
	return Cache{capacity: capacity, path: path, content: make(map[int]interface{}), cached: 0, requestIndex: 1, lastRequested: make(map[int]int)}
}

//Get returns a pointer to a parition specific object with an index
func (cache *Cache) Get(index int) *interface{} {
	cache.addIndex(index)
	if cache.cached > cache.capacity {
		cache.removeLeastUsed()
	}

	//cache hit
	if val, ok := cache.content[index]; ok {
		return &val
	}

	//cache miss
	path := cache.path + strconv.Itoa(index) + ".gob"

	var object interface{}
	err := utils.ReadGob(path, object)
	if err != nil {
		log.Errorf("Error while reading object for partition %q: %v\n", index, err)
	}
	return &object
}

//Delete removes a cached object from the cache
func (cache *Cache) Delete(index int) {
	delete(cache.lastRequested, index)
	delete(cache.content, index)
	cache.cached = len(cache.lastRequested)
}

func (cache *Cache) addIndex(index int) {
	cache.lastRequested[index] = cache.requestIndex
	cache.requestIndex++
	cache.cached = len(cache.lastRequested)
}

func (cache *Cache) removeLeastUsed() {
	var minValue = 0
	var minIndex = -1

	for k, v := range cache.lastRequested {
		if v < minValue {
			minValue = v
			minIndex = k
		}
	}
	delete(cache.lastRequested, minIndex)
	delete(cache.content, minIndex)
}
