package configuration

import (
    "github.com/MrWhok/IMK-FP-BACKEND/exception"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "math/rand"
    "os"
    "strconv"
    "time"
)

func NewDatabase(config Config) *gorm.DB {
    username := config.Get("DATASOURCE_USERNAME")
    password := config.Get("DATASOURCE_PASSWORD")
    host := config.Get("DATASOURCE_HOST")
    port := config.Get("DATASOURCE_PORT")
    dbName := config.Get("DATASOURCE_DB_NAME")
    maxPoolOpen, err := strconv.Atoi(config.Get("DATASOURCE_POOL_MAX_CONN"))
    maxPoolIdle, err := strconv.Atoi(config.Get("DATASOURCE_POOL_IDLE_CONN"))
    maxPollLifeTime, err := strconv.Atoi(config.Get("DATASOURCE_POOL_LIFE_TIME"))
    exception.PanicLogging(err)

    loggerDb := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,
            LogLevel:                  logger.Info,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    dsn := "host=" + host + " user=" + username + " password=" + password + " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: loggerDb,
    })
    exception.PanicLogging(err)

    sqlDB, err := db.DB()
    exception.PanicLogging(err)

    sqlDB.SetMaxOpenConns(maxPoolOpen)
    sqlDB.SetMaxIdleConns(maxPoolIdle)
    sqlDB.SetConnMaxLifetime(time.Duration(rand.Int31n(int32(maxPollLifeTime))) * time.Millisecond)

    return db
}
