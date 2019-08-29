package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id       int
	UserName string
	Password string
	Mobile   string
	Avatar   string
	Status   uint8
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User))
}
