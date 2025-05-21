# sse – Server‑Sent Events for Go

> A tiny, dependency‑free package that makes it easy to use SSE streaming with any HTTP framework.

---

## ✨ Features

* **Spec‑compliant**: Follows the [official WHATWG spec](https://html.spec.whatwg.org/multipage/server-sent-events.html) for server-side events.
* **Zero dependencies**: Depends only on the go standard library.
* **Well-tested** Comprehensive test suite with >90% coverage, including end-to-end tests.
* **Framework Agnostic**: Designed to be used with any Go HTTP framework.
* **Flexible Encoding**: Supports both JSON and raw data encoding.

## Quick Start

### Server (Sending Events)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    
    "go.jetify.com/sse"
)

func eventsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to SSE
    conn, err := sse.Upgrade(r.Context(), w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close()
    
    // Send events until client disconnects
    count := 0
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Send a simple event with JSON data
            event := &sse.Event{
                ID:    fmt.Sprintf("%d", count),
                Data:  map[string]any{"count": count, "time": time.Now().Format(time.RFC3339)},
            }
            
            if err := conn.SendEvent(r.Context(), event); err != nil {
                log.Printf("Error: %v", err)
                return
            }
            count++
            
        case <-r.Context().Done():
            return // Client disconnected
        }
    }
}

func main() {
    http.HandleFunc("/events", eventsHandler)
    log.Println("SSE server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Client (Receiving Events)

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    "go.jetify.com/sse"
)

func main() {
    // Make request to SSE endpoint
    resp, err := http.Get("http://localhost:8080/events")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer resp.Body.Close()
    
    // Create decoder for the SSE stream
    decoder := sse.NewDecoder(resp.Body)
    
    // Process events as they arrive
    for {
        var event sse.Event
        err := decoder.Decode(&event)
        if err != nil {
            log.Printf("Stream ended: %v", err)
            break
        }
        
        fmt.Printf("Event ID: %s, Data: %v\n", event.ID, event.Data)
    }
}
```

## Examples

See the [examples](examples/) directory for more patterns:
- [Basic Server](examples/1_basic_server/main.go)
- [Custom Events](examples/2_custom_events/main.go)
- [Reconnection Client](examples/3_reconnection_client/main.go)
- [Middleware Integration](examples/4_middleware_integration/main.go)