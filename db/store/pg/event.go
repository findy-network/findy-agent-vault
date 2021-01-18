package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

const (
	sqlEventFields = "tenant_id, connection_id, job_id, description, read"
	sqlEventInsert = "INSERT INTO event " + "(" + sqlEventFields + ") " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id, created, cursor"
	sqlEventSelect = "SELECT id, " + sqlEventFields + ", created, cursor FROM"
)

func (pg *Database) AddEvent(e *model.Event) (n *model.Event, err error) {
	defer returnErr("AddEvent", &err)

	rows, err := pg.db.Query(
		sqlEventInsert,
		e.TenantID,
		e.ConnectionID,
		e.JobID,
		e.Description,
		e.Read,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewEvent(e.TenantID, e)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return n, err
}

func (pg *Database) MarkEventRead(id, tenantID string) (e *model.Event, err error) {
	defer returnErr("MarkEventRead", &err)

	const sqlEventUpdate = "UPDATE event SET read=true WHERE id = $1 AND tenant_id = $2" +
		" RETURNING id," + sqlEventFields + ", created, cursor"

	rows, err := pg.db.Query(
		sqlEventUpdate,
		id,
		tenantID,
	)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		e, err = readRowToEvent(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return e, err
}

func readRowToEvent(rows *sql.Rows) (*model.Event, error) {
	n := model.NewEvent("", nil)

	err := rows.Scan(
		&n.ID,
		&n.TenantID,
		&n.ConnectionID,
		&n.JobID,
		&n.Description,
		&n.Read,
		&n.Created,
		&n.Cursor,
	)
	return n, err
}

func (pg *Database) GetEvent(id, tenantID string) (e *model.Event, err error) {
	defer returnErr("GetEvent", &err)

	const sqlEventSelectByID = sqlEventSelect + " event" +
		" WHERE event.id=$1 AND tenant_id=$2"

	rows, err := pg.db.Query(sqlEventSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		e, err = readRowToEvent(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (pg *Database) getEventsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,
	initialArgs []interface{},
) (e *model.Events, err error) {
	defer returnErr("GetEvents", &err)

	query, args := getBatchQuery(queries, batch, tenantID, initialArgs)
	rows, err := pg.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	e = &model.Events{
		Events:          make([]*model.Event, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var event *model.Event
	for rows.Next() {
		event, err = readRowToEvent(rows)
		err2.Check(err)
		e.Events = append(e.Events, event)
	}

	err = rows.Err()
	err2.Check(err)

	if batch.Count < len(e.Events) {
		e.Events = e.Events[:batch.Count]
		if batch.Tail {
			e.HasPreviousPage = true
		} else {
			e.HasNextPage = true
		}
	}

	if batch.After > 0 {
		e.HasPreviousPage = true
	}
	if batch.Before > 0 {
		e.HasNextPage = true
	}

	// Reverse order for tail first
	if batch.Tail {
		sort.Slice(e.Events, func(i, j int) bool {
			return e.Events[i].Created.Sub(e.Events[j].Created) < 0
		})
	}

	return e, err
}

func sqlEventBatchWhere(cursorParam, connectionParam, limitParam string, desc, before bool) string {
	const whereTenantID = " WHERE tenant_id=$1 "
	cursorOrder := sqlOrderByCursorAsc
	cursor := ""
	connection := ""
	compareChar := sqlGreaterThan
	if before {
		compareChar = sqlLessThan
	}
	if connectionParam != "" {
		connection = " AND connection_id = " + connectionParam + " "
	}
	if cursorParam != "" {
		cursor = " AND cursor " + compareChar + cursorParam + " "
		if desc {
			cursor = " AND cursor " + compareChar + cursorParam + " "
		}
	}
	if desc {
		cursorOrder = sqlOrderByCursorDesc
	}
	where := whereTenantID + cursor + connection
	return sqlEventSelect + " event " + where + cursorOrder + " " + limitParam
}

func (pg *Database) GetEvents(info *paginator.BatchInfo, tenantID string, connectionID *string) (c *model.Events, err error) {
	if connectionID == nil {
		return pg.getEventsForQuery(&queryInfo{
			Asc:        sqlEventBatchWhere("", "", "$2", false, false),
			Desc:       sqlEventBatchWhere("", "", "$2", true, false),
			AfterAsc:   sqlEventBatchWhere("$2", "", "$3", false, false),
			AfterDesc:  sqlEventBatchWhere("$2", "", "$3", true, false),
			BeforeAsc:  sqlEventBatchWhere("$2", "", "$3", false, true),
			BeforeDesc: sqlEventBatchWhere("$2", "", "$3", true, true),
		},
			info,
			tenantID,
			[]interface{}{},
		)
	}
	return pg.getEventsForQuery(&queryInfo{
		Asc:        sqlEventBatchWhere("", "$2", "$3", false, false),
		Desc:       sqlEventBatchWhere("", "$2", "$3", true, false),
		AfterAsc:   sqlEventBatchWhere("$2", "$3", "$4", false, false),
		AfterDesc:  sqlEventBatchWhere("$2", "$3", "$4", true, false),
		BeforeAsc:  sqlEventBatchWhere("$2", "$3", "$4", false, true),
		BeforeDesc: sqlEventBatchWhere("$2", "$3", "$4", true, true),
	},
		info,
		tenantID,
		[]interface{}{*connectionID},
	)
}

func (pg *Database) GetEventCount(tenantID string, connectionID *string) (count int, err error) {
	defer returnErr("GetEventCount", &err)
	const (
		sqlEventBatchWhere           = " WHERE tenant_id=$1 "
		sqlEventBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2"
	)
	count, err = pg.getCount(
		"event",
		sqlEventBatchWhere,
		sqlEventBatchWhereConnection,
		tenantID,
		connectionID,
	)
	err2.Check(err)
	return
}

func (pg *Database) GetConnectionForEvent(id, tenantID string) (*model.Connection, error) {
	return pg.getConnectionForObject("event", "connection_id", id, tenantID)
}

func (pg *Database) GetJobForEvent(id, tenantID string) (*model.Job, error) {
	return pg.getJobForObject("event", id, tenantID)
}
