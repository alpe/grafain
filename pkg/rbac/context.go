package rbac

import (
	"context"

	"github.com/iov-one/weave"
)

type contextKey int // local to the rbac module

const (
	contextRBACConditions  contextKey = iota
	contextRBACPermissions contextKey = iota
)

// withRBAC creates a new context from parent context with Roles conditions and permissions added.
func withRBAC(ctx weave.Context, roleIDsToRoles map[string]Role) weave.Context {
	conds, _ := ctx.Value(contextRBACConditions).([]weave.Condition)
	perms, _ := ctx.Value(contextRBACPermissions).(map[Permission]struct{})
	if perms == nil {
		perms = make(map[Permission]struct{})
	}

	for id, role := range roleIDsToRoles {
		conds = append(conds, RoleCondition([]byte(id)))
		for _, v := range role.Permissions {
			perms[v] = struct{}{}
		}
	}

	newCtx := context.WithValue(ctx, contextRBACConditions, conds)
	newCtx = context.WithValue(newCtx, contextRBACPermissions, perms)
	return newCtx
}

// Authenticate gets/sets permissions on the given context key
type Authenticate struct {
}

// GetConditions returns permissions previously set on this context
func (a Authenticate) GetConditions(ctx weave.Context) []weave.Condition {
	val, _ := ctx.Value(contextRBACConditions).([]weave.Condition)
	if val == nil {
		return nil
	}
	return val
}

// HasAddress returns true iff this address is in GetConditions
func (a Authenticate) HasAddress(ctx weave.Context, addr weave.Address) bool {
	for _, s := range a.GetConditions(ctx) {
		if addr.Equals(s.Address()) {
			return true
		}
	}
	return false
}

type Authorize struct{}

// HasPermission for authorization
func (Authorize) HasPermission(ctx weave.Context, p Permission) bool {
	val, _ := ctx.Value(contextRBACPermissions).(map[Permission]struct{})
	if val == nil {
		return false
	}
	_, ok := val[p]
	return ok
}
