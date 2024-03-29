//nolint:goconst
package pg

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/lib/pq"
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

var (
	credentialFields      = []string{"tenant_id", "connection_id", "role", "schema_id", "cred_def_id", "initiated_by_us", "archived"}
	credentialExtraFields = []string{"created", "approved", "issued", "failed", "cursor"}

	sqlCredentialBaseFields = sqlFields("", credentialFields)
	sqlCredentialSelect     = "SELECT credential.id, " + sqlCredentialBaseFields + "," + sqlFields("", credentialExtraFields) +
		", credential_attribute.id, name, value FROM"
)

const (
	//#nosec
	sqlCredentialJoin = " INNER JOIN credential_attribute on credential_attribute.credential_id = credential.id"
)

func (pg *Database) getCredentialForObject(objectName, columnName, objectID, tenantID string) (cred *model.Credential, err error) {
	defer err2.Handle(&err, "getCredentialForObject")

	sqlCredentialJoinSelect := "SELECT credential.id, " +
		sqlFields("credential", credentialFields) + "," + sqlFields("credential", credentialExtraFields) +
		", credential_attribute.id, credential_attribute.name, credential_attribute.value FROM"
	sqlCredentialSelectByObjectID := sqlCredentialJoinSelect + " credential " + sqlCredentialJoin +
		" INNER JOIN " + objectName + " ON " + objectName +
		"." + columnName + "=credential.id WHERE " + objectName + ".id = $1 AND credential.tenant_id = $2"

	cred = &model.Credential{}
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		cred = try.To1(readRowToCredential(rows, cred))
		return
	}, sqlCredentialSelectByObjectID, objectID, tenantID))

	return
}

func (pg *Database) addCredentialAttributes(id string, attributes []*graph.CredentialValue) (a []*graph.CredentialValue, err error) {
	defer err2.Handle(&err, "addCredentialAttributes")

	query := constructCredentialAttributeInsert(len(attributes))
	args := make([]interface{}, 0)
	for index, a := range attributes {
		args = append(args, []interface{}{id, a.Name, a.Value, index}...)
	}

	index := 0
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		try.To(rows.Scan(&attributes[index].ID))
		index++
		return
	}, query, args...))

	return attributes, nil
}

func (pg *Database) AddCredential(c *model.Credential) (cred *model.Credential, err error) {
	defer err2.Handle(&err, "AddCredential")

	sqlCredentialInsert := "INSERT INTO credential " + "(" + sqlCredentialBaseFields + ") " +
		"VALUES (" + sqlArguments(credentialFields) + ") RETURNING " + sqlInsertFields

	if len(c.Attributes) == 0 {
		panic("Attributes are always required for credential.")
	}

	cred = &model.Credential{}
	*cred = *c
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&cred.ID, &cred.Created, &cred.Cursor)
		},
		sqlCredentialInsert,
		c.TenantID,
		c.ConnectionID,
		c.Role,
		c.SchemaID,
		c.CredDefID,
		c.InitiatedByUs,
		c.Archived,
	))

	attributes := try.To1(pg.addCredentialAttributes(cred.ID, cred.Attributes))

	cred.Attributes = attributes
	return cred, err
}

func (pg *Database) UpdateCredential(c *model.Credential) (n *model.Credential, err error) {
	defer err2.Handle(&err, "UpdateCredential")

	//#nosec
	const sqlCredentialUpdate = "UPDATE credential SET approved=$1, issued=$2, failed=$3 WHERE id = $4" // TODO: tenant_id, connection_id?

	try.To1(pg.db.Exec(
		sqlCredentialUpdate,
		c.Approved,
		c.Issued,
		c.Failed,
		c.ID,
	))
	return c, err
}

