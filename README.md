# Go Git API Client

A simple Go client for interacting with REST APIs.

## Installation

```bash
go get github.com/ritu-p/go-git/api-client
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/go-git/api-client"
)

func main() {
    client := apiclient.NewClient("api-service-url")
    user, err := client.GetUser("test-id")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("User Name: %+v\n", user.name)
}
```

## Features

- Create User
- Get User
- Update User Info

## Configuration

Pass the url of you API to the clients

## License

MIT