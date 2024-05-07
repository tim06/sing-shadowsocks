package ntp

import (
 "github.com/beevik/ntp"
 "sync"
 "time"
)

type NTPClient struct {
 ntpServer string
 mutex     sync.Mutex
 ntpTime   time.Time
}

func NewNTPClient(ntpServer string) *NTPClient {
 return &NTPClient{
  ntpServer: ntpServer,
  ntpTime:   time.Now(),
 }
}

func (c *NTPClient) UpdateTime() error {
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

 c.ntpTime = time.Now().Add(response.ClockOffset)
 return nil
}

func (c *NTPClient) Now() time.Time {
 c.mutex.Lock()
 defer c.mutex.Unlock()

 // Возвращаем скорректированное время.
 return c.ntpTime.Add(time.Since(c.ntpTime))
}