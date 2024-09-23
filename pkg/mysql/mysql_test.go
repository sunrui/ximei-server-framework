/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-12-05 22:24:45
 */

package mysql

import (
	"encoding/json"
	"fmt"
	"framework/pkg/hypersonic"
	"framework/pkg/id"
	"testing"

	"gorm.io/sharding"
)

type Model struct {
	Name string `json:"name" gorm:"unique;type:varchar(256); comment:名称"`
}

type User struct {
	ModelId

	Model `gorm:"embedded"`
	Age   int    `json:"age" gorm:"default:18; comment:年龄"`                  // 年龄
	Class string `json:"class" gorm:"type:varchar(32); not null;comment:班级"` // 班级

	ModelTime
}

func (User) TableName() string {
	return "t_user"
}

type UserScore struct {
	ModelId

	Model `gorm:"embedded"`
	Score int `json:"score" gorm:"not null;check:score>=0&&score<=100;comment:分数"` // 分数

	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserId"`
	UserId string `json:"userId" gorm:"comment:用户 id"`

	ModelTime
}

func (UserScore) TableName() string {
	return "t_user_score"
}

var (
	db *Mysql
)

// TestMain 初始化前准备
func TestMain(m *testing.M) {
	var err error

	// 测试数据库连接
	db, err = New(Config{
		User:          "root",
		Password:      "root",
		Host:          "127.0.0.1",
		Port:          3306,
		Database:      "hypersonic_test",
		MaxOpenConns:  1,
		MaxIdleConns:  1,
		SlowThreshold: 50,
	})
	if err != nil {
		panic(err.Error())
	}

	// 删除数据库
	db.DB.Exec(`
		TRUNCATE TABLE t_user_score;
	`)
	db.DB.Exec(`
		TRUNCATE TABLE t_user;
	`)
	// 测试多个迁移一次
	_ = db.DB.AutoMigrate(User{}, UserScore{})

	m.Run()
}

func Test_AutoMigrate(t *testing.T) {
	_ = db.DB.AutoMigrate(User{}, UserScore{})
}

func Test_Insert(t *testing.T) {
	user := User{
		Model: Model{
			Name: "张三",
		},
		Age: 19,
	}

	userRepository := NewRepository[User](db)
	if u := userRepository.FindOne("name = ? And age = ?", "张三", 19); u == nil {
		userRepository.Save(&user)
	} else {
		user = *u
	}

	userScoreRepository := NewRepository[UserScore](db)
	count := userScoreRepository.Count()
	for i := count + 1; i < count+1+10; i++ {
		// 测试 beyond to
		userScore := UserScore{
			Model: Model{
				Name: fmt.Sprintf("语文 - %03d", +i),
			},
			UserId: user.Id,
			Score:  80,
		}
		userScoreRepository.Save(&userScore)
	}

	t.Log("ok")
}

func TestRepository_FindOne(t *testing.T) {
	userRepository := NewRepository[User](db)
	if user := userRepository.FindOne(&User{
		Model: Model{
			Name: "张三",
		},
	}); user == nil {
		t.Error("not have this id")
	} else {
		t.Log("\n" + user.Name + "\n")
	}
}

func TestRepository_FindAll(t *testing.T) {
	var userId string

	userRepository := NewRepository[User](db)
	if user := userRepository.FindOne(&User{
		Model: Model{
			Name: "张三",
		},
	}); user == nil {
		t.Error("not have this id")
		return
	} else {
		userId = user.Id
	}

	userScoreRepository := NewRepository[UserScore](db)
	userScoreList := userScoreRepository.FindAll("name ASC", &UserScore{
		UserId: userId,
	})

	t.Log(userScoreList)
}

func TestRepository_FindPage(t *testing.T) {
	Test_Insert(t)

	var userId string

	userRepository := NewRepository[User](db)
	if user := userRepository.FindOne(&User{
		Model: Model{
			Name: "张三",
		},
	}); user == nil {
		t.Error("not have this id")
		return
	} else {
		userId = user.Id
	}

	userScoreRepository := NewRepository[UserScore](db)

	for i := 0; ; i++ {
		userScorePage, pagination := userScoreRepository.FindPage(hypersonic.Page{
			Page:     i,
			PageSize: 4,
		}, "name ASC", &UserScore{
			UserId: userId,
		})

		for index, userScore := range userScorePage {
			userScoreBytes, _ := json.Marshal(userScore)
			t.Log(index, string(userScoreBytes))
		}

		if int64(pagination.Page.Page) == pagination.TotalPage {
			break
		}
	}
}

