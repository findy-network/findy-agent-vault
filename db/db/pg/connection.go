package pg

const (
	sqlConnectionFields = "tenant_id, our_id, their_did, their_endpoint, their_label, invited"
	sqlConnectionInsert = "INSERT INTO connection " +
		"(" + sqlConnectionFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6)"
	sqlConnectionSelect     = "SELECT id, " + sqlConnectionFields + ", created, approved, cursor FROM connection"
	sqlConnectionSelectByID = sqlConnectionSelect + " WHERE id=$1"
)

/*func (p *Database) AddConnection(c *model.Connection) (err error) {
	defer returnErr("AddConnection", &err)

	_, err = p.db.Exec(sqlConnectionInsert, a.AgentID, a.Label)
	err2.Check(err)

	return
}

func (p *Database) GetAgent(id, agentID *string) (a *model.Agent, err error) {
	defer returnErr("GetAgent", &err)

	if id == nil && agentID == nil {
		panic(fmt.Errorf("either id or agent id is required"))
	}
	query := sqlAgentSelectByID
	queryID := id
	if id == nil {
		query = sqlAgentSelectByAgentID
		queryID = agentID
	}

	rows, err := p.db.Query(query, *queryID)
	err2.Check(err)
	defer rows.Close()

	a = &model.Agent{}
	if rows.Next() {
		err = rows.Scan(&a.ID, &a.AgentID, &a.Label, &a.Created)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}*/
