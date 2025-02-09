package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
	s "golang.org/x/sync/singleflight"
)

type Product struct {
	ID   int
	Name string
}

type Singleflight[T any] struct {
	Group *s.Group
	Key   string
}

func (single *Singleflight[T]) ProccesWrapper(fn func() (T, error)) (T, error) {
	wrapperFn := func() (interface{}, error) {
		return fn()
	}

	res, err, _ := single.Group.Do(single.Key, wrapperFn)

	// Type assertion check
	if result, ok := res.(T); ok {
		return result, err
	}

	// Handle type assertion failure gracefully
	err = errors.New("unexpected type assertion failure")
	return *new(T), err
}

func (single *Singleflight[T]) Forget(keys ...string) {
	for _, key := range keys {
		single.Group.Forget(key)
	}
}

func getProductFromCache(rdb *redis.Client, sGroup *s.Group, productID int, currIdx int) (*Product, error) {

	singleflightInstance := Singleflight[*Product]{
		Group: sGroup,
		Key:   fmt.Sprintf("singleflight:product:%v", productID),
	}

	if currIdx == 2 {
		singleflightInstance.Forget(fmt.Sprintf("singleflight:product:%v", productID))
	}

	// get the product from cache
	res, err := singleflightInstance.ProccesWrapper(func() (*Product, error) {
		val, err := rdb.Get(context.Background(), fmt.Sprintf("product:%v", productID)).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, nil
			}
			msg := fmt.Sprintf("Error: %v", err)
			fmt.Println(msg)
			return nil, err
		}

		// unmarshal the value
		var product Product
		if err := json.Unmarshal([]byte(val), &product); err != nil {
			err = errors.Wrap(err, "Failed to unmarshal product")
			return nil, err
		}
		return &product, nil
	})

	if err != nil {
		err = errors.Wrap(err, "Failed to get product from cache")
		return nil, err
	}
	return res, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("Connected to Redis")

	// example product instance
	product := Product{
		ID:   1,
		Name: "Product 1",
	}
	sGroup := s.Group{}

	// marshal the product instance
	productBytes, err := json.Marshal(product)
	if err != nil {
		msg := fmt.Sprintln("Failed to marshal product")
		fmt.Println(msg)
		return
	}

	// set the product instance to redis
	err = rdb.Set(context.Background(), fmt.Sprintf("product:%v", product.ID), productBytes, 0).Err()
	if err != nil {
		msg := fmt.Sprintf("Failed to set product to cache %v", err)
		fmt.Println(msg)
		return
	}

	var wg sync.WaitGroup
	// get the product instance from cache
	// use a loop to simulate multiple requests
	for i := 0; i < 10; i++ {
		wg.Add(1)
		idxPtr := &i
		go func(idx *int) {
			// if idx is 2, wait for 5 seconds
			if *idx == 2 {
				fmt.Println("Sleeping for 5 seconds")
				time.Sleep(5 * time.Second)
			}
			defer wg.Done()
			_, err := getProductFromCache(rdb, &sGroup, product.ID, *idx)
			if err != nil {
				msg := fmt.Sprintf("Error: %v", err)
				fmt.Println(msg)
				return
			}
		}(idxPtr)
	}

	wg.Wait()
}