func readRowToCredential(rows *sql.Rows, previous *model.Credential) (*model.Credential, error) {
	a := &graph.CredentialValue{}

	cred := &model.Credential{}

	err := rows.Scan(
		&cred.ID,
		&cred.TenantID,
		&cred.ConnectionID,
		&cred.Role,
		&cred.SchemaID,
		&cred.CredDefID,
		&cred.InitiatedByUs,
		&cred.Archived,
		&cred.Created,
		&cred.Approved,
		&cred.Issued,
		&cred.Failed,
		&cred.Cursor,
		&a.ID,
		&a.Name,
		&a.Value,
	)

	cred.Attributes = make([]*graph.CredentialValue, 0)
	if previous.ID == cred.ID {
		cred.Attributes = append(cred.Attributes, previous.Attributes...)
		cred.Attributes = append(cred.Attributes, a)
	} else {
		cred.Attributes = append(cred.Attributes, a)
	}

	return cred, err
}

func (pg *Database) GetCredential(id, tenantID string) (cred *model.Credential, err error) {
	defer err2.Handle(&err, "GetCredential")

	sqlCredentialSelectByID := sqlCredentialSelect + " credential" + sqlCredentialJoin +
		" WHERE credential.id=$1 AND tenant_id=$2" +
		" ORDER BY credential_attribute.index"

	cred = &model.Credential{}
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		cred = try.To1(readRowToCredential(rows, cred))
		return
	}, sqlCredentialSelectByID, id, tenantID))

	return
}

func (pg *Database) getCredentialsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,

	initialArgs []interface{},
) (c *model.Credentials, err error) {
	defer err2.Handle(&err, "GetCredentials")

	query, args := getBatchQuery(queries, batch, tenantID, initialArgs)

	c = &model.Credentials{
		Credentials:     make([]*model.Credential, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	prevCredential := &model.Credential{}
	var credential *model.Credential
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		credential = try.To1(readRowToCredential(rows, prevCredential))
		if prevCredential.ID != "" && prevCredential.ID != credential.ID {
			c.Credentials = append(c.Credentials, prevCredential)
		}
		prevCredential = credential
		return
	}, query, args...))

	// ensure also last credential is added
	lastCredentialID := ""
	if len(c.Credentials) > 0 {
		lastCredentialID = c.Credentials[len(c.Credentials)-1].ID
	}

	if prevCredential.ID != lastCredentialID {
		c.Credentials = append(c.Credentials, prevCredential)
	}

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

func sqlCredentialBatchWhere(cursorParam, connectionParam, limitParam string, desc, before bool) string {
	const issuedNotNull = " AND issued > timestamp '0001-01-01' "
	const whereTenantID = " WHERE tenant_id=$1 "
	order := sqlAsc
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
		order = sqlDesc
		cursorOrder = sqlOrderByCursorDesc
	}
	where := whereTenantID + cursor + connection + issuedNotNull
	return sqlCredentialSelect + " (SELECT * FROM credential " + where + cursorOrder + " " + limitParam + ") AS credential " +
		sqlCredentialJoin + " ORDER BY cursor " + order + ", credential_attribute.index"
}

func (pg *Database) GetCredentials(info *paginator.BatchInfo, tenantID string, connectionID *string) (c *model.Credentials, err error) {
	if connectionID == nil {
		return pg.getCredentialsForQuery(&queryInfo{
			Asc:        sqlCredentialBatchWhere("", "", "$2", false, false),
			Desc:       sqlCredentialBatchWhere("", "", "$2", true, false),
			AfterAsc:   sqlCredentialBatchWhere("$2", "", "$3", false, false),
			AfterDesc:  sqlCredentialBatchWhere("$2", "", "$3", true, false),
			BeforeAsc:  sqlCredentialBatchWhere("$2", "", "$3", false, true),
			BeforeDesc: sqlCredentialBatchWhere("$2", "", "$3", true, true),
		},
			info,
			tenantID,
			[]interface{}{},
		)
	}
	return pg.getCredentialsForQuery(&queryInfo{
		Asc:        sqlCredentialBatchWhere("", "$2", "$3", false, false),
		Desc:       sqlCredentialBatchWhere("", "$2", "$3", true, false),
		AfterAsc:   sqlCredentialBatchWhere("$2", "$3", "$4", false, false),
		AfterDesc:  sqlCredentialBatchWhere("$2", "$3", "$4", true, false),
		BeforeAsc:  sqlCredentialBatchWhere("$2", "$3", "$4", false, true),
		BeforeDesc: sqlCredentialBatchWhere("$2", "$3", "$4", true, true),
	},
		info,
		tenantID,
		[]interface{}{*connectionID},
	)
}

