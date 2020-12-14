package pg

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

func constructCredentialAttributeInsert(count int) string {
	const sqlCredentialAttributeInsert = "INSERT INTO credential_attribute (credential_id, name, value, index) VALUES "

	result := sqlCredentialAttributeInsert
	paramCount := 4
	for i := 0; i < count; i++ {
		if i >= 1 {
			result += ","
		}
		nbr := i*paramCount + 1
		params := ""
		for j := 0; j < paramCount; j++ {
			if j >= 1 {
				params += ","
			}
			params = fmt.Sprintf("%s$%d", params, (nbr + j))
		}
		result += fmt.Sprintf("(%s)", params)
	}
	return result + " RETURNING id"
}

func sqlCredentialSelectBatchFor(tenantOrder, limit, cursorOrder string) string {
	return sqlCredentialSelect +
		" (SELECT * FROM credential " + tenantOrder + " " + limit + ") AS credential " +
		sqlCredentialJoin + " ORDER BY cursor " + cursorOrder + ", credential_attribute.index"
}

const (
	sqlCredentialBatchWhere           = " WHERE tenant_id=$1 AND issued IS NOT NULL "
	sqlCredentialBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2 AND issued IS NOT NULL "

	sqlCredentialFields = "tenant_id, connection_id, role, schema_id, cred_def_id, initiated_by_us"
	sqlCredentialInsert = "INSERT INTO credential " + "(" + sqlCredentialFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, cursor"
	sqlCredentialSelect = "SELECT credential.id, " + sqlCredentialFields +
		", created, approved, issued, failed, cursor, credential_attribute.id, name, value FROM"
	sqlCredentialJoin = " INNER JOIN credential_attribute on credential_id = credential.id"
)

func (p *Database) addCredentialAttributes(id string, attributes []*graph.CredentialValue) (a []*graph.CredentialValue, err error) {
	defer returnErr("addCredentialAttributes", &err)

	query := constructCredentialAttributeInsert(len(attributes))
	args := make([]interface{}, 0)
	for index, a := range attributes {
		args = append(args, []interface{}{id, a.Name, a.Value, index}...)
	}

	rows, err := p.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	index := 0
	for rows.Next() {
		err = rows.Scan(&attributes[index].ID)
		err2.Check(err)
		index++
	}

	err = rows.Err()
	err2.Check(err)

	return attributes, nil
}

