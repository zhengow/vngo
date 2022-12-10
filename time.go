package vngo

import (
    "database/sql/driver"
    "fmt"
    "time"
)

type DatabaseTime struct {
    time.Time
}

func (t *DatabaseTime) UnmarshalJSON(data []byte) (err error) {
    now, err := time.ParseInLocation(DateFormat, string(data), time.Local)
    *t = DatabaseTime{
        now,
    }
    return
}

func (t *DatabaseTime) MarshalJSON() ([]byte, error) {
    tTime := t.Time
    return []byte(tTime.Format(DateFormat)), nil
}

func (t DatabaseTime) Value() (driver.Value, error) {
    return t.Format(DateFormat), nil
}

func (t *DatabaseTime) Scan(v interface{}) error {
    switch vt := v.(type) {
    case string:
        tTime, err := time.Parse(DateFormat, vt)
        if err != nil {
            return err
        }
        *t = DatabaseTime{
            tTime,
        }
    case time.Time:
        *t = DatabaseTime{
            vt,
        }
    default:
        return fmt.Errorf("unknown err: %v", v)
    }
    return nil
}