func (pg *Database) GetCredentialCount(tenantID string, connectionID *string) (count int, err error) {
	defer err2.Handle(&err, "GetCredentialCount")
	const (
		sqlCredentialBatchWhere           = " WHERE tenant_id=$1 AND issued > timestamp '0001-01-01' "
		sqlCredentialBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2 AND issued > timestamp '0001-01-01' "
	)
	count = try.To1(pg.getCount(
		"credential",
		sqlCredentialBatchWhere,
		sqlCredentialBatchWhereConnection,
		tenantID,
		connectionID,
	))
	return
}

func (pg *Database) GetConnectionForCredential(id, tenantID string) (*model.Connection, error) {
	return pg.getConnectionForObject("credential", "connection_id", id, tenantID)
}

func (pg *Database) ArchiveCredential(id, tenantID string) (err error) {
	defer err2.Handle(&err, "ArchiveCredential")

	var (
		//#nosec
		sqlCredentialArchive = "UPDATE credential SET archived=$1 WHERE id = $2 and tenant_id = $3 RETURNING id"
	)

	now := utils.CurrentTime()
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&id)
		},
		sqlCredentialArchive,
		now,
		id,
		tenantID,
	))
	return
}

func (pg *Database) SearchCredentials(
	tenantID string,
	proofAttributes []*graph.ProofAttribute,
) (res []*graph.ProvableAttribute, err error) {
	defer err2.Handle(&err, "SearchCredentials")

	assert.That(len(proofAttributes) > 0, "cannot search credentials for empty proof")

	credDefIDs := make([]string, 0)
	names := make([]string, 0)
	for _, attr := range proofAttributes {
		if attr.CredDefID != "" {
			credDefIDs = append(credDefIDs, attr.CredDefID)
		}
		names = append(names, attr.Name)
	}

	utils.LogMed().Infof("Searching for credentials with cred def id %v, name %v", credDefIDs, names)

	var (
		sqlCredentialSearch = "SELECT credential.id, name, cred_def_id, value FROM credential " + sqlCredentialJoin +
			" WHERE tenant_id=$1 AND issued > timestamp '0001-01-01' AND (cred_def_id=ANY($2::varchar[]) OR name=ANY($3::varchar[]))" +
			" ORDER BY credential.created"
	)

	type searchResult struct {
		credID    string
		attrName  string
		credDefID string
		credValue string
	}
	searchResults := make([]*searchResult, 0)
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		s := &searchResult{}
		try.To(rows.Scan(&s.credID, &s.attrName, &s.credDefID, &s.credValue))
		searchResults = append(searchResults, s)
		return
	}, sqlCredentialSearch, tenantID, pq.Array(credDefIDs), pq.Array(names)))

	res = make([]*graph.ProvableAttribute, 0)
	for _, attr := range proofAttributes {
		provableAttr := &graph.ProvableAttribute{}
		provableAttr.ID = attr.ID
		provableAttr.Attribute = attr
		provableAttr.Credentials = make([]*graph.CredentialMatch, 0)
		for _, value := range searchResults {
			if value.attrName == attr.Name && (attr.CredDefID == "" || attr.CredDefID == value.credDefID) {
				provableAttr.Credentials = append(provableAttr.Credentials, &graph.CredentialMatch{
					ID:           attr.ID + "/" + value.credID,
					CredentialID: value.credID,
					Value:        value.credValue,
				})
			}
		}
		res = append(res, provableAttr)
	}

	return res, err
}
