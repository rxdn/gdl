package events

import "github.com/rxdn/gdl/objects/entitlement"

type EntitlementCreate struct {
	entitlement.Entitlement
}

type EntitlementUpdate struct {
	entitlement.Entitlement
}

type EntitlementDelete struct {
	entitlement.Entitlement
}