func (p *Database) AddCredential(c *model.Credential) (n *model.Credential, err error) {
	defer returnErr("AddCredential", &err)

	if len(c.Attributes) == 0 {
		panic("Attributes are always required for credential.")
	}
	rows, err := p.db.Query(
		sqlCredentialInsert,
		c.TenantID,
		c.ConnectionID,
		c.Role,
		c.SchemaID,
		c.CredDefID,
		c.InitiatedByUs,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewCredential(c)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	attributes, err := p.addCredentialAttributes(n.ID, n.Attributes)
	err2.Check(err)

	n.Attributes = attributes
	return n, err
}

func (p *Database) UpdateCredential(c *model.Credential) (n *model.Credential, err error) {
	defer returnErr("UpdateCredential", &err)

	const sqlCredentialUpdate = "UPDATE credential SET approved=$1, issued=$2, failed=$3 WHERE id = $4" // TODO: tenant_id, connection_id?

	_, err = p.db.Exec(
		sqlCredentialUpdate,
		c.Approved,
		c.Issued,
		c.Failed,
		c.ID,
	)
	err2.Check(err)
	return c, err
}

func readRowToCredential(rows *sql.Rows, previous *model.Credential) (*model.Credential, error) {
	a := &graph.CredentialValue{}
	var approved sql.NullTime
	var issued sql.NullTime
	var failed sql.NullTime

	n := model.NewCredential(nil)

	err := rows.Scan(
		&n.ID,
		&n.TenantID,
		&n.ConnectionID,
		&n.Role,
		&n.SchemaID,
		&n.CredDefID,
		&n.InitiatedByUs,
		&n.Created,
		&approved,
		&issued,
		&failed,
		&n.Cursor,
		&a.ID,
		&a.Name,
		&a.Value,
	)

	if approved.Valid {
		n.Approved = &approved.Time
	}
	if issued.Valid {
		n.Issued = &issued.Time
	}
	if failed.Valid {
		n.Failed = &failed.Time
	}

	n.Attributes = make([]*graph.CredentialValue, 0)
	if previous.ID == n.ID {
		n.Attributes = append(n.Attributes, previous.Attributes...)
		n.Attributes = append(n.Attributes, a)
	} else {
		n.Attributes = append(n.Attributes, a)
	}

	return n, err
}

func (p *Database) GetCredential(id, tenantID string) (c *model.Credential, err error) {
	defer returnErr("GetCredential", &err)

	const sqlCredentialSelectByID = sqlCredentialSelect + " credential" + sqlCredentialJoin +
		" WHERE credential.id=$1 AND tenant_id=$2" +
		" ORDER BY credential_attribute.index"

	rows, err := p.db.Query(sqlCredentialSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	c = model.NewCredential(nil)
	for rows.Next() {
		c, err = readRowToCredential(rows, c)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (p *Database) getCredentialsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	initialArgs []interface{},
) (c *model.Credentials, err error) {
	defer returnErr("GetCredentials", &err)

	query, args := getBatchQuery(queries, batch, initialArgs)
	rows, err := p.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	c = &model.Credentials{
		Credentials:     make([]*model.Credential, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	prevCredential := model.NewCredential(nil)
	var credential *model.Credential
	for rows.Next() {
		credential, err = readRowToCredential(rows, prevCredential)
		err2.Check(err)
		if prevCredential.ID != "" && prevCredential.ID != credential.ID {
			c.Credentials = append(c.Credentials, prevCredential)
		}
		prevCredential = credential
	}

	// ensure also last credential is added
	if prevCredential.ID != c.Credentials[len(c.Credentials)-1].ID {
		c.Credentials = append(c.Credentials, prevCredential)
	}

	err = rows.Err()
	err2.Check(err)

	if batch.Count < len(c.Credentials) {
		c.Credentials = c.Credentials[:batch.Count]
		if batch.Tail {
			c.HasPreviousPage = true
		} else {
			c.HasNextPage = true
		}
	}

	if batch.After > 0 {
		c.HasPreviousPage = true
	}
	if batch.Before > 0 {
		c.HasNextPage = true
	}

	// Reverse order for tail first
	if batch.Tail {
		sort.Slice(c.Credentials, func(i, j int) bool {
			return c.Credentials[i].Created.Sub(c.Credentials[j].Created) < 0
		})
	}

	return c, err
}

func (p *Database) GetCredentials(info *paginator.BatchInfo, tenantID string) (c *model.Credentials, err error) {
	return p.getCredentialsForQuery(&queryInfo{
		Asc:        sqlCredentialSelectBatchFor(sqlCredentialBatchWhere+sqlOrderByAsc(""), "$2", "ASC"),
		Desc:       sqlCredentialSelectBatchFor(sqlCredentialBatchWhere+sqlOrderByDesc(""), "$2", "DESC"),
		AfterAsc:   sqlCredentialSelectBatchFor(sqlCredentialBatchWhere+" AND cursor > $2"+sqlOrderByAsc(""), "$3", "ASC"),
		AfterDesc:  sqlCredentialSelectBatchFor(sqlCredentialBatchWhere+" AND cursor > $2"+sqlOrderByDesc(""), "$3", "DESC"),
		BeforeAsc:  sqlCredentialSelectBatchFor(sqlCredentialBatchWhere+" AND cursor < $2"+sqlOrderByAsc(""), "$3", "ASC"),
		BeforeDesc: sqlCredentialSelectBatchFor(sqlCredentialBatchWhere+" AND cursor < $2"+sqlOrderByDesc(""), "$3", "DESC"),
	},
		info,
		[]interface{}{tenantID},
	)
}

func (p *Database) GetCredentialCount(tenantID string) (count int, err error) {
	defer returnErr("GetCredentialCount", &err)

	const sqlCredentialSelectCount = "SELECT count(id) FROM credential " + sqlCredentialBatchWhere

	rows, err := p.db.Query(sqlCredentialSelectCount, tenantID)
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

func (p *Database) GetConnectionCredentials(
	info *paginator.BatchInfo,
	tenantID,
	connectionID string,
) (connections *model.Credentials, err error) {
	return p.getCredentialsForQuery(&queryInfo{
		Asc:        sqlCredentialSelectBatchFor(sqlCredentialBatchWhereConnection+sqlOrderByAsc(""), "$3", "ASC"),
		Desc:       sqlCredentialSelectBatchFor(sqlCredentialBatchWhereConnection+sqlOrderByDesc(""), "$3", "DESC"),
		AfterAsc:   sqlCredentialSelectBatchFor(sqlCredentialBatchWhereConnection+" AND cursor > $3"+sqlOrderByAsc(""), "$4", "ASC"),
		AfterDesc:  sqlCredentialSelectBatchFor(sqlCredentialBatchWhereConnection+" AND cursor > $3"+sqlOrderByDesc(""), "$4", "DESC"),
		BeforeAsc:  sqlCredentialSelectBatchFor(sqlCredentialBatchWhereConnection+" AND cursor < $3"+sqlOrderByAsc(""), "$4", "ASC"),
		BeforeDesc: sqlCredentialSelectBatchFor(sqlCredentialBatchWhereConnection+" AND cursor < $3"+sqlOrderByDesc(""), "$4", "DESC"),
	},
		info,
		[]interface{}{tenantID, connectionID},
	)
}

func (p *Database) GetConnectionCredentialCount(tenantID, connectionID string) (count int, err error) {
	defer returnErr("GetCredentialCount", &err)

	const sqlCredentialSelectCount = "SELECT count(id) FROM credential " + sqlCredentialBatchWhereConnection

	rows, err := p.db.Query(sqlCredentialSelectCount, tenantID, connectionID)
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
