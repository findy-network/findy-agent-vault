package pg

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

var (
	messageFields = []string{"tenant_id", "connection_id", "message", "sent_by_me", "delivered"}

	sqlBaseMessageFields = sqlFields("", messageFields)
	sqlMessageInsert     = "INSERT INTO message " + "(" + sqlBaseMessageFields + ") " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id, created, cursor"
	sqlMessageSelect = "SELECT id, " + sqlBaseMessageFields + ", created, cursor FROM"
)

func (pg *Database) getMessageForObject(objectName, columnName, objectID, tenantID string) (c *model.Message, err error) {
	defer returnErr("getMessageForObject", &err)

	sqlMessageSelectByObjectID := "SELECT message.id, " +
		sqlFields("message", messageFields) + ", message.created, message.cursor FROM" +
		" message INNER JOIN " + objectName + " ON " + objectName +
		"." + columnName + "=message.id WHERE " + objectName + ".id = $1 AND message.tenant_id = $2"

	rows, err := pg.db.Query(sqlMessageSelectByObjectID, objectID, tenantID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		c, err = readRowToMessage(rows)
	} else {
		err = fmt.Errorf("not found message for %s id %s", objectName, objectID)
	}
	err2.Check(err)

	return
}

func readRowToMessage(rows *sql.Rows) (*model.Message, error) {
	n := model.NewMessage("", nil)

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

	n = model.NewMessage(arg.TenantID, arg)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
	} else {
		err = fmt.Errorf("no rows returned from insert message query")
	}
	err2.Check(err)

	return n, err
}

func (pg *Database) UpdateMessage(arg *model.Message) (m *model.Message, err error) {
	defer returnErr("UpdateMessage", &err)

	sqlMessageUpdate := "UPDATE message SET delivered=$1 WHERE id = $2 AND tenant_id = $3" +
		" RETURNING id," + sqlBaseMessageFields + ", created, cursor"

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
	} else {
		err = fmt.Errorf("no rows returned from update message query")
	}
	err2.Check(err)

	return m, err
}

func (pg *Database) GetMessage(id, tenantID string) (m *model.Message, err error) {
	defer returnErr("GetMessage", &err)

	sqlMessageSelectByID := sqlMessageSelect + " message WHERE id=$1 AND tenant_id=$2"

	rows, err := pg.db.Query(sqlMessageSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	m = model.NewMessage("", nil)
	if rows.Next() {
		m, err = readRowToMessage(rows)
	} else {
		err = fmt.Errorf("no rows returned from select message query (%s)", id)
	}
	err2.Check(err)

	return
}

func (pg *Database) getMessagesForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,
	initialArgs []interface{},
) (m *model.Messages, err error) {
	defer returnErr("GetMessages", &err)

	query, args := getBatchQuery(queries, batch, tenantID, initialArgs)
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

func sqlMessageBatchWhere(cursorParam, connectionParam, limitParam string, desc, before bool) string {
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
	return sqlMessageSelect + " message " + where + cursorOrder + " " + limitParam
}

func (pg *Database) GetMessages(info *paginator.BatchInfo, tenantID string, connectionID *string) (m *model.Messages, err error) {
	if connectionID == nil {
		return pg.getMessagesForQuery(&queryInfo{
			Asc:        sqlMessageBatchWhere("", "", "$2", false, false),
			Desc:       sqlMessageBatchWhere("", "", "$2", true, false),
			AfterAsc:   sqlMessageBatchWhere("$2", "", "$3", false, false),
			AfterDesc:  sqlMessageBatchWhere("$2", "", "$3", true, false),
			BeforeAsc:  sqlMessageBatchWhere("$2", "", "$3", false, true),
			BeforeDesc: sqlMessageBatchWhere("$2", "", "$3", true, true),
		},
			info,
			tenantID,
			[]interface{}{},
		)
	}
	return pg.getMessagesForQuery(&queryInfo{
		Asc:        sqlMessageBatchWhere("", "$2", "$3", false, false),
		Desc:       sqlMessageBatchWhere("", "$2", "$3", true, false),
		AfterAsc:   sqlMessageBatchWhere("$2", "$3", "$4", false, false),
		AfterDesc:  sqlMessageBatchWhere("$2", "$3", "$4", true, false),
		BeforeAsc:  sqlMessageBatchWhere("$2", "$3", "$4", false, true),
		BeforeDesc: sqlMessageBatchWhere("$2", "$3", "$4", true, true),
	},
		info,
		tenantID,
		[]interface{}{*connectionID},
	)
}

func (pg *Database) GetMessageCount(tenantID string, connectionID *string) (count int, err error) {
	defer returnErr("GetMessageCount", &err)
	const (
		sqlMessageBatchWhere           = " WHERE tenant_id=$1 "
		sqlMessageBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2"
	)
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
	return pg.getConnectionForObject("message", "connection_id", id, tenantID)
}

func (pg *Database) ArchiveMessage(m *model.Message) (*model.Message, error) {
	return nil, nil
}
