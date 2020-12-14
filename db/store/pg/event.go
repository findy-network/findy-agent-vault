package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

func sqlEventSelectBatchFor(where, limitArg string) string {
	return sqlEventSelect + " event " + where + " " + limitArg
}

const (
	sqlEventBatchWhere           = " WHERE tenant_id=$1 "
	sqlEventBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2"

	sqlEventFields = "tenant_id, connection_id, job_id, description, read"
	sqlEventInsert = "INSERT INTO event " + "(" + sqlEventFields + ") " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id, created, cursor"
	sqlEventSelect = "SELECT id, " + sqlEventFields + ", created, cursor FROM"
)

func (p *Database) AddEvent(e *model.Event) (n *model.Event, err error) {
	defer returnErr("AddEvent", &err)

	rows, err := p.db.Query(
		sqlEventInsert,
		e.TenantID,
		e.ConnectionID,
		e.JobID,
		e.Description,
		e.Read,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewEvent(e)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return n, err
}

func (p *Database) MarkEventRead(id, tenantID string) (e *model.Event, err error) {
	defer returnErr("MarkEventRead", &err)

	const sqlEventUpdate = "UPDATE event SET read=true WHERE id = $1 AND tenant_id = $2" +
		" RETURNING id," + sqlEventFields + ", created, cursor"

	rows, err := p.db.Query(
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
	n := model.NewEvent(nil)

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

func (p *Database) GetEvent(id, tenantID string) (e *model.Event, err error) {
	defer returnErr("GetEvent", &err)

	const sqlEventSelectByID = sqlEventSelect + " event" +
		" WHERE event.id=$1 AND tenant_id=$2"

	rows, err := p.db.Query(sqlEventSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	e = model.NewEvent(nil)
	if rows.Next() {
		e, err = readRowToEvent(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (p *Database) getEventsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	initialArgs []interface{},
) (e *model.Events, err error) {
	defer returnErr("GetEvents", &err)

	query, args := getBatchQuery(queries, batch, initialArgs)
	rows, err := p.db.Query(query, args...)
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

func (p *Database) GetEvents(info *paginator.BatchInfo, tenantID string) (c *model.Events, err error) {
	return p.getEventsForQuery(&queryInfo{
		Asc:        sqlEventSelectBatchFor(sqlEventBatchWhere+sqlOrderByAsc(""), "$2"),
		Desc:       sqlEventSelectBatchFor(sqlEventBatchWhere+sqlOrderByDesc(""), "$2"),
		AfterAsc:   sqlEventSelectBatchFor(sqlEventBatchWhere+" AND cursor > $2"+sqlOrderByAsc(""), "$3"),
		AfterDesc:  sqlEventSelectBatchFor(sqlEventBatchWhere+" AND cursor > $2"+sqlOrderByDesc(""), "$3"),
		BeforeAsc:  sqlEventSelectBatchFor(sqlEventBatchWhere+" AND cursor < $2"+sqlOrderByAsc(""), "$3"),
		BeforeDesc: sqlEventSelectBatchFor(sqlEventBatchWhere+" AND cursor < $2"+sqlOrderByDesc(""), "$3"),
	},
		info,
		[]interface{}{tenantID},
	)
}

func (p *Database) GetEventCount(tenantID string) (count int, err error) {
	defer returnErr("GetEventCount", &err)

	const sqlEventSelectCount = "SELECT count(id) FROM event " + sqlEventBatchWhere

	rows, err := p.db.Query(sqlEventSelectCount, tenantID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (p *Database) GetConnectionEvents(
	info *paginator.BatchInfo,
	tenantID,
	connectionID string,
) (connections *model.Events, err error) {
	return p.getEventsForQuery(&queryInfo{
		Asc:        sqlEventSelectBatchFor(sqlEventBatchWhereConnection+sqlOrderByAsc(""), "$3"),
		Desc:       sqlEventSelectBatchFor(sqlEventBatchWhereConnection+sqlOrderByDesc(""), "$3"),
		AfterAsc:   sqlEventSelectBatchFor(sqlEventBatchWhereConnection+" AND cursor > $3"+sqlOrderByAsc(""), "$4"),
		AfterDesc:  sqlEventSelectBatchFor(sqlEventBatchWhereConnection+" AND cursor > $3"+sqlOrderByDesc(""), "$4"),
		BeforeAsc:  sqlEventSelectBatchFor(sqlEventBatchWhereConnection+" AND cursor < $3"+sqlOrderByAsc(""), "$4"),
		BeforeDesc: sqlEventSelectBatchFor(sqlEventBatchWhereConnection+" AND cursor < $3"+sqlOrderByDesc(""), "$4"),
	},
		info,
		[]interface{}{tenantID, connectionID},
	)
}

func (p *Database) GetConnectionEventCount(tenantID, connectionID string) (count int, err error) {
	defer returnErr("GetEventCount", &err)

	const sqlEventSelectCount = "SELECT count(id) FROM event " + sqlEventBatchWhereConnection

	rows, err := p.db.Query(sqlEventSelectCount, tenantID, connectionID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}
