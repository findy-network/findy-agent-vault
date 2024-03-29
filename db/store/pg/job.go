package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/lib/pq"
)

var (
	jobFields = []string{"id", "tenant_id", "protocol_type", "protocol_connection_id", "protocol_credential_id", "protocol_proof_id",
		"protocol_message_id", "connection_id", "status", "result", "initiated_by_us", "updated"}
	sqlJobBaseFields = sqlFields("", jobFields)
	sqlJobInsert     = "INSERT INTO job " + "(" + sqlJobBaseFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, (now() at time zone 'UTC')) RETURNING " + sqlInsertFields
	sqlJobSelect = "SELECT " + sqlJobBaseFields + ", created, cursor FROM"
)

func (pg *Database) getJobForObject(objectName, objectID, tenantID string) (j *model.Job, err error) {
	defer err2.Handle(&err, "getJobForObject")

	sqlJobSelectJoin := "SELECT " + sqlFields("job", jobFields) + ", job.created, job.cursor FROM"
	sqlJobSelectByObjectID := sqlJobSelectJoin +
		" job INNER JOIN " + objectName + " ON " + objectName +
		".job_id=job.id WHERE " + objectName + ".id = $1 AND job.tenant_id = $2"

	j = &model.Job{}
	try.To(pg.doRowQuery(
		readRowToJob(j),
		sqlJobSelectByObjectID,
		objectID,
		tenantID,
	))

	return
}

func (pg *Database) AddJob(j *model.Job) (job *model.Job, err error) {
	defer err2.Handle(&err, "AddJob")

	job = &model.Job{}
	*job = *j
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&job.ID, &job.Created, &job.Cursor)
		},
		sqlJobInsert,
		j.ID,
		j.TenantID,
		j.ProtocolType,
		j.ProtocolConnectionID,
		j.ProtocolCredentialID,
		j.ProtocolProofID,
		j.ProtocolMessageID,
		j.ConnectionID,
		j.Status,
		j.Result,
		j.InitiatedByUs,
	))

	return job, err
}

func (pg *Database) UpdateJob(arg *model.Job) (j *model.Job, err error) {
	defer err2.Handle(&err, "UpdateJob")

	sqlJobUpdate := "UPDATE job " +
		"SET protocol_connection_id=$1, protocol_credential_id=$2, protocol_proof_id=$3, protocol_message_id=$4," +
		" connection_id=$5, status=$6, result=$7, updated=(now() at time zone 'UTC')" +
		" WHERE id = $8 AND tenant_id = $9" +
		" RETURNING " + sqlJobBaseFields + ", created, cursor"

	j = &model.Job{}
	try.To(pg.doRowQuery(
		readRowToJob(j),
		sqlJobUpdate,
		arg.ProtocolConnectionID,
		arg.ProtocolCredentialID,
		arg.ProtocolProofID,
		arg.ProtocolMessageID,
		arg.ConnectionID,
		arg.Status,
		arg.Result,
		arg.ID,
		arg.TenantID,
	))
	return j, err
}

func rowToJob(rows *sql.Rows) (n *model.Job, err error) {
	n = &model.Job{}
	return n, readRowToJob(n)(rows)
}

func readRowToJob(n *model.Job) func(*sql.Rows) error {
	return func(rows *sql.Rows) error {
		return rows.Scan(
			&n.ID,
			&n.TenantID,
			&n.ProtocolType,
			&n.ProtocolConnectionID,
			&n.ProtocolCredentialID,
			&n.ProtocolProofID,
			&n.ProtocolMessageID,
			&n.ConnectionID,
			&n.Status,
			&n.Result,
			&n.InitiatedByUs,
			&n.Updated,
			&n.Created,
			&n.Cursor,
		)
	}
}

func (pg *Database) GetJob(id, tenantID string) (job *model.Job, err error) {
	defer err2.Handle(&err, "GetJob")

	sqlJobSelectByID := sqlJobSelect + " job WHERE id=$1 AND tenant_id=$2"

	job = &model.Job{}
	try.To(pg.doRowQuery(
		readRowToJob(job),
		sqlJobSelectByID,
		id,
		tenantID,
	))

	return
}

func (pg *Database) getJobsForQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,
	initialArgs []interface{},
) (j *model.Jobs, err error) {
	defer err2.Handle(&err, "GetJobs")

	query, args := getBatchQuery(queries, batch, tenantID, initialArgs)

	j = &model.Jobs{
		Jobs:            make([]*model.Job, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var job *model.Job
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		job = try.To1(rowToJob(rows))
		j.Jobs = append(j.Jobs, job)
		return
	}, query, args...))

	if batch.Count < len(j.Jobs) {
		j.Jobs = j.Jobs[:batch.Count]
		if batch.Tail {
			j.HasPreviousPage = true
		} else {
			j.HasNextPage = true
		}
	}

	if batch.After > 0 {
		j.HasPreviousPage = true
	}
	if batch.Before > 0 {
		j.HasNextPage = true
	}

	// Reverse order for tail first
	if batch.Tail {
		sort.Slice(j.Jobs, func(i, k int) bool {
			return j.Jobs[i].Created.Sub(j.Jobs[k].Created) < 0
		})
	}

	return j, err
}

func sqlJobBatchWhere(fetchAll bool, cursorParam, connectionParam, limitParam string, desc, before bool) string {
	const whereTenantID = " WHERE tenant_id=$1 "
	whereStatus := " AND status != 'COMPLETE' "
	cursorOrder := sqlOrderByCursorAsc
	cursor := ""
	connection := ""
	compareChar := sqlGreaterThan
	if fetchAll {
		whereStatus = ""
	}
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
	where := whereTenantID + cursor + connection + whereStatus
	return sqlJobSelect + " job " + where + cursorOrder + " " + limitParam
}

