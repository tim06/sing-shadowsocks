package ntp

import (
 "github.com/beevik/ntp"
 "sync"
 "time"
)

type NTPClient struct {
 ntpServer   string
 offset      time.Duration
 mutex       sync.Mutex
 lastSynced  time.Time
}

func NewNTPClient(ntpServer string) (*NTPClient, error) {
 client := &NTPClient{
  ntpServer: ntpServer,
 }

 if err := client.syncTime(); err != nil {
  return nil, err
 }

 return client, nil
}

func (c *NTPClient) syncTime() error {
 c.mutex.Lock()
 defer c.mutex.Unlock()

 response, err := ntp.Query(c.ntpServer)
 if err != nil {
  return err
 }

 err = response.Validate()
 if err != nil {
  return err
 }

 c.offset = response.ClockOffset
 c.lastSynced = time.Now()

 return nil
}

func (c *NTPClient) Now() time.Time {
 c.mutex.Lock()
 defer c.mutex.Unlock()

 return time.Now().Add(c.offset)
}