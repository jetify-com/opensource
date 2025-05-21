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
    "context"
    "log"
    "net/http"
    "time"
    
    "go.jetify.com/sse"
)

func eventsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade the HTTP connection to SSE
    conn, err := sse.Upgrade(r.Context(), w,
        sse.WithHeartbeatInterval(15*time.Second),
        sse.WithRetryDelay(3*time.Second),
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close()
    
    // Send events in a loop
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for i := 0; ; i++ {
        select {
        case <-ticker.C:
            event := &sse.Event{
                ID:    fmt.Sprintf("event-%d", i),
                Event: "message",
                Data:  map[string]any{
                    "count": i,
                    "time":  time.Now().Format(time.RFC3339),
                },
            }
            
            if err := conn.SendEvent(r.Context(), event); err != nil {
                log.Printf("Error sending event: %v", err)
                return
            }
            
        case <-r.Context().Done():
            log.Println("Client disconnected")
            return
        }
    }
}

func main() {
    http.HandleFunc("/events", eventsHandler)
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Client (Receiving Events)

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
    
    "go.jetify.com/sse"
)

func main() {
    ctx := context.Background()
    
    // Client with reconnection handling
    for {
        if err := connectAndProcess(ctx); err != nil {
            if !errors.Is(err, io.EOF) {
                log.Printf("Connection error: %v, reconnecting...", err)
                time.Sleep(3 * time.Second) // Default reconnection delay
                continue
            }
            log.Println("Stream ended normally")
            break
        }
    }
}

func connectAndProcess(ctx context.Context) error {
    resp, err := http.Get("http://localhost:8080/events")
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    decoder := sse.NewDecoder(resp.Body)
    var lastEventID string
    var retryDelay time.Duration = 3 * time.Second // Default retry delay
    
    for {
        var event sse.Event
        err := decoder.Decode(&event)
        if err != nil {
            // On reconnect, we'll use the last event ID and retry delay we received
            return err
        }
        
        // Update our reconnection state
        if event.ID != "" {
            lastEventID = event.ID
        }
        
        // Use server-suggested retry delay if provided
        if serverDelay := decoder.RetryDelay(); serverDelay > 0 {
            retryDelay = serverDelay
            log.Printf("Server requested retry delay: %v", retryDelay)
        }
        
        fmt.Printf("Received event (ID: %s, Type: %s): %v\n", 
            event.ID, event.Event, event.Data)
    }
}