func (pg *Database) GetJobs(info *paginator.BatchInfo, tenantID string, connectionID *string, completed *bool) (c *model.Jobs, err error) {
	fetchAll := completed != nil && *completed

	if connectionID == nil {
		return pg.getJobsForQuery(&queryInfo{
			Asc:        sqlJobBatchWhere(fetchAll, "", "", "$2", false, false),
			Desc:       sqlJobBatchWhere(fetchAll, "", "", "$2", true, false),
			AfterAsc:   sqlJobBatchWhere(fetchAll, "$2", "", "$3", false, false),
			AfterDesc:  sqlJobBatchWhere(fetchAll, "$2", "", "$3", true, false),
			BeforeAsc:  sqlJobBatchWhere(fetchAll, "$2", "", "$3", false, true),
			BeforeDesc: sqlJobBatchWhere(fetchAll, "$2", "", "$3", true, true),
		},
			info,
			tenantID,
			[]interface{}{},
		)
	}
	return pg.getJobsForQuery(&queryInfo{
		Asc:        sqlJobBatchWhere(fetchAll, "", "$2", "$3", false, false),
		Desc:       sqlJobBatchWhere(fetchAll, "", "$2", "$3", true, false),
		AfterAsc:   sqlJobBatchWhere(fetchAll, "$2", "$3", "$4", false, false),
		AfterDesc:  sqlJobBatchWhere(fetchAll, "$2", "$3", "$4", true, false),
		BeforeAsc:  sqlJobBatchWhere(fetchAll, "$2", "$3", "$4", false, true),
		BeforeDesc: sqlJobBatchWhere(fetchAll, "$2", "$3", "$4", true, true),
	},
		info,
		tenantID,
		[]interface{}{*connectionID},
	)
}

func (pg *Database) GetJobCount(tenantID string, connectionID *string, completed *bool) (count int, err error) {
	defer err2.Handle(&err, "GetJobCount")
	const (
		sqlJobBatchWhere              = " WHERE tenant_id=$1 AND status != 'COMPLETE'"
		sqlJobBatchWhereConnection    = " WHERE tenant_id=$1 AND connection_id=$2 AND status != 'COMPLETE'"
		sqlJobBatchWhereAll           = " WHERE tenant_id=$1"
		sqlJobBatchWhereConnectionAll = " WHERE tenant_id=$1 AND connection_id=$2"
	)

	fetchAll := completed != nil && *completed

	qWhere := sqlJobBatchWhere
	qWhereConnection := sqlJobBatchWhereConnection
	if fetchAll {
		qWhere = sqlJobBatchWhereAll
		qWhereConnection = sqlJobBatchWhereConnectionAll
	}

	count = try.To1(pg.getCount(
		"job",
		qWhere,
		qWhereConnection,
		tenantID,
		connectionID,
	))
	return
}

func (pg *Database) GetConnectionForJob(id, tenantID string) (*model.Connection, error) {
	return pg.getConnectionForObject("job", "connection_id", id, tenantID)
}

func (pg *Database) GetJobOutput(id, tenantID string, protocolType graph.ProtocolType) (output *model.JobOutput, err error) {
	defer err2.Handle(&err)
	switch protocolType {
	case graph.ProtocolTypeConnection:
		connection := try.To1(pg.getConnectionForObject("job", "protocol_connection_id", id, tenantID))
		return &model.JobOutput{Connection: connection}, nil
	case graph.ProtocolTypeCredential:
		credential := try.To1(pg.getCredentialForObject("job", "protocol_credential_id", id, tenantID))
		return &model.JobOutput{Credential: credential}, nil
	case graph.ProtocolTypeProof:
		proof := try.To1(pg.getProofForObject("job", "protocol_proof_id", id, tenantID))
		return &model.JobOutput{Proof: proof}, nil
	case graph.ProtocolTypeBasicMessage:
		message := try.To1(pg.getMessageForObject("job", "protocol_message_id", id, tenantID))
		return &model.JobOutput{Message: message}, nil
	case graph.ProtocolTypeNone:
		break
	}
	return &model.JobOutput{}, nil
}

func (pg *Database) GetOpenProofJobs(tenantID string, proofAttributes []*graph.ProofAttribute) (jobs []*model.Job, err error) {
	defer err2.Handle(&err)

	credDefIDs := make([]string, 0)
	names := make([]string, 0)
	for _, attr := range proofAttributes {
		if attr.CredDefID != "" {
			credDefIDs = append(credDefIDs, attr.CredDefID)
		}
		names = append(names, attr.Name)
	}

	query := "SELECT DISTINCT " + sqlFields("job", jobFields) + ", created, cursor FROM job " +
		"INNER JOIN proof_attribute ON proof_attribute.proof_id = job.protocol_proof_id " +
		"WHERE tenant_id=$1 AND status = 'BLOCKED' AND protocol_type = 'PROOF' AND " +
		"(cred_def_id=ANY($2::varchar[]) OR name=ANY($3::varchar[]))"

	jobs = make([]*model.Job, 0)
	var job *model.Job
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		job = try.To1(rowToJob(rows))
		jobs = append(jobs, job)
		return
	}, query, tenantID, pq.Array(credDefIDs), pq.Array(names)))

	return jobs, nil
}
