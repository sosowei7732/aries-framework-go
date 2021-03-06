/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package introduce

import (
	"github.com/google/uuid"

	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
)

const (
	// protocol instance ID
	metaPIID         = Introduce + "_pi_id"
	metaSkipProposal = Introduce + "_skip_proposal"
	metaInvitation   = Introduce + "_invitation"
	metaRecipients   = Introduce + "_recipients"
)

// Opt describes option signature for the Continue function
type Opt func(m map[string]interface{})

// WithInvitation is used when introducee wants to provide invitation.
// NOTE: Introducee can provide invitation only after receiving ProposalMsgType
// USAGE: event.Continue(WithInvitation(inv))
func WithInvitation(inv *didexchange.Invitation) Opt {
	return func(m map[string]interface{}) {
		m[metaInvitation] = service.NewDIDCommMsgMap(inv)
	}
}

// WithPublicInvitation is used when introducer wants to provide public invitation.
// NOTE: Introducer can provide invitation only after receiving RequestMsgType
// USAGE: event.Continue(WithPublicInvitation(inv, to))
func WithPublicInvitation(inv *didexchange.Invitation, to *To) Opt {
	return func(m map[string]interface{}) {
		m[metaInvitation] = service.NewDIDCommMsgMap(inv)
		m[metaSkipProposal] = true
		m[metaRecipients] = []interface{}{&Recipient{
			To: to,
		}}
	}
}

// WithRecipients is used when the introducer does not have a public invitation
// but he is willing to introduce agents to each other.
// NOTE: Introducer can provide recipients only after receiving RequestMsgType.
// USAGE: event.Continue(WithRecipients(to, recipient))
func WithRecipients(to *To, recipient *Recipient) Opt {
	return func(m map[string]interface{}) {
		m[metaRecipients] = []interface{}{
			&Recipient{To: to}, recipient,
		}
	}
}

// WrapWithMetadataPIID wraps message with metadata.
// The function is used by the introduce client to define that a few messages are related to each other.
// e.g When two proposals are sent simultaneously piID helps the protocol to determine that messages are related.
func WrapWithMetadataPIID(msgMap ...service.DIDCommMsg) {
	var piID = uuid.New().String()

	for _, msg := range msgMap {
		msg.Metadata()[metaPIID] = piID
	}
}

// WrapWithMetadataPublicInvitation wraps message with metadata.
// The function is used by the introduce client to define skip proposal.
// It also saves invitation and will provide it later to the introducee.
func WrapWithMetadataPublicInvitation(msg service.DIDCommMsgMap, inv *didexchange.Invitation) {
	msg.Metadata()[metaInvitation] = service.NewDIDCommMsgMap(inv)
	msg.Metadata()[metaSkipProposal] = true
}

func copyMetadata(from, to service.DIDCommMsg) {
	for k, v := range from.Metadata() {
		to.Metadata()[k] = v
	}
}
