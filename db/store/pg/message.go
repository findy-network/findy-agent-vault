package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

func sqlMessageSelectBatchFor(where, limitArg string) string {
	return sqlMessageSelect + " message " + where + " " + limitArg
}

const (
	sqlMessageBatchWhere           = " WHERE tenant_id=$1 "
	sqlMessageBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2"

	sqlMessageFields = "tenant_id, connection_id, message, sent_by_me, delivered"
	sqlMessageInsert = "INSERT INTO message " + "(" + sqlMessageFields + ") " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id, created, cursor"
	sqlMessageSelect = "SELECT id, " + sqlMessageFields + ", created, cursor FROM"
)

func readRowToMessage(rows *sql.Rows) (*model.Message, error) {
	n := model.NewMessage(nil)

	err := rows.Scan(
		&n.ID,
		&n.TenantID,
		&n.ConnectionID,
		&n.Message,
		&n.SentByMe,
		&n.Delivered,
		&n.Created,
		&n.Cursor,
	)
	return n, err
}

func (pg *Database) AddMessage(arg *model.Message) (n *model.Message, err error) {
	defer returnErr("AddMessage", &err)

	rows, err := pg.db.Query(
		sqlMessageInsert,
		arg.TenantID,
		arg.ConnectionID,
		arg.Message,
		arg.SentByMe,
		arg.Delivered,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewMessage(arg)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return n, err
}

func (pg *Database) UpdateMessage(arg *model.Message) (m *model.Message, err error) {
	defer returnErr("UpdateMessage", &err)

	const sqlMessageUpdate = "UPDATE message SET delivered=$1 WHERE id = $2 AND tenant_id = $3" +
		" RETURNING id," + sqlMessageFields + ", created, cursor"

	rows, err := pg.db.Query(
		sqlMessageUpdate,
		arg.Delivered,
		arg.ID,
		arg.TenantID,
	)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		m, err = readRowToMessage(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return m, err
}

func (pg *Database) GetMessage(id, tenantID string) (m *model.Message, err error) {
	defer returnErr("GetMessage", &err)

	const sqlMessageSelectByID = sqlMessageSelect + " message" +
		" WHERE id=$1 AND tenant_id=$2"

	rows, err := pg.db.Query(sqlMessageSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	m = model.NewMessage(nil)
	if rows.Next() {
		m, err = readRowToMessage(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (pg *Database) getMessagesForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	initialArgs []interface{},
) (m *model.Messages, err error) {
	defer returnErr("GetMessages", &err)

	query, args := getBatchQuery(queries, batch, initialArgs)
	rows, err := pg.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	m = &model.Messages{
		Messages:        make([]*model.Message, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var message *model.Message
	for rows.Next() {
		message, err = readRowToMessage(rows)
		err2.Check(err)
		m.Messages = append(m.Messages, message)
	}

	err = rows.Err()
	err2.Check(err)

	if batch.Count < len(m.Messages) {
		m.Messages = m.Messages[:batch.Count]
		if batch.Tail {
			m.HasPreviousPage = true
		} else {
			m.HasNextPage = true
		}
	}

	if batch.After > 0 {
		m.HasPreviousPage = true
	}
	if batch.Before > 0 {
		m.HasNextPage = true
	}

	// Reverse order for tail first
	if batch.Tail {
		sort.Slice(m.Messages, func(i, j int) bool {
			return m.Messages[i].Created.Sub(m.Messages[j].Created) < 0
		})
	}

	return m, err
}

func (pg *Database) GetMessages(info *paginator.BatchInfo, tenantID string, connectionID *string) (m *model.Messages, err error) {
	if connectionID == nil {
		return pg.getMessagesForQuery(&queryInfo{
			Asc:        sqlMessageSelectBatchFor(sqlMessageBatchWhere+sqlOrderByAsc(""), "$2"),
			Desc:       sqlMessageSelectBatchFor(sqlMessageBatchWhere+sqlOrderByDesc(""), "$2"),
			AfterAsc:   sqlMessageSelectBatchFor(sqlMessageBatchWhere+" AND cursor > $2"+sqlOrderByAsc(""), "$3"),
			AfterDesc:  sqlMessageSelectBatchFor(sqlMessageBatchWhere+" AND cursor > $2"+sqlOrderByDesc(""), "$3"),
			BeforeAsc:  sqlMessageSelectBatchFor(sqlMessageBatchWhere+" AND cursor < $2"+sqlOrderByAsc(""), "$3"),
			BeforeDesc: sqlMessageSelectBatchFor(sqlMessageBatchWhere+" AND cursor < $2"+sqlOrderByDesc(""), "$3"),
		},
			info,
			[]interface{}{tenantID},
		)
	}
	return pg.getMessagesForQuery(&queryInfo{
		Asc:        sqlMessageSelectBatchFor(sqlMessageBatchWhereConnection+sqlOrderByAsc(""), "$3"),
		Desc:       sqlMessageSelectBatchFor(sqlMessageBatchWhereConnection+sqlOrderByDesc(""), "$3"),
		AfterAsc:   sqlMessageSelectBatchFor(sqlMessageBatchWhereConnection+" AND cursor > $3"+sqlOrderByAsc(""), "$4"),
		AfterDesc:  sqlMessageSelectBatchFor(sqlMessageBatchWhereConnection+" AND cursor > $3"+sqlOrderByDesc(""), "$4"),
		BeforeAsc:  sqlMessageSelectBatchFor(sqlMessageBatchWhereConnection+" AND cursor < $3"+sqlOrderByAsc(""), "$4"),
		BeforeDesc: sqlMessageSelectBatchFor(sqlMessageBatchWhereConnection+" AND cursor < $3"+sqlOrderByDesc(""), "$4"),
	},
		info,
		[]interface{}{tenantID, *connectionID},
	)
}

func (pg *Database) GetMessageCount(tenantID string, connectionID *string) (count int, err error) {
	defer returnErr("GetMessageCount", &err)
	count, err = pg.getCount(
		"message",
		sqlMessageBatchWhere,
		sqlMessageBatchWhereConnection,
		tenantID,
		connectionID,
	)
	err2.Check(err)
	return
}

func (pg *Database) GetConnectionForMessage(id, tenantID string) (*model.Connection, error) {
	return pg.getConnectionForObject("message", id, tenantID)
}
