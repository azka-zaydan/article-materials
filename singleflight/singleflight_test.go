package main

import (
	"sync"
	"testing"
	"time"

	s "golang.org/x/sync/singleflight"
)

// Simulated function that fetches product data (slow operation)
func fetchProduct() (*Product, error) {
	time.Sleep(300 * time.Millisecond) // Simulated expensive DB query
	return &Product{ID: 1, Name: "Expensive Product"}, nil
}

// Function fetching product without singleflight (independent calls)
func getProductWithoutSingleflight() (*Product, error) {
	return fetchProduct()
}

// Function using singleflight to fetch product
func getProductWithSingleflight(sGroup *s.Group) (*Product, error) {
	singleflightInstance := Singleflight[*Product]{
		Group: sGroup,
		Key:   "singleflight:product:1",
	}

	return singleflightInstance.ProccesWrapper(fetchProduct)
}

// Benchmark: Without singleflight (each request executes separately)
func BenchmarkWithoutSingleflight(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = getProductWithoutSingleflight()
		}()
	}
	wg.Wait()
}

// Benchmark: With singleflight (grouped calls)
func BenchmarkWithSingleflight(b *testing.B) {
	var sGroup s.Group
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = getProductWithSingleflight(&sGroup)
		}()
	}
	wg.Wait()
}
