package app

var (
	_ MessageClientInfoImpl = RestrictedMessageClientInfo(nil)
)

type RestrictedMessageClientInfo MessageClient
