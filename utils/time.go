package utils

import (
    "database/sql/driver"
    "fmt"
    "time"

    "github.com/zhengow/vngo/consts"
)

type DatabaseTime time.Time

func (t *DatabaseTime) UnmarshalJSON(data []byte) (err error) {
    now, err := time.ParseInLocation(consts.DateFormat, string(data), time.Local)
    *t = DatabaseTime(now)
    return
}

func (t *DatabaseTime) MarshalJSON() ([]byte, error) {
    tTime := time.Time(*t)
    return []byte(tTime.Format(consts.DateFormat)), nil
}

func (t DatabaseTime) Value() (driver.Value, error) {
    return time.Time(t).Format(consts.DateFormat), nil
}

func (t *DatabaseTime) Scan(v interface{}) error {
    switch vt := v.(type) {
    case string:
        tTime, err := time.Parse(consts.DateFormat, vt)
        if err != nil {
            return err
        }
        *t = DatabaseTime(tTime)
    case time.Time:
        *t = DatabaseTime(vt)
    default:
        return fmt.Errorf("unknown err: %v", v)
    }
    return nil
}

func (t DatabaseTime) Format() string {
    return time.Time(t).Format(consts.DateFormat)
}
