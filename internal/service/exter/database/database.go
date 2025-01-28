package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/erx/erx"
	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var initOnce sync.Once

type zerologWriter struct{}

func (z *zerologWriter) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func New() error {
	var err error

	initOnce.Do(func() {
		log.Info().Msg("New db")
		config := gorm.Config{
			// NamingStrategy: schema.NamingStrategy{
			// 	NoLowerCase: true,
			// },
			Logger: logger.New(
				&zerologWriter{},
				logger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  logger.Warn,
					IgnoreRecordNotFoundError: false,
					Colorful:                  true,
				},
			),
		}

		switch global.Database {
		case "mysql":
			log.Info().Msg("Connecting to MySQL")
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
				"admin",
				"admin",
				"127.0.0.1",
				3306,
				"mydb",
				"UTC",
			)

			global.DB, err = gorm.Open(mysql.Open(dsn), &config)
		case "postgres":
			log.Info().Msg("Connecting to Postgres")
			connStr := fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s TimeZone=%s",
				viper.GetString("db.host"),
				viper.GetString("db.port"),
				viper.GetString("db.user"),
				viper.GetString("db.password"),
				viper.GetString("db.db"),
				viper.GetString("db.timezone"),
			)
			global.DB, err = gorm.Open(postgres.Open(connStr), &config)

			if err == nil {
				sql := `
						CREATE EXTENSION IF NOT EXISTS hstore;

		CREATE OR REPLACE FUNCTION check_permission(user_id TEXT, permission_type TEXT, object_id TEXT)
		RETURNS BOOLEAN AS $$
		DECLARE
		    queue TEXT[] := ARRAY[]::TEXT[];
		    visited hstore := hstore('');  -- Use hstore to track visited nodes
		    current TEXT;
		BEGIN
		    -- Initialize queue with nodes from 'belongs_to' relation
		    queue := queue || ARRAY(
		        SELECT to_v
		        FROM "edges"
		        WHERE from_v = user_id
		          AND relation = 'belongs_to'
		    );

		    -- BFS traversal
		    WHILE array_length(queue, 1) > 0 LOOP
		        -- Dequeue the first element
		        current := queue[1];
		        queue := queue[2:array_length(queue, 1)];

		        -- Skip if already visited
		        IF visited -> current IS NOT NULL THEN
		            CONTINUE;
		        END IF;

		        -- Mark the node as visited
		        visited := visited || hstore(current, 'visited');

		        -- Check if the permission exists
		        PERFORM 1
		        FROM "edges"
		        WHERE from_v = current
		          AND relation = permission_type
		          AND to_v = object_id;

		        IF FOUND THEN
		            RETURN TRUE;
		        END IF;

		        -- Enqueue neighbors (leader_of relation) if not visited
		        queue := queue || ARRAY(
		            SELECT to_v
		            FROM "edges"
		            WHERE from_v = current
		              AND relation = 'leader_of'
		              AND visited -> to_v IS NULL  -- Only enqueue if not visited
		        );
		    END LOOP;

		    RETURN FALSE;
		END;
		$$ LANGUAGE plpgsql;
				`
				if global.DB.Raw(sql); err != nil {
					err = erx.W(err, "Initialize database failed")
					return
				}
			}
		}

		if err != nil {
			err = erx.W(err, "Initialize database failed")
			return
		}

		if global.AutoMigrate {
			if err = global.DB.AutoMigrate(
				&model.Edge{},
			); err != nil {
				err = erx.W(err, "Migration failed")
				return
			}
		}
	})
	return err
}

func Close() error {
	if global.DB == nil {
		return nil
	}

	db, err := global.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
