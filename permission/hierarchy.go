package permission

import (
	"context"
	"errors"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/gateway"
)

func CanSelfInteractWith(ctx context.Context, shard *gateway.Shard, guildId, targetId uint64) (bool, error) {
	self, err := shard.Cache.GetSelf(ctx)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return CanInteractWith(ctx, shard, self.Id, guildId, targetId), nil
}

func CanInteractWith(ctx context.Context, shard *gateway.Shard, guildId, userId, targetId uint64) bool {
	return GetHighestRolePosition(ctx, shard, guildId, userId) > GetHighestRolePosition(ctx, shard, guildId, targetId)
}

func GetHighestRolePosition(ctx context.Context, shard *gateway.Shard, guildId, userId uint64) int {
	member, err := shard.GetGuildMember(ctx, guildId, userId)
	if err != nil {
		return 0
	}

	roles, err := shard.GetGuildRoles(ctx, guildId)
	if err != nil {
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
