package ntp

import (
    "github.com/beevik/ntp"
    "log"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

type NTPClient struct {
    ntpServers []string
    offset     time.Duration
    mutex      sync.Mutex
    lastSynced time.Time
    synced     atomic.Bool
}

func NewNTPClient(ntpServers ...string) *NTPClient {
    client := &NTPClient{
        ntpServers: ntpServers,
    }
    go client.UpdateTime()
    return client
}

func (c *NTPClient) UpdateTime() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    for _, server := range c.ntpServers {
        response, err := ntp.Query(server)
        if err != nil {
            log.Println(fmt.Sprintf("Error querying NTP server %s: %v", server, err))
            continue
        }
        err = response.Validate()
        if err != nil {
            log.Println(fmt.Sprintf("Validation error from NTP server %s: %v", server, err))
            continue
        }
        c.offset = response.ClockOffset
        c.lastSynced = time.Now()
        c.synced.Store(true)
        log.Println(fmt.Sprintf("Time synchronized with NTP server %s", server))
        return
    }

    c.synced.Store(false)
    log.Println(fmt.Sprintf("Failed to synchronize time with all provided NTP servers"))
}

func (c *NTPClient) Now() time.Time {
    if !c.synced.Load() {
        go c.UpdateTime()
        if !c.synced.Load() {
            return time.Now()
        }
    }

    c.mutex.Lock()
    defer c.mutex.Unlock()

    return time.Now().Add(c.offset)
}
