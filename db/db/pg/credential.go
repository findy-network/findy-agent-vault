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

var (
	sqlCredentialAttributeInsert = "INSERT INTO credential_attribute (credential_id, name, value, index) VALUES "

	sqlCredentialFields = "tenant_id, connection_id, role, schema_id, cred_def_id, initiated_by_us"
	sqlCredentialInsert = "INSERT INTO credential " + "(" + sqlCredentialFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, cursor"
	sqlCredentialSelect = "SELECT credential.id, " + sqlCredentialFields +
		", created, approved, issued, failed, cursor, credential_attribute.id, name, value FROM"
	sqlCredentialJoin       = " INNER JOIN credential_attribute on credential_id = credential.id"
	sqlCredentialSelectByID = sqlCredentialSelect + " credential" + sqlCredentialJoin +
		" WHERE credential.id=$1 AND tenant_id=$2" +
		" ORDER BY credential_attribute.index"
	sqlCredentialSelectBatch           = sqlCredentialSelectBatchFor(sqlWhereTenantAsc(""), "$2", "ASC")
	sqlCredentialSelectBatchTail       = sqlCredentialSelectBatchFor(sqlWhereTenantDesc(""), "$2", "DESC")
	sqlCredentialSelectBatchAfter      = sqlCredentialSelectBatchFor(sqlWhereTenantAscAfter(""), "$3", "ASC")
	sqlCredentialSelectBatchAfterTail  = sqlCredentialSelectBatchFor(sqlWhereTenantDescAfter(""), "$3", "DESC")
	sqlCredentialSelectBatchBefore     = sqlCredentialSelectBatchFor(sqlWhereTenantAscBefore(""), "$3", "ASC")
	sqlCredentialSelectBatchBeforeTail = sqlCredentialSelectBatchFor(sqlWhereTenantDescBefore(""), "$3", "DESC")
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

	n = c.Copy()
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

func readRowToCredential(rows *sql.Rows, previous *model.Credential) (*model.Credential, error) {
	a := &graph.CredentialValue{}
	var approved sql.NullTime
	var issued sql.NullTime
	var failed sql.NullTime

	n := model.NewCredential()

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

	rows, err := p.db.Query(sqlCredentialSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	c = model.NewCredential()
	for rows.Next() {
		c, err = readRowToCredential(rows, c)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (p *Database) GetCredentials(info *paginator.BatchInfo, tenantID string) (c *model.Credentials, err error) {
	defer returnErr("GetCredentials", &err)

	query := ""
	args := make([]interface{}, 0)
	args = append(args, tenantID)
	if info.Tail {
		query = sqlCredentialSelectBatchTail
		if info.After > 0 {
			query = sqlCredentialSelectBatchAfterTail
		} else if info.Before > 0 {
			query = sqlCredentialSelectBatchBeforeTail
		}
	} else {
		query = sqlCredentialSelectBatch
		if info.After > 0 {
			query = sqlCredentialSelectBatchAfter
		} else if info.Before > 0 {
			query = sqlCredentialSelectBatchBefore
		}
	}
	if info.After > 0 {
		args = append(args, info.After)
	} else if info.Before > 0 {
		args = append(args, info.Before)
	}

	args = append(args, info.Count+1)

	rows, err := p.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	c = &model.Credentials{
		Credentials:     make([]*model.Credential, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	prevCredential := model.NewCredential()
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

	if info.Count < len(c.Credentials) {
		c.Credentials = c.Credentials[:info.Count]
		if info.Tail {
			c.HasPreviousPage = true
		} else {
			c.HasNextPage = true
		}
	}

	if info.After > 0 {
		c.HasPreviousPage = true
	}
	if info.Before > 0 {
		c.HasNextPage = true
	}

	// Reverse order for tail first
	if info.Tail {
		sort.Slice(c.Credentials, func(i, j int) bool {
			return c.Credentials[i].Created.Sub(c.Credentials[j].Created) < 0
		})
	}

	return c, err
}
