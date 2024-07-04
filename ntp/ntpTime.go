package ntp

import (
    "github.com/beevik/ntp"
    "log"
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
            log.Printf("Error querying NTP server %s: %v", server, err)
            continue
        }
        err = response.Validate()
        if err != nil {
            log.Printf("Validation error from NTP server %s: %v", server, err)
            continue
        }
        c.offset = response.ClockOffset
        c.lastSynced = time.Now()
        c.synced.Store(true)
        log.Printf("Time synchronized with NTP server %s", server)
        return
    }

    c.synced.Store(false)
    log.Println("Failed to synchronize time with all provided NTP servers")
}

func (c *NTPClient) Now() time.Time {
    if !c.synced.Load() {
        go c.UpdateTime() // Запуск UpdateTime в другом потоке, если еще не синхронизировано
        if !c.synced.Load() {
            return time.Now()
        }
    }

    c.mutex.Lock()
    defer c.mutex.Unlock()

    return time.Now().Add(c.offset)
}
