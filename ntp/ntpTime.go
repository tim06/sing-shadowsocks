package ntp

import (
    "github.com/beevik/ntp"
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
    go client.UpdateTime() // Запуск UpdateTime в другом потоке
    return client
}

func (c *NTPClient) UpdateTime() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    for _, server := range c.ntpServers {
        response, err := ntp.Query(server)
        if err != nil {
            continue
        }
        err = response.Validate()
        if err != nil {
            continue
        }
        c.offset = response.ClockOffset
        c.lastSynced = time.Now()
        c.synced.Store(true)
        return
    }

    c.synced.Store(false)
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
