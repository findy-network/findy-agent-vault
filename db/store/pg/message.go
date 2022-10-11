package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

var (
	messageFields = []string{"tenant_id", "connection_id", "message", "sent_by_me", "delivered", "archived"}

	sqlBaseMessageFields = sqlFields("", messageFields)
	sqlMessageSelect     = "SELECT id, " + sqlBaseMessageFields + ", created, cursor FROM"
)

func (pg *Database) getMessageForObject(objectName, columnName, objectID, tenantID string) (m *model.Message, err error) {
	defer err2.Returnf(&err, "getMessageForObject")

	sqlMessageSelectByObjectID := "SELECT message.id, " +
		sqlFields("message", messageFields) + ", message.created, message.cursor FROM" +
		" message INNER JOIN " + objectName + " ON " + objectName +
		"." + columnName + "=message.id WHERE " + objectName + ".id = $1 AND message.tenant_id = $2"

	m = &model.Message{}
	try.To(pg.doRowQuery(
		readRowToMessage(m),
		sqlMessageSelectByObjectID,
		objectID,
		tenantID,
	))

	return
}

func rowToMessage(rows *sql.Rows) (n *model.Message, err error) {
	n = &model.Message{}
	return n, readRowToMessage(n)(rows)
}

func readRowToMessage(n *model.Message) func(*sql.Rows) error {
	return func(rows *sql.Rows) error {
		return rows.Scan(
			&n.ID,
			&n.TenantID,
			&n.ConnectionID,
			&n.Message,
			&n.SentByMe,
			&n.Delivered,
			&n.Archived,
			&n.Created,
			&n.Cursor,
		)
	}
}

func (pg *Database) AddMessage(arg *model.Message) (msg *model.Message, err error) {
	defer err2.Returnf(&err, "AddMessage")

	var (
		sqlMessageInsert = "INSERT INTO message " + "(" + sqlBaseMessageFields + ") " +
			"VALUES (" + sqlArguments(messageFields) + ") RETURNING " + sqlInsertFields
	)

	msg = &model.Message{}
	*msg = *arg
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&msg.ID, &msg.Created, &msg.Cursor)
		},
		sqlMessageInsert,
		arg.TenantID,
		arg.ConnectionID,
		arg.Message,
		arg.SentByMe,
		arg.Delivered,
		arg.Archived,
	))

	return msg, err
}

func (pg *Database) UpdateMessage(arg *model.Message) (m *model.Message, err error) {
	defer err2.Returnf(&err, "UpdateMessage")

	sqlMessageUpdate := "UPDATE message SET delivered=$1 WHERE id = $2 AND tenant_id = $3" +
		" RETURNING id," + sqlBaseMessageFields + ", created, cursor"

	m = &model.Message{}
	try.To(pg.doRowQuery(
		readRowToMessage(m),
		sqlMessageUpdate,
		arg.Delivered,
		arg.ID,
		arg.TenantID,
	))
	return m, err
}

func (pg *Database) GetMessage(id, tenantID string) (m *model.Message, err error) {
	defer err2.Returnf(&err, "GetMessage")

	sqlMessageSelectByID := sqlMessageSelect + " message WHERE id=$1 AND tenant_id=$2"

	m = &model.Message{}
	try.To(pg.doRowQuery(
		readRowToMessage(m),
		sqlMessageSelectByID,
		id,
		tenantID,
	))

	return
}

func (pg *Database) getMessagesForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,
	initialArgs []interface{},
) (m *model.Messages, err error) {
	defer err2.Returnf(&err, "GetMessages")

	query, args := getBatchQuery(queries, batch, tenantID, initialArgs)

	m = &model.Messages{
		Messages:        make([]*model.Message, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var message *model.Message
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Return(&err)
		message = try.To1(rowToMessage(rows))
		m.Messages = append(m.Messages, message)
		return
	}, query, args...))

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
	defer err2.Returnf(&err, "GetMessageCount")
	const (
		sqlMessageBatchWhere           = " WHERE tenant_id=$1 "
		sqlMessageBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2"
	)
	count = try.To1(pg.getCount(
		"message",
		sqlMessageBatchWhere,
		sqlMessageBatchWhereConnection,
		tenantID,
		connectionID,
	))
	return
}

func (pg *Database) GetConnectionForMessage(id, tenantID string) (*model.Connection, error) {
	return pg.getConnectionForObject("message", "connection_id", id, tenantID)
}

func (pg *Database) ArchiveMessage(id, tenantID string) (err error) {
	defer err2.Returnf(&err, "ArchiveMessage")

	var (
		sqlMessageArchive = "UPDATE message SET archived=$1 WHERE id = $2 and tenant_id = $3 RETURNING id"
	)

	now := utils.CurrentTime()
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&id)
		},
		sqlMessageArchive,
		now,
		id,
		tenantID,
	))
	return
}
