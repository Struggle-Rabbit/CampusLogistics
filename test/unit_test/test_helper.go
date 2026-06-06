package unittest

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	testCounter int64
	counterMu   sync.Mutex
)

func SetupTestDB() *gorm.DB {
	counterMu.Lock()
	testCounter++
	idx := testCounter
	counterMu.Unlock()

	tmpFile := fmt.Sprintf("%s/test_campus_db_%d_%d", os.TempDir(), time.Now().UnixNano(), idx)
	os.Remove(tmpFile)

	db, err := gorm.Open(sqlite.Open(tmpFile+"?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret",
			AccessExpire:  3600,
			RefreshExpire: 86400,
		},
	}
	config.GlobalConfig = cfg

	utils.InitSnowflake()

	db.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysOperationLog{},
		&model.RepairOrder{},
		&model.RepairRecord{},
		&model.Campus{},
		&model.Building{},
		&model.DormRoom{},
		&model.DormUser{},
		&model.DormUtility{},
		&model.UtilityPrice{},
		&model.Notice{},
	)

	return db
}
