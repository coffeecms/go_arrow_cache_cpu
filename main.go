package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
)

// CacheItem đại diện cho một mục trong cache
type CacheItem struct {
	Value     []byte
	ExpiresAt time.Time
}

// Cache là cấu trúc lưu trữ các mục cache với Apache Arrow
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	pool  *memory.GoAllocator
}

// NewCache khởi tạo cache mới
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
		pool:  memory.NewGoAllocator(),
	}
}

// SetWithArrow sử dụng Apache Arrow để lưu trữ dữ liệu
func (c *Cache) SetWithArrow(key string, value []byte, ttl time.Duration) {
	// Sử dụng Arrow để xây dựng bộ nhớ nhị phân
	builder := array.NewBinaryBuilder(c.pool, arrow.BinaryTypes.Binary)
	defer builder.Release()

	// Thêm giá trị vào builder Arrow
	builder.Append(value)
	arr := builder.NewArray()
	defer arr.Release()

	// Lưu trữ vào cache
	c.mu.Lock()
	c.items[key] = CacheItem{
		Value:     arr.(*array.Binary).Value(0),
		ExpiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}

// Get lấy dữ liệu từ cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return item.Value, true
}

// CleanExpiredItems dọn dẹp các mục đã hết hạn
func (c *Cache) CleanExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if time.Now().After(v.ExpiresAt) {
			delete(c.items, k)
		}
	}
}

// BenchmarkSet đo tổng thời gian thực hiện Set với Apache Arrow cho 1,000,000 kết nối
func BenchmarkSet(cache *Cache, wg *sync.WaitGroup, numConnections int) {
	startTime := time.Now()

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.SetWithArrow(fmt.Sprintf("key%d", i), []byte("some data"), 1*time.Minute)
		}(i)
	}

	wg.Wait()

	endTime := time.Now()
	fmt.Printf("Thời gian Set với Apache Arrow cho %d kết nối: %v\n", numConnections, endTime.Sub(startTime))
}

// BenchmarkGet đo tổng thời gian thực hiện Get cho 1,000,000 kết nối
func BenchmarkGet(cache *Cache, wg *sync.WaitGroup, numConnections int) {
	startTime := time.Now()

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Get(fmt.Sprintf("key%d", i))
		}(i)
	}

	wg.Wait()

	endTime := time.Now()
	fmt.Printf("Thời gian Get cho %d kết nối: %v\n", numConnections, endTime.Sub(startTime))
}

func main() {
	// Khởi tạo cache
	cache := NewCache()
	var wg sync.WaitGroup

	// Số lượng kết nối giả định
	numConnections := 1000000

	// Chạy benchmark cho Set với Apache Arrow
	BenchmarkSet(cache, &wg, numConnections)

	// Chạy benchmark cho Get
	BenchmarkGet(cache, &wg, numConnections)
}