func TestRepository_DeleteById(t *testing.T) {
	var userId string

	userRepository := NewRepository[User](db)
	if user := userRepository.FindOne(&User{
		Model: Model{
			Name: "张三",
		},
	}); user == nil {
		t.Error("not have this id")
		return
	} else {
		userId = user.Id
	}

	userScoreRepository := NewRepository[UserScore](db)

	userScorePage, pagination := userScoreRepository.FindPage(hypersonic.Page{
		Page:     1,
		PageSize: 10,
	}, "name ASC", &UserScore{
		UserId: userId,
	})

	t.Log(pagination.TotalPage, pagination.TotalSize)

	for _, userScore := range userScorePage {
		var r bool
		r = userScoreRepository.SoftDeleteById(userScore.Id)
		t.Log("\n"+"SoftDeleteById userScore by id "+userScore.Id+", result =", r)
		r = userScoreRepository.DeleteById(userScore.Id)
		t.Log("\n"+"DeleteById userScore by id "+userScore.Id+", result =", r)
	}
}

func TestRepository_SoftDeleteById(t *testing.T) {
	userRepository := NewRepository[User](db)
	if user := userRepository.FindOne(&User{
		Model: Model{
			Name: "张三",
		},
	}); user == nil {
		t.Error("not have this id")
		return
	}

	userPage, pagination := userRepository.FindPage(hypersonic.Page{
		Page:     1,
		PageSize: 10,
	}, "name ASC", nil)

	t.Log(pagination.TotalPage, pagination.TotalSize)

	for _, user := range userPage {
		var r bool
		r = userRepository.SoftDeleteById(user.Id)
		t.Log("delete "+user.Id+" ", r)
	}
}

type Order struct {
	ID        int64 `gorm:"primarykey"`
	UserID    string
	ProductID int64
}

func TestSharding(t *testing.T) {
	shardingNumber := 16

	for i := 0; i < shardingNumber; i += 1 {
		table := fmt.Sprintf("orders_%02d", i)
		db.DB.Exec(`DROP TABLE IF EXISTS ` + table)
		db.DB.Exec(`CREATE TABLE ` + table + ` (
			id bigint PRIMARY KEY,
			user_id varchar(32),
			product_id bigint
		)`)
	}

	middleware := sharding.Register(sharding.Config{
		ShardingKey:         "user_id",
		NumberOfShards:      uint(shardingNumber),
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, "orders")
	err := db.DB.Use(middleware)
	if err != nil {
		fmt.Println(err)
	}

	// this record will insert to orders_02
	err = db.DB.Create(&Order{UserID: "2"}).Error
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 1000; i++ {
		// this record will insert to orders_03
		err = db.DB.Exec("INSERT INTO orders(user_id,product_id) VALUES(?,?)", id.Uuid(), i).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	// this will throw ErrMissingShardingKey err
	err = db.DB.Exec("INSERT INTO orders(product_id) VALUES(1)").Error
	fmt.Println(err)

	// this will redirect query to orders_02
	var orders []Order
	err = db.DB.Model(&Order{}).Where("user_id", int64(2)).Find(&orders).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", orders)

	// Raw SQL also supported
	db.DB.Raw("SELECT * FROM orders WHERE user_id = ?", int64(3)).Scan(&orders)
	fmt.Printf("%#v\n", orders)

	// this will throw ErrMissingShardingKey err
	err = db.DB.Model(&Order{}).Where("product_id", "1").Find(&orders).Error
	fmt.Println(err)

	// Update and Delete are similar to create and query
	err = db.DB.Exec("UPDATE orders SET product_id = ? WHERE user_id = ?", 2, int64(3)).Error
	fmt.Println(err) // nil

	err = db.DB.Exec("DELETE FROM orders WHERE product_id = 3").Error
	fmt.Println(err) // ErrMissingShardingKey
}
