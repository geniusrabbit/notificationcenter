# Postgres event subscriber

- [Integration with PostgreSQL](#integration-with-postgresql)
  - [Trigger function](#trigger-function)
  - [Apply the trigger](#apply-the-trigger)
- [Using example](#using-example)
- [Dependencies](#dependencies)

The module provides the minimal functionality to subscribe to the PostgreSQL events.
In the implementation used to NOTIFY/LISTEN functionality from the PostgreSQL.
For more details read the documentation or follow the [link](https://www.postgresql.org/docs/current/sql-notify.html).

The NOTIFY command sends a notification event together with an optional “payload” string to each client application that has previously executed LISTEN channel for the specified channel name in the current database. Notifications are visible to all users.

NOTIFY provides a simple interprocess communication mechanism for a collection of processes accessing the same PostgreSQL database. A payload string can be sent along with the notification, and higher-level mechanisms for passing structured data can be built by using tables in the database to pass additional data from notifier to listener(s).

The information passed to the client for a notification event includes the notification channel name, the notifying session's server process PID, and the payload string, which is an empty string if it has not been specified.

It is up to the database designer to define the channel names that will be used in a given database and what each one means. Commonly, the channel name is the same as the name of some table in the database, and the notify event essentially means, “I changed this table, take a look at it to see what's new”. But no such association is enforced by the NOTIFY and LISTEN commands. For example, a database designer could use several different channel names to signal different sorts of changes to a single table. Alternatively, the payload string could be used to differentiate various cases.

When NOTIFY is used to signal the occurrence of changes to a particular table, a useful programming technique is to put the NOTIFY in a statement trigger that is triggered by table updates. In this way, notification happens automatically when the table is changed, and the application programmer cannot accidentally forget to do it.

NOTIFY interacts with SQL transactions in some important ways. Firstly, if a NOTIFY is executed inside a transaction, the notify events are not delivered until and unless the transaction is committed. This is appropriate, since if the transaction is aborted, all the commands within it have had no effect, including NOTIFY. But it can be disconcerting if one is expecting the notification events to be delivered immediately. Secondly, if a listening session receives a notification signal while it is within a transaction, the notification event will not be delivered to its connected client until just after the transaction is completed (either committed or aborted). Again, the reasoning is that if a notification were delivered within a transaction that was later aborted, one would want the notification to be undone somehow — but the server cannot “take back” a notification once it has sent it to the client. So notification events are only delivered between transactions. The upshot of this is that applications using NOTIFY for real-time signaling should try to keep their transactions short.

If the same channel name is signaled multiple times from the same transaction with identical payload strings, the database server can decide to deliver a single notification only. On the other hand, notifications with distinct payload strings will always be delivered as distinct notifications. Similarly, notifications from different transactions will never get folded into one notification. Except for dropping later instances of duplicate notifications, NOTIFY guarantees that notifications from the same transaction get delivered in the order they were sent. It is also guaranteed that messages from different transactions are delivered in the order in which the transactions committed.

It is common for a client that executes NOTIFY to be listening on the same notification channel itself. In that case it will get back a notification event, just like all the other listening sessions. Depending on the application logic, this could result in useless work, for example, reading a database table to find the same updates that that session just wrote out. It is possible to avoid such extra work by noticing whether the notifying session's server process PID (supplied in the notification event message) is the same as one's own session's PID (available from libpq). When they are the same, the notification event is one's own work bouncing back, and can be ignored.

* https://www.postgresql.org/docs/current/sql-notify.html

## Integration with PostgreSQL

To support event processing in PostgreSQL you have to define manually what exactly you want to subscribe.
This possibility provides PostgreSQL triggers.

### Trigger function

```sql
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE
        data json;
        notification json;

    BEGIN

        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
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
        PERFORM pg_notify('events', notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;

$$ LANGUAGE plpgsql;
```

### Apply the trigger

```sql
CREATE TRIGGER products_notify_event
AFTER INSERT OR UPDATE OR DELETE ON products
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
```

## Using example

```go
import (
  nc "github.com/geniusrabbit/notificationcenter"
  "github.com/geniusrabbit/notificationcenter/pg"
)

const connectionDNS = "postgres://connection"

func main() {
  db, err := sql.Open("postgres", connectionDNS)
  if err != nil {
    t.Error(err)
    return
  }

  events := pg.MustNewSubscriber(connectionDNS, "events", nil)
  nc.Register("events", events)

  // Add new handler to process the stream "events"
  nc.Subscribe("events", nc.FuncHandler(func(msg nc.Message) error {
    fmt.Printf("%v\n", msg.Data())
    return nil
  }))

  // Run programm executers
  // ...

  // Run subscriber listeners
  nc.Listen()
}
```

## Dependencies

* github.com/lib/pq
