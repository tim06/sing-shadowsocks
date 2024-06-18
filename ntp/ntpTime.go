package ntp

import (
    "github.com/beevik/ntp"
    "sync"
    "time"
)

type NTPClient struct {
    ntpServers  []string
    offset      time.Duration
    mutex       sync.Mutex
    lastSynced  time.Time
    synced      bool
   }

   func NewNTPClient(ntpServers ...string) *NTPClient {
    client := &NTPClient{
     ntpServers: ntpServers,
    }
    client.UpdateTime()
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
     c.synced = true
     return
    }

    c.synced = false
   }

   func (c *NTPClient) Now() time.Time {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if !c.synced {
     c.UpdateTime()
     if !c.synced {
      return time.Now()
     }
    }

    return time.Now().Add(c.offset)
   }