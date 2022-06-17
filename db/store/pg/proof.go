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
	"github.com/lainio/err2/try"
)

func constructProofAttributeInsert(count int) string {
	const sqlProofAttributeInsert = "INSERT INTO proof_attribute (proof_id, name, value, cred_def_id, index) VALUES "

	result := sqlProofAttributeInsert
	paramCount := 5
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
	proofFields      = []string{"tenant_id", "connection_id", "role", "initiated_by_us", "result", "archived", "provable"}
	proofExtraFields = []string{"created", "approved", "verified", "failed", "cursor"}

	sqlProofBaseFields = sqlFields("", proofFields)
	sqlProofSelect     = "SELECT proof.id, " + sqlProofBaseFields +
		", " + sqlFields("", proofExtraFields) + ", proof_attribute.id, name, value, cred_def_id FROM"
)

const (
	sqlProofJoin = " INNER JOIN proof_attribute on proof_attribute.proof_id = proof.id"
)

func (pg *Database) getProofForObject(objectName, columnName, objectID, tenantID string) (c *model.Proof, err error) {
	defer err2.Annotate("getProofForObject", &err)

	sqlProofJoinSelect := "SELECT proof.id, " + sqlFields("proof", proofFields) +
		", " + sqlFields("proof", proofExtraFields) +
		", proof_attribute.id, proof_attribute.name, proof_attribute.value, proof_attribute.cred_def_id FROM"
	sqlProofSelectByObjectID := sqlProofJoinSelect + " proof " + sqlProofJoin +
		" INNER JOIN " + objectName + " ON " + objectName +
		"." + columnName + "=proof.id WHERE " + objectName + ".id = $1 AND proof.tenant_id = $2"

	c = &model.Proof{}
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Return(&err)
		c = try.To1(readRowToProof(rows, c))
		return
	}, sqlProofSelectByObjectID, objectID, tenantID))

	return
}

func (pg *Database) addProofAttributes(id string, attributes []*graph.ProofAttribute) (a []*graph.ProofAttribute, err error) {
	defer err2.Annotate("addProofAttributes", &err)

	query := constructProofAttributeInsert(len(attributes))
	args := make([]interface{}, 0)
	for index, a := range attributes {
		// TODO: save values when received
		args = append(args, []interface{}{id, a.Name, "", a.CredDefID, index}...)
	}

	index := 0
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Return(&err)
		try.To(rows.Scan(&attributes[index].ID))
		index++
		return
	}, query, args...))

	return attributes, nil
}

func (pg *Database) AddProof(p *model.Proof) (proof *model.Proof, err error) {
	defer err2.Annotate("AddProof", &err)

	if len(p.Attributes) == 0 {
		panic("Attributes are always required for proof.")
	}

	var (
		sqlProofInsert = "INSERT INTO proof " + "(" + sqlProofBaseFields + ") " +
			"VALUES (" + sqlArguments(proofFields) + ") RETURNING " + sqlInsertFields
	)

	proof = &model.Proof{}
	*proof = *p

	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&proof.ID, &proof.Created, &proof.Cursor)
		},
		sqlProofInsert,
		p.TenantID,
		p.ConnectionID,
		p.Role,
		p.InitiatedByUs,
		p.Result,
		p.Archived,
		p.Provable,
	))

	attributes := try.To1(pg.addProofAttributes(proof.ID, proof.Attributes))

	proof.Attributes = attributes
	return proof, err
}

func (pg *Database) UpdateProof(p *model.Proof) (n *model.Proof, err error) {
	defer err2.Annotate("UpdateProof", &err)

	const (
		// TODO: tenant id + connection id
		sqlProofUpdate          = "UPDATE proof SET provable=$1, approved=$2, verified=$3, failed=$4 WHERE id = $5"
		sqlProofAttributeUpdate = "UPDATE proof_attribute SET value = (CASE %s END) WHERE id IN (%s)"
	)

	_, err = pg.db.Exec(
		sqlProofUpdate,
		p.Provable,
		p.Approved,
		p.Verified,
		p.Failed,
		p.ID,
	)
	try.To(err)

	valueUpdate := ""
	ids := ""
	args := make([]interface{}, 0)
	for i, value := range p.Values {
		round := i*2 + 1
		valueUpdate += fmt.Sprintf("WHEN id = $%d THEN $%d ", round, round+1)
		args = append(args, value.AttributeID, value.Value)
		if ids != "" {
			ids += ","
		}
		ids += fmt.Sprintf("'%s'", value.AttributeID)
	}

	if valueUpdate != "" {
		_, err = pg.db.Exec(
			fmt.Sprintf(sqlProofAttributeUpdate, valueUpdate, ids),
			args...,
		)
		try.To(err)
	}

	for i, value := range p.Values {
		p.Values[i].ID = value.AttributeID
	}

	return p, err
}

func readRowToProof(rows *sql.Rows, previous *model.Proof) (*model.Proof, error) {
	a := &graph.ProofAttribute{}

	n := &model.Proof{}

	value := &graph.ProofValue{}

	err := rows.Scan(
		&n.ID,
		&n.TenantID,
		&n.ConnectionID,
		&n.Role,
		&n.InitiatedByUs,
		&n.Result,
		&n.Archived,
		&n.Provable,
		&n.Created,
		&n.Approved,
		&n.Verified,
		&n.Failed,
		&n.Cursor,
		&a.ID,
		&a.Name,
		&value.Value,
		&a.CredDefID,
	)

	n.Attributes = make([]*graph.ProofAttribute, 0)
	if previous.ID == n.ID {
		n.Attributes = append(n.Attributes, previous.Attributes...)
	}
	n.Attributes = append(n.Attributes, a)

	n.Values = make([]*graph.ProofValue, 0)
	if value.Value != "" {
		value.ID = a.ID
		value.AttributeID = a.ID
		if previous.ID == n.ID {
			n.Values = append(n.Values, previous.Values...)
		}
		n.Values = append(n.Values, value)
	}

	return n, err
}

