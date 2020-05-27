package structures

import (
	"container/list"
	"strconv"

	utils "dlocate/osutils"
	log "github.com/sirupsen/logrus"
)

//Cache is a generic struct to LRU cache different paritions data
type Cache struct {
	capacity int
	path     string
	content  map[int]interface{}
	cached   int
	queue    list.List
	freq     map[int]int
}

//GetCache create and return a new cache object using a path and a capacity limit
func GetCache(capacity int, path string) Cache {
	return Cache{capacity: capacity, path: path, content: make(map[int]interface{}), cached: 0, queue: *list.New(), freq: make(map[int]int)}
}

//Get returns a pointer to a parition specific object with an index
func (cache *Cache) Get(index int) *interface{} {
	cache.addIndex(index)
	cache.removeLeastUsed()

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

func (cache *Cache) addIndex(index int) {
	cache.queue.PushBack(index)
	cache.freq[index]++
	if cache.freq[index] == 1 {
		cache.cached++
	}
}

func (cache *Cache) removeLeastUsed() {
	for cache.queue.Len() > 0 && cache.cached > cache.capacity {
		cache.removeQueueFront()
	}
}

func (cache *Cache) removeQueueFront() {
	e := cache.queue.Front()
	cache.queue.Remove(e)

	val := e.Value.(int)

	if cache.freq[val] == 1 {
		delete(cache.freq, val)
		delete(cache.content, val)
		cache.cached--
	} else {
		cache.freq[val]--
	}
}
