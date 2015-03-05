package shared

import (
	"github.com/jinzhu/gorm"
)

type SharedContext struct {
    PersistentDB gorm.DB
    DB gorm.DB
}