func (pg *Database) GetProof(id, tenantID string) (p *model.Proof, err error) {
	defer err2.Annotate("GetProof", &err)

	sqlProofSelectByID := sqlProofSelect + " proof" + sqlProofJoin +
		" WHERE proof.id=$1 AND tenant_id=$2" +
		" ORDER BY proof_attribute.index"

	p = &model.Proof{}
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Return(&err)
		p = try.To1(readRowToProof(rows, p))
		return
	}, sqlProofSelectByID, id, tenantID))

	return
}

func (pg *Database) getProofsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,
	initialArgs []interface{},
) (p *model.Proofs, err error) {
	defer err2.Annotate("GetProofs", &err)

	query, args := getBatchQuery(queries, batch, tenantID, initialArgs)

	p = &model.Proofs{
		Proofs:          make([]*model.Proof, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	prevProof := &model.Proof{}
	var proof *model.Proof
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Return(&err)
		proof = try.To1(readRowToProof(rows, prevProof))
		if prevProof.ID != "" && prevProof.ID != proof.ID {
			p.Proofs = append(p.Proofs, prevProof)
		}
		prevProof = proof
		return
	}, query, args...))

	// ensure also last proof is added
	lastProofID := ""
	if len(p.Proofs) > 0 {
		lastProofID = p.Proofs[len(p.Proofs)-1].ID
	}
	if prevProof.ID != lastProofID {
		p.Proofs = append(p.Proofs, prevProof)
	}

	if batch.Count < len(p.Proofs) {
		p.Proofs = p.Proofs[:batch.Count]
		if batch.Tail {
			p.HasPreviousPage = true
		} else {
			p.HasNextPage = true
		}
	}

	if batch.After > 0 {
		p.HasPreviousPage = true
	}
	if batch.Before > 0 {
		p.HasNextPage = true
	}

	// Reverse order for tail first
	if batch.Tail {
		sort.Slice(p.Proofs, func(i, j int) bool {
			return p.Proofs[i].Created.Sub(p.Proofs[j].Created) < 0
		})
	}

	return p, err
}

func sqlProofBatchWhere(cursorParam, connectionParam, limitParam string, desc, before bool) string {
	const verifiedNotNull = " AND verified IS NOT NULL "
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
	where := whereTenantID + cursor + connection + verifiedNotNull
	return sqlProofSelect + " (SELECT * FROM proof " + where + cursorOrder + " " + limitParam + ") AS proof " +
		sqlProofJoin + " ORDER BY cursor " + order + ", proof_attribute.index"
}

func (pg *Database) GetProofs(info *paginator.BatchInfo, tenantID string, connectionID *string) (c *model.Proofs, err error) {
	if connectionID == nil {
		return pg.getProofsForQuery(&queryInfo{
			Asc:        sqlProofBatchWhere("", "", "$2", false, false),
			Desc:       sqlProofBatchWhere("", "", "$2", true, false),
			AfterAsc:   sqlProofBatchWhere("$2", "", "$3", false, false),
			AfterDesc:  sqlProofBatchWhere("$2", "", "$3", true, false),
			BeforeAsc:  sqlProofBatchWhere("$2", "", "$3", false, true),
			BeforeDesc: sqlProofBatchWhere("$2", "", "$3", true, true),
		},
			info,
			tenantID,
			[]interface{}{},
		)
	}
	return pg.getProofsForQuery(&queryInfo{
		Asc:        sqlProofBatchWhere("", "$2", "$3", false, false),
		Desc:       sqlProofBatchWhere("", "$2", "$3", true, false),
		AfterAsc:   sqlProofBatchWhere("$2", "$3", "$4", false, false),
		AfterDesc:  sqlProofBatchWhere("$2", "$3", "$4", true, false),
		BeforeAsc:  sqlProofBatchWhere("$2", "$3", "$4", false, true),
		BeforeDesc: sqlProofBatchWhere("$2", "$3", "$4", true, true),
	},
		info,
		tenantID,
		[]interface{}{*connectionID},
	)
}

func (pg *Database) GetProofCount(tenantID string, connectionID *string) (count int, err error) {
	defer err2.Annotate("GetProofCount", &err)
	const (
		sqlProofBatchWhere           = " WHERE tenant_id=$1 AND verified IS NOT NULL "
		sqlProofBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2 AND verified IS NOT NULL "
	)
	count, err = pg.getCount(
		"proof",
		sqlProofBatchWhere,
		sqlProofBatchWhereConnection,
		tenantID,
		connectionID,
	)
	try.To(err)
	return
}

func (pg *Database) GetConnectionForProof(id, tenantID string) (*model.Connection, error) {
	return pg.getConnectionForObject("proof", "connection_id", id, tenantID)
}

func (pg *Database) ArchiveProof(id, tenantID string) (err error) {
	defer err2.Annotate("ArchiveProof", &err)

	var (
		sqlProofArchive = "UPDATE proof SET archived=$1 WHERE id = $2 and tenant_id = $3 RETURNING id"
	)

	now := utils.CurrentTime()
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&id)
		},
		sqlProofArchive,
		now,
		id,
		tenantID,
	))
	return
}
