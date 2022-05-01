package handler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
# SERVER_PORT
ServerPort : 8888

# DB_PORT
DBPort : 3306

# DB_ID
DBId : root

# DB_Password
DBPassword : root1234

# DB_Database(scheme)
DBDatabase : chatbot
*/

// SINGLETON //
var (
	instance *ConfigHandler
	once     sync.Once
)

type ConfigHandler struct {
	ServerPort string
	DBPort     string
	DBId       string
	DBPassword string
	DBDatabase string

	sync.RWMutex
}

func GetCHInstance() *ConfigHandler {

	once.Do(func() {
		instance = new(ConfigHandler)
	})
	return instance

}

func (ch *ConfigHandler) Read(configPath string) error {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("no config.txt ...")
	}

	ioReader, err := os.Open(configPath)
	if err != nil {
		return errors.New("[Read] opening ioReader failed")
	}

	scanner := bufio.NewScanner(ioReader)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "#") && len(line) != 0 {
			ch.setValue(line)
		}
	}

	return nil

}

func (ch *ConfigHandler) setValue(line string) {

	lineItems := strings.Split(line, ":")
	key := strings.Trim(lineItems[0], " ")
	value := strings.Trim(lineItems[1], " ")

	if key == "ServerPort" {
		ch.Lock()
		defer ch.Unlock()
		ch.ServerPort = value
	} else if key == "DBPort" {
		ch.Lock()
		defer ch.Unlock()
		ch.DBPort = value
	} else if key == "DBId" {
		ch.Lock()
		defer ch.Unlock()
		ch.DBId = value
	} else if key == "DBPassword" {
		ch.Lock()
		defer ch.Unlock()
		ch.DBPassword = value
	} else if key == "DBDatabase" {
		ch.Lock()
		defer ch.Unlock()
		ch.DBDatabase = value
	}

}

func (ch *ConfigHandler) ToString() {

	encryptedPw := ch.DBPassword[:2]
	for i := 1; i <= len(ch.DBPassword)-3; i++ {
		encryptedPw = encryptedPw + "*"
	}

	fmt.Sprintf("server port : %s", ch.ServerPort)
	fmt.Sprintf("DB port : %s, DB id : %s, DB passwd : %s", ch.DBPort, ch.DBId, ch.DBPassword)
	fmt.Sprintf("DB database : %s", ch.DBDatabase)
}

func (ch *ConfigHandler) IsEmpty() bool {
	if len(ch.ServerPort) == 0 || len(ch.DBPort) == 0 || len(ch.DBId) == 0 || len(ch.DBPassword) == 0 || len(ch.DBDatabase) == 0 {
		return true
	} else {
		return false
	}
}

func (ch *ConfigHandler) GetDBconnection() (*gorm.DB, error) {

	if ch.IsEmpty() {
		return nil, errors.New("empty handler")
	}
	datasource := ch.DBId + ":" + ch.DBPassword + "@tcp(localhost:" + ch.DBPort + ")/" + ch.DBDatabase + "?charset=utf8"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       datasource,
		DefaultStringSize:         512,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
