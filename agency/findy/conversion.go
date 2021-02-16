package findy

import (
	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
)

func statusToConnection(status *agency.ProtocolStatus) *model.Connection {
	connection := status.GetConnection()
	if connection != nil {
		return &model.Connection{
			OurDID:        connection.MyDid,
			TheirDID:      connection.TheirDid,
			TheirEndpoint: connection.TheirEndpoint,
			TheirLabel:    connection.TheirLabel,
		}
	}
	return nil
}

func statusToCredential(status *agency.ProtocolStatus) *model.Credential {
	credential := status.GetIssue()
	if credential != nil {
		role := graph.CredentialRoleHolder
		if status.State.GetProtocolId().Role != agency.Protocol_ADDRESSEE {
			role = graph.CredentialRoleIssuer
		}
		values := make([]*graph.CredentialValue, 0)
		for _, v := range credential.Attrs {
			values = append(values, &graph.CredentialValue{
				Name:  v.Name,
				Value: v.Value,
			})
		}
		return &model.Credential{
			Role:          role,
			SchemaID:      credential.SchemaId,
			CredDefID:     credential.CredDefId,
			Attributes:    values,
			InitiatedByUs: false,
		}
	}
	return nil
}

func statusToProof(status *agency.ProtocolStatus) *model.Proof {
	proof := status.GetProof()
	if proof != nil {
		role := graph.ProofRoleProver
		if status.State.GetProtocolId().Role != agency.Protocol_ADDRESSEE {
			role = graph.ProofRoleVerifier
		}
		attributes := make([]*graph.ProofAttribute, 0)
		for _, v := range proof.Attrs {
			attributes = append(attributes, &graph.ProofAttribute{
				Name:      v.Name,
				CredDefID: v.CredDefId,
			})
		}
		return &model.Proof{
			Role:          role,
			Attributes:    attributes,
			InitiatedByUs: false,
		}
	}
	return nil
}

func statusToMessage(status *agency.ProtocolStatus) *model.Message {
	message := status.GetBasicMessage()
	if message != nil {
		return &model.Message{
			Message: message.Content,
			// TODO: remove SentByMe from agency API
			SentByMe: status.State.GetProtocolId().Role != agency.Protocol_ADDRESSEE,
		}
	}
	return nil
}
