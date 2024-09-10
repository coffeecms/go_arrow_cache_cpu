# High-Performance Cache Management System with Apache Arrow - https://blog.lowlevelforest.com/

This project implements a high-performance key-value cache system in Go, utilizing **Apache Arrow** for efficient memory handling and supporting large-scale connections. The system supports caching data with a **Time To Live (TTL)** mechanism to automatically expire stale items.

## Features

- **Key-Value Cache**: Stores and retrieves key-value pairs in memory.
- **Apache Arrow Integration**: Uses Apache Arrow to optimize memory usage for binary data storage.
- **TTL Support**: Each cache entry expires after a user-defined period (Time To Live).
- **Concurrency**: Efficient thread-safe management of concurrent cache access.

## Prerequisites

- Go 1.18 or higher
- Apache Arrow Go library (`github.com/apache/arrow/go/v14`)
  
  Install using:
  ```bash
  go get github.com/apache/arrow/go/v14
  ```

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

### 1. Initializing the Cache

You need to create a new cache instance using the `NewCache()` function:

```go
cache := NewCache()
```

### 2. Setting a Key-Value Pair with TTL

To store a key-value pair in the cache with a **Time To Live (TTL)**, use the `SetWithArrow` method. The following example demonstrates how to set a key with a value and a TTL of 1 minute:

```go
cache.SetWithArrow("example_key", []byte("This is some data"), 1*time.Minute)
```

- **Key**: A string representing the key for the cache entry.
- **Value**: A byte slice representing the data you want to store.
- **TTL (Time To Live)**: The duration after which the cache entry will expire. In this example, it's set to `1*time.Minute`, but you can adjust it as needed.

#### Example:
```go
cache.SetWithArrow("user123", []byte("User data for ID 123"), 5*time.Minute)
```

This stores the data `"User data for ID 123"` under the key `user123` and the entry will expire after 5 minutes.

### 3. Retrieving a Value from the Cache

To retrieve a cached value by its key, use the `Get` method. It returns the value and a boolean indicating whether the key was found and is still valid (not expired).

```go
value, found := cache.Get("example_key")
if found {
    fmt.Printf("Value: %s\n", string(value))
} else {
    fmt.Println("Key not found or expired")
}
```

- **Key**: The key for which you want to retrieve the value.
- **Return Values**: 
  - **Value**: The byte slice stored under the key (if found).
  - **Found**: A boolean (`true` if the key exists and hasn't expired, `false` otherwise).

#### Example:
```go
value, found := cache.Get("user123")
if found {
    fmt.Printf("Retrieved Value: %s\n", string(value))
} else {
    fmt.Println("The key has expired or does not exist.")
}
```

This retrieves the value for the key `user123`. If the key has expired or doesn't exist, the `found` variable will be `false`.

### 4. Removing Expired Items from the Cache

The cache automatically handles TTL, but you can manually clean up expired items using the `CleanExpiredItems` method:

```go
cache.CleanExpiredItems()
```

This function iterates over all items and removes those that have expired.

## Structure

- `CacheItem`: Represents a single cache entry with the data and expiration time.
- `Cache`: Manages cache storage, including setting and getting key-value pairs, with TTL.
- `Apache Arrow`: Used to store binary data efficiently, improving memory management when handling large datasets.

## How It Works

1. **Setting Cache Items**: 
   - Data is stored in memory using **Apache Arrow**'s `BinaryBuilder`, which helps efficiently manage binary data storage.
   - Each item is associated with a TTL (Time To Live), ensuring that stale data is automatically removed.

2. **Retrieving Cache Items**:
   - The cache checks if the key exists and whether the TTL has expired before returning the value.

3. **Concurrency**:
   - The cache system uses `sync.RWMutex` to handle concurrent reads and writes safely.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
