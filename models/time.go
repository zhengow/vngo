package models

import (
    "database/sql/driver"
    "fmt"
    "time"
)

type VnTime struct {
    time.Time
}

func NewVnTime(t time.Time) VnTime {
    return VnTime{
        t,
    }
}

func (t *VnTime) Format() string {
    return t.Time.Format(DateFormat)
}

func (t *VnTime) UnmarshalJSON(data []byte) (err error) {
    now, err := time.ParseInLocation(DateFormat, string(data), time.Local)
    *t = VnTime{
        now,
    }
    return
}

func (t *VnTime) MarshalJSON() ([]byte, error) {
    tTime := t.Time
    return []byte(tTime.Format(DateFormat)), nil
}

func (t VnTime) Value() (driver.Value, error) {
    return t.Format(), nil
}

func (t *VnTime) Scan(v interface{}) error {
    switch vt := v.(type) {
    case string:
        tTime, err := time.Parse(DateFormat, vt)
        if err != nil {
            return err
        }
        *t = VnTime{
            tTime,
        }
    case time.Time:
        *t = VnTime{
            vt,
        }
    default:
        return fmt.Errorf("unknown err: %v", v)
    }
    return nil
}
