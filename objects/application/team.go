package application

import "github.com/rxdn/gdl/objects/user"

type (
	Team struct {
		Icon        *string      `json:"icon"`
		Id          uint64       `json:"id,string"`
		Members     []TeamMember `json:"members"`
		Name        string       `json:"name"`
		OwnerUserId uint64       `json:"owner_user_id,string"`
	}

	TeamMember struct {
		MembershipState TeamMembershipState `json:"membership_state"`
		TeamId          uint64              `json:"team_id,string"`
		User            *user.User          `json:"user,omitempty"`
		Role            TeamMemberRole      `json:"role"`
	}

	TeamMembershipState int

	TeamMemberRole string
)

const (
	TeamMembershipStateInvited  TeamMembershipState = 1
	TeamMembershipStateAccepted TeamMembershipState = 2

	TeamMemberRoleAdmin     TeamMemberRole = "admin"
	TeamMemberRoleDeveloper TeamMemberRole = "developer"
	TeamMemberRoleReadOnly  TeamMemberRole = "read_only"
)
