package models

import (
	"fmt"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/juju/loggo"
)

var db *gorm.DB
var logger *loggo.Logger

func CloseDB() {
	db.Close()
	return
}

func DecodeEngine(engine string) (dialect string, args string) {
	pgRe, err := regexp.Compile(`postgresql://([\w]*):([\w\-.~:/?#\[\]!$&'()*+,;=]*)@([\w.]*)/([\w]*)`)
	if err != nil {
		logger.Criticalf("Regex compile error: %s", err)
		panic("PANIC!")
	}

	if pgRe.MatchString(engine) {
		dialect = "postgres"
		match := pgRe.FindStringSubmatch(engine)
		args = fmt.Sprintf("host=%s user=%s dbname=%s password=%s", match[3], match[1], match[4], match[2])
	} else {
		logger.Criticalf("Could not parse DB_ENGINE: %s", err)
		panic("PANIC!")
	}

	return
}

func InitDB(connectionString string) {
	newLogger :=  loggo.GetLogger("puphaus.models")
	logger = &newLogger

	var err error
	dialect, dbArgs := DecodeEngine(connectionString)
	db, err = gorm.Open(dialect, dbArgs)
	if err != nil {
		logger.Criticalf("Coud not connect to database: %s", err)
		panic(err)
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return "api_" + defaultTableName;
	}

	db.AutoMigrate(&User{})
	db.Model(&User{}).AddUniqueIndex("uidx_api_users_username_key", "lower(username)")

	logger.Infof("Connected to %s database", dialect)

	// Create admin if no users present
	if uc := UserCount();uc == 0 {
		logger.Infof("Creating admin user")
		SetUser("admin","admin")
	}

	return
}
