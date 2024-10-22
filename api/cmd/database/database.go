package database

import "xorm.io/xorm"

type Database interface {
  GetDb() *xorm.Engine
}