package findy

import (
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
)

func statusToConnection(status *agency.ProtocolStatus) *model.Connection {
	connection := status.GetDIDExchange()
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
	credential := status.GetIssueCredential()
	if credential != nil {
		role := graph.CredentialRoleHolder
		if status.State.GetProtocolID().Role != agency.Protocol_ADDRESSEE {
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
			SchemaID:      credential.SchemaID,
			CredDefID:     credential.CredDefID,
			Attributes:    values,
			InitiatedByUs: false,
		}
	}
	return nil
}

func statusToProof(status *agency.ProtocolStatus) *model.Proof {
	proof := status.GetPresentProof()
	if proof != nil {
		role := graph.ProofRoleProver
		if status.State.GetProtocolID().Role != agency.Protocol_ADDRESSEE {
			role = graph.ProofRoleVerifier
		}
		attributes := make([]*graph.ProofAttribute, 0)
		for _, v := range proof.Attributes {
			attributes = append(attributes, &graph.ProofAttribute{
				Name:      v.Name,
				CredDefID: v.CredDefID,
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
			SentByMe: status.State.GetProtocolID().Role != agency.Protocol_ADDRESSEE,
		}
	}
	return nil
}
