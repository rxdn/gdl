package permission

import "github.com/rxdn/gdl/gateway"

func CanSelfInteractWith(shard *gateway.Shard, guildId, targetId uint64) bool {
	return CanInteractWith(shard, (*shard.Cache).GetSelf().Id, guildId, targetId)
}

func CanInteractWith(shard *gateway.Shard, guildId, userId, targetId uint64) bool {
	return GetHighestRolePosition(shard, guildId, userId) > GetHighestRolePosition(shard, guildId, targetId)
}

func GetHighestRolePosition(shard *gateway.Shard, guildId, userId uint64) int {
	member, err := shard.GetGuildMember(guildId, userId); if err != nil {
		return 0
	}

	roles, err := shard.GetGuildRoles(guildId); if err != nil {
		return 0
	}

	highest := 0 // @everyone has a position of 0
	for _, roleId := range member.Roles {
		for _, role := range roles {
			if role.Id == roleId {
				if role.Position > highest {
					highest = role.Position
				}
			}
		}
	}

	return highest
}
