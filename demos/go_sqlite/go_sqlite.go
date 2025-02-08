package go_sqlite

import (
	"fmt"
	"strconv"
	"sync"

	"time"

	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbfile = "demos/go_sqlite/test.db"

func Run() {
	gormDB, sqlDB, err := InitDB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	users := []User{}
	for i := 0; i < 1000; i++ {
		user := User{
			Name:      "user_" + strconv.Itoa(i),
			Age:       uint8(i % 100),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		users = append(users, user)
	}
	err = BatchInsertUsers(gormDB, users)
	if err != nil {
		panic(err)
	}

	users, err = GetUsers(gormDB)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(users))
	fmt.Println(users[0])
}

type User struct {
	ID        uint
	Name      string
	Age       uint8
	CreatedAt int64
	UpdatedAt int64
}

func InitDB() (*gorm.DB, *sql.DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, _ := gormDB.DB()
	gormDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDB, sqlDB, nil
}

func BatchInsertUsers(gormDB *gorm.DB, users []User) error {
	batchSize := 100
	batchCount := (len(users) + batchSize - 1) / batchSize
	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(users) {
			end = len(users)
		}
		batch := users[start:end]
		tx := gormDB.Begin()
		if err := tx.Error; err != nil {
			return err
		}
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

func GetUsers(gormDB *gorm.DB) ([]User, error) {
	var users []User
	err := gormDB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

var wg sync.WaitGroup

func Run2() {
	gormDB, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB, _ := gormDB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	wg.Add(2000)

	// 并发写入 1000 条数据
	for i := 0; i < 1000; i++ {
		go func(i int) {
			defer wg.Done()
			err := gormDB.Transaction(func(tx *gorm.DB) error {
				user := User{Name: fmt.Sprintf("user_%d", i)}
				result := tx.Create(&user)
				return result.Error
			})
			if err != nil {
				fmt.Printf("failed to write data: %v\n", err)
			}
		}(i)
	}

	// 并发读取数据
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			var users []User
			err := gormDB.Transaction(func(tx *gorm.DB) error {
				result := tx.Find(&users)
				return result.Error
			})
			if err != nil {
				fmt.Printf("failed to read data: %v\n", err)
			} else {
				fmt.Printf("read %d records\n", len(users))
			}
		}()
	}

	wg.Wait()

	fmt.Println("done")
}
