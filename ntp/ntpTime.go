package ntp

import (
 "github.com/beevik/ntp"
 "sync"
 "time"
)

type NTPClient struct {
 ntpServer    string
 offset       time.Duration
 mutex        sync.Mutex
 lastSynced   time.Time
 synced       bool
}

func NewNTPClient(ntpServer string) *NTPClient {
 client := &NTPClient{
  ntpServer: ntpServer,
 }
 client.UpdateTime()
 return client
}

func (c *NTPClient) UpdateTime() {
 c.mutex.Lock()
 defer c.mutex.Unlock()

 response, err := ntp.Query(c.ntpServer)
 if err != nil {
  c.synced = false
  return
 }
 err = response.Validate()
 if err != nil {
  c.synced = false
  return
 }
 c.offset = response.ClockOffset
 c.lastSynced = time.Now()
 c.synced = true
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