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

func sqlProofSelectBatchFor(tenantOrder, limit, cursorOrder string) string {
	return sqlProofSelect +
		" (SELECT * FROM proof " + tenantOrder + " " + limit + ") AS proof " +
		sqlProofJoin + " ORDER BY cursor " + cursorOrder + ", proof_attribute.index"
}

const (
	sqlProofBatchWhere           = " WHERE tenant_id=$1 AND verified IS NOT NULL "
	sqlProofBatchWhereConnection = " WHERE tenant_id=$1 AND connection_id=$2 AND verified IS NOT NULL "

	sqlProofFields = "tenant_id, connection_id, role, initiated_by_us, result"
	sqlProofInsert = "INSERT INTO proof " + "(" + sqlProofFields + ") " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id, created, cursor"
	sqlProofSelect = "SELECT proof.id, " + sqlProofFields +
		", created, approved, verified, failed, cursor, proof_attribute.id, name, value, cred_def_id FROM"
	sqlProofJoin = " INNER JOIN proof_attribute on proof_id = proof.id"
)

func (pg *Database) addProofAttributes(id string, attributes []*graph.ProofAttribute) (a []*graph.ProofAttribute, err error) {
	defer returnErr("addProofAttributes", &err)

	query := constructProofAttributeInsert(len(attributes))
	args := make([]interface{}, 0)
	for index, a := range attributes {
		args = append(args, []interface{}{id, a.Name, a.Value, a.CredDefID, index}...)
	}

	rows, err := pg.db.Query(query, args...)
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

func (pg *Database) AddProof(p *model.Proof) (n *model.Proof, err error) {
	defer returnErr("AddProof", &err)

	if len(p.Attributes) == 0 {
		panic("Attributes are always required for proof.")
	}
	rows, err := pg.db.Query(
		sqlProofInsert,
		p.TenantID,
		p.ConnectionID,
		p.Role,
		p.InitiatedByUs,
		p.Result,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewProof(p)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	attributes, err := pg.addProofAttributes(n.ID, n.Attributes)
	err2.Check(err)

	n.Attributes = attributes
	return n, err
}

func (pg *Database) UpdateProof(p *model.Proof) (n *model.Proof, err error) {
	defer returnErr("UpdateProof", &err)

	const sqlProofUpdate = "UPDATE proof SET approved=$1, verified=$2, failed=$3 WHERE id = $4" // TODO: tenant_id, connection_id?

	_, err = pg.db.Exec(
		sqlProofUpdate,
		p.Approved,
		p.Verified,
		p.Failed,
		p.ID,
	)
	err2.Check(err)
	return p, err
}

func readRowToProof(rows *sql.Rows, previous *model.Proof) (*model.Proof, error) {
	a := &graph.ProofAttribute{}
	var approved sql.NullTime
	var verified sql.NullTime
	var failed sql.NullTime

	n := model.NewProof(nil)

	err := rows.Scan(
		&n.ID,
		&n.TenantID,
		&n.ConnectionID,
		&n.Role,
		&n.InitiatedByUs,
		&n.Result,
		&n.Created,
		&approved,
		&verified,
		&failed,
		&n.Cursor,
		&a.ID,
		&a.Name,
		&a.Value,
		&a.CredDefID,
	)

	if approved.Valid {
		n.Approved = &approved.Time
	}
	if verified.Valid {
		n.Verified = &verified.Time
	}
	if failed.Valid {
		n.Failed = &failed.Time
	}

	n.Attributes = make([]*graph.ProofAttribute, 0)
	if previous.ID == n.ID {
		n.Attributes = append(n.Attributes, previous.Attributes...)
		n.Attributes = append(n.Attributes, a)
	} else {
		n.Attributes = append(n.Attributes, a)
	}

	return n, err
}

func (pg *Database) GetProof(id, tenantID string) (p *model.Proof, err error) {
	defer returnErr("GetProof", &err)

	const sqlProofSelectByID = sqlProofSelect + " proof" + sqlProofJoin +
		" WHERE proof.id=$1 AND tenant_id=$2" +
		" ORDER BY proof_attribute.index"

	rows, err := pg.db.Query(sqlProofSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	p = model.NewProof(nil)
	for rows.Next() {
		p, err = readRowToProof(rows, p)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (pg *Database) getProofsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	initialArgs []interface{},
) (p *model.Proofs, err error) {
	defer returnErr("GetProofs", &err)

	query, args := getBatchQuery(queries, batch, initialArgs)
	rows, err := pg.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	p = &model.Proofs{
		Proofs:          make([]*model.Proof, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	prevProof := model.NewProof(nil)
	var proof *model.Proof
	for rows.Next() {
		proof, err = readRowToProof(rows, prevProof)
		err2.Check(err)
		if prevProof.ID != "" && prevProof.ID != proof.ID {
			p.Proofs = append(p.Proofs, prevProof)
		}
		prevProof = proof
	}

	// ensure also last proof is added
	if prevProof.ID != p.Proofs[len(p.Proofs)-1].ID {
		p.Proofs = append(p.Proofs, prevProof)
	}

	err = rows.Err()
	err2.Check(err)

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

func (pg *Database) GetProofs(info *paginator.BatchInfo, tenantID string, connectionID *string) (c *model.Proofs, err error) {
	if connectionID == nil {
		return pg.getProofsForQuery(&queryInfo{
			Asc:        sqlProofSelectBatchFor(sqlProofBatchWhere+sqlOrderByAsc(""), "$2", "ASC"),
			Desc:       sqlProofSelectBatchFor(sqlProofBatchWhere+sqlOrderByDesc(""), "$2", "DESC"),
			AfterAsc:   sqlProofSelectBatchFor(sqlProofBatchWhere+" AND cursor > $2"+sqlOrderByAsc(""), "$3", "ASC"),
			AfterDesc:  sqlProofSelectBatchFor(sqlProofBatchWhere+" AND cursor > $2"+sqlOrderByDesc(""), "$3", "DESC"),
			BeforeAsc:  sqlProofSelectBatchFor(sqlProofBatchWhere+" AND cursor < $2"+sqlOrderByAsc(""), "$3", "ASC"),
			BeforeDesc: sqlProofSelectBatchFor(sqlProofBatchWhere+" AND cursor < $2"+sqlOrderByDesc(""), "$3", "DESC"),
		},
			info,
			[]interface{}{tenantID},
		)
	}
	return pg.getProofsForQuery(&queryInfo{
		Asc:        sqlProofSelectBatchFor(sqlProofBatchWhereConnection+sqlOrderByAsc(""), "$3", "ASC"),
		Desc:       sqlProofSelectBatchFor(sqlProofBatchWhereConnection+sqlOrderByDesc(""), "$3", "DESC"),
		AfterAsc:   sqlProofSelectBatchFor(sqlProofBatchWhereConnection+" AND cursor > $3"+sqlOrderByAsc(""), "$4", "ASC"),
		AfterDesc:  sqlProofSelectBatchFor(sqlProofBatchWhereConnection+" AND cursor > $3"+sqlOrderByDesc(""), "$4", "DESC"),
		BeforeAsc:  sqlProofSelectBatchFor(sqlProofBatchWhereConnection+" AND cursor < $3"+sqlOrderByAsc(""), "$4", "ASC"),
		BeforeDesc: sqlProofSelectBatchFor(sqlProofBatchWhereConnection+" AND cursor < $3"+sqlOrderByDesc(""), "$4", "DESC"),
	},
		info,
		[]interface{}{tenantID, *connectionID},
	)
}

func (pg *Database) GetProofCount(tenantID string, connectionID *string) (count int, err error) {
	defer returnErr("GetProofCount", &err)
	count, err = pg.getCount(
		"proof",
		sqlProofBatchWhere,
		sqlProofBatchWhereConnection,
		tenantID,
		connectionID,
	)
	err2.Check(err)
	return
}
