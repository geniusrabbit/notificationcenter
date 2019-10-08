package pg

import (
	"database/sql"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	nc "github.com/geniusrabbit/notificationcenter"
)

const (
	testPostgresConnection = "TEST_PGCONNECTION"
	triggerFunction        = `
CREATE OR REPLACE FUNCTION test_notify_event() RETURNS TRIGGER AS $$
    DECLARE
        data json;
        notification json;
    BEGIN
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);


        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('test_events', notification::text);

        RAISE NOTICE '(%)', notification::text;


        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;`
	testData = `
CREATE TABLE IF NOT EXISTS test_notifies
( id BIGSERIAL PRIMARY KEY
, vl TEXT
);
CREATE TRIGGER test_notifies_event
AFTER INSERT OR UPDATE OR DELETE ON test_notifies
    FOR EACH ROW EXECUTE PROCEDURE test_notify_event();
`
	testDataOperation = `INSERT INTO test_notifies (vl) VALUES(md5(random()::text))`
	testDataEracer    = `
DROP TABLE IF EXISTS test_notifies;
DROP FUNCTION IF EXISTS test_notify_event;
`
	testEventCount = 10
)

// docker run --rm -p 54320:5432 postgres:11
// export TEST_PGCONNECTION=postgres://postgres@localhost:54320/postgres?sslmode=disable

func Test_EventListening(t *testing.T) {
	connection := os.Getenv(testPostgresConnection)
	if connection == "" {
		t.Skip()
		return
	}

	db, err := sql.Open("postgres", connection)
	if err != nil {
		t.Error(err)
		return
	}

	execSQL(t, db, triggerFunction)
	execSQL(t, db, testData)
	defer execSQL(t, db, testDataEracer)

	// Subsribe on the notification
	subscriber, err := NewSubscriber(connection, "test_events", nil)
	if err != nil {
		t.Error(err)
		return
	}

	messageCount := int32(0)
	subscriber.Subscribe(nc.FuncHandler(func(msg nc.Message) error {
		if !strings.HasPrefix(msg.ID(), "test_events") {
			t.Error("invalid event message:", msg.ID())
		} else {
			atomic.AddInt32(&messageCount, 1)
		}
		return nil
	}))
	go subscriber.Listen()

	for i := 0; i <= testEventCount; i++ {
		db.Exec(testDataOperation)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)

	if cnt := atomic.LoadInt32(&messageCount); cnt != testEventCount {
		t.Errorf("not all events was delivered: %d of %d", cnt, testEventCount)
	}
}

func execSQL(t *testing.T, conn *sql.DB, query string) {
	if _, err := conn.Exec(query); err != nil {
		t.Error(err)
	}
}
