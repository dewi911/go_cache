# go_cache_in_memory

```bash
go get github.com/dewi911/go_cache
```
simple code
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dewi911/go_cache"
)

func main() {
	cacheConfig := cache.Config{
		CleanupInterval: 5 * time.Minute,
		MaxSize:         100,
	}
	myCache := cache.New(cacheConfig)
	
	ctx := context.Background()
	
	err := myCache.Set(ctx, "key1", "value1", 10*time.Minute)
	if err != nil {
		fmt.Printf("Ошибка при добавлении key1: %v\n", err)
	}

	err = myCache.Set(ctx, "key2", 42, 1*time.Hour)
	if err != nil {
		fmt.Printf("Ошибка при добавлении key2: %v\n", err)
	}

	// Получаем значения из кэша
	value1, err := myCache.Get(ctx, "key1")
	if err != nil {
		fmt.Printf("Ошибка при получении key1: %v\n", err)
	} else {
		fmt.Printf("Значение key1: %v\n", value1)
	}

	value2, err := myCache.Get(ctx, "key2")
	if err != nil {
		fmt.Printf("Ошибка при получении key2: %v\n", err)
	} else {
		fmt.Printf("Значение key2: %v\n", value2)
	}

	// Пробуем получить несуществующий ключ
	value3, err := myCache.Get(ctx, "key3")
	if err != nil {
		fmt.Printf("Ошибка при получении key3: %v\n", err)
	} else {
		fmt.Printf("Значение key3: %v\n", value3)
	}

	// Удаляем элемент из кэша
	err = myCache.Del(ctx, "key1")
	if err != nil {
		fmt.Printf("Ошибка при удалении key1: %v\n", err)
	}

	// Пробуем получить удаленный элемент
	value1Again, err := myCache.Get(ctx, "key1")
	if err != nil {
		fmt.Printf("Ошибка при получении удаленного key1: %v\n", err)
	} else {
		fmt.Printf("Значение удаленного key1: %v\n", value1Again)
	}

	// Очищаем весь кэш
	err = myCache.Flush(ctx)
	if err != nil {
		fmt.Printf("Ошибка при очистке кэша: %v\n", err)
	}

	fmt.Println("Кэш очищен")
}
```