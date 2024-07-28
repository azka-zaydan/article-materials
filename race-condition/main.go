package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pool := goredis.NewPool(rdb)

	rs := redsync.New(pool)
	err := AddToBankAccountWithMutex("", 100, rs)
	if err != nil {
		return
	}
}

func AddToBankAccountWithMutex(accountId string, amount int, redSync *redsync.Redsync) (err error) {
	// create the mutex with account id
	mutex := redSync.NewMutex(fmt.Sprintf("add-account:{%s}", accountId))

	// lock the mutex, it will fail if the mutex with the same name already exists
	if err = mutex.Lock(); err != nil {
		return
	}

	// we unlock after the function has done running or if an error occurs
	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return
		}
	}()

	// put logic here

	return
}

func AddToBankAccount(accountId string, amount int, rdb *redis.Client) (err error) {
	// we first check if the key already exist, if not then continue\
	exist := true
	err = rdb.Get(context.Background(), fmt.Sprintf("add-account:{%s}", accountId)).
		Err()
	if err != nil {
		// if the error is anything other than redis nil, than we return the error
		if err != redis.Nil {
			return
		}
		exist = false
	}

	if exist {
		return
	}
	// set the key
	err = rdb.Set(context.Background(), fmt.Sprintf("add-account:{%s}", accountId), accountId, time.Minute*10).Err()
	if err != nil {
		return
	}
	// delete the key after the function is done
	defer func() {
		rdb.Del(context.Background(), fmt.Sprintf("add-account:{%s}", accountId))
	}()

	// put logic here

	return
}
