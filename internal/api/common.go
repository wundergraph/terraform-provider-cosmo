package api

import (
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
)

const (
	GraphQLWebsocketSubprotocolDefault            = "auto"
	GraphQLWebsocketSubprotocolGraphQLWS          = "graphql-ws"
	GraphQLWebsocketSubprotocolGraphQLTransportWS = "graphql-transport-ws"
)

func ResolveWebsocketSubprotocol(protocol *string) *common.GraphQLWebsocketSubprotocol {
	if protocol == nil {
		return nil
	}

	switch *protocol {
	case GraphQLWebsocketSubprotocolGraphQLWS:
		return common.GraphQLWebsocketSubprotocol_GRAPHQL_WEBSOCKET_SUBPROTOCOL_WS.Enum()
	case GraphQLWebsocketSubprotocolGraphQLTransportWS:
		return common.GraphQLWebsocketSubprotocol_GRAPHQL_WEBSOCKET_SUBPROTOCOL_TRANSPORT_WS.Enum()
	// GraphQLWebsocketSubprotocolDefault
	default:
		return common.GraphQLWebsocketSubprotocol_GRAPHQL_WEBSOCKET_SUBPROTOCOL_AUTO.Enum()
	}
}

const (
	GraphQLSubscriptionProtocolWS      = "ws"
	GraphQLSubscriptionProtocolSSE     = "sse"
	GraphQLSubscriptionProtocolSSEPost = "sse_post"
)

func ResolveSubscriptionProtocol(protocol *string) *common.GraphQLSubscriptionProtocol {
	if protocol == nil {
		return nil
	}

	switch *protocol {
	case GraphQLSubscriptionProtocolSSE:
		return common.GraphQLSubscriptionProtocol_GRAPHQL_SUBSCRIPTION_PROTOCOL_SSE.Enum()
	case GraphQLSubscriptionProtocolSSEPost:
		return common.GraphQLSubscriptionProtocol_GRAPHQL_SUBSCRIPTION_PROTOCOL_SSE_POST.Enum()
	// GraphQLSubscriptionProtocolWS
	default:
		return common.GraphQLSubscriptionProtocol_GRAPHQL_SUBSCRIPTION_PROTOCOL_WS.Enum()
	}
}
