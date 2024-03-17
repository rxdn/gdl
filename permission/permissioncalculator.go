package permission

import (
	"context"
	"errors"
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/sirupsen/logrus"
)

func HasPermissionsChannel(ctx context.Context, shard *gateway.Shard, guildId, userId, channelId uint64, permissions ...Permission) bool {
	sum, err := GetEffectivePermissionsChannel(ctx, shard, guildId, userId, channelId)
	if err != nil {
		return false
	}

	if HasPermissionRaw(sum, Administrator) {
		return true
	}

	hasPermission := true

	for _, permission := range permissions {
		if !HasPermissionRaw(sum, permission) {
			hasPermission = false
			break
		}
	}

	return hasPermission
}

func HasPermissions(ctx context.Context, shard *gateway.Shard, guildId, userId uint64, permissions ...Permission) bool {
	sum, err := GetEffectivePermissions(ctx, shard, guildId, userId)
	if err != nil {
		return false
	}

	if HasPermissionRaw(sum, Administrator) {
		return true
	}

	hasPermission := true

	for _, permission := range permissions {
		if !HasPermissionRaw(sum, permission) {
			hasPermission = false
			break
		}
	}

	return hasPermission
}

func GetAllPermissionsChannel(ctx context.Context, shard *gateway.Shard, guildId, userId, channelId uint64) []Permission {
	permissions := make([]Permission, 0)

	sum, err := GetEffectivePermissionsChannel(ctx, shard, guildId, userId, channelId)
	if err != nil {
		if shard.ShardManager.ShardOptions.Debug {
			logrus.Infof("shard %d: error retrieving permissions: %s", shard.ShardId, err.Error())
		}

		return permissions
	}

	for _, permission := range AllPermissions {
		if HasPermissionRaw(sum, permission) {
			permissions = append(permissions, permission)
		}
	}

	return permissions
}

func GetAllPermissions(ctx context.Context, shard *gateway.Shard, guildId, userId uint64) []Permission {
	permissions := make([]Permission, 0)

	sum, err := GetEffectivePermissions(ctx, shard, guildId, userId)
	if err != nil {
		if shard.ShardManager.ShardOptions.Debug {
			logrus.Infof("shard %d: error retrieving permissions: %s", shard.ShardId, err.Error())
		}

		return permissions
	}

	for _, permission := range AllPermissions {
		if HasPermissionRaw(sum, permission) {
			permissions = append(permissions, permission)
		}
	}

	return permissions
}

func GetEffectivePermissionsChannel(ctx context.Context, shard *gateway.Shard, guildId, userId, channelId uint64) (uint64, error) {
	permissions, err := GetBasePermissions(ctx, shard, guildId)
	if err != nil {
		return 0, err
	}

	permissions, err = GetGuildTotalRolePermissions(ctx, shard, guildId, userId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelBasePermissions(ctx, shard, guildId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelTotalRolePermissions(ctx, shard, guildId, userId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelMemberPermissions(ctx, shard, userId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	return permissions, nil
}

func GetEffectivePermissions(ctx context.Context, shard *gateway.Shard, guildId, userId uint64) (uint64, error) {
	permissions, err := GetBasePermissions(ctx, shard, guildId)
	if err != nil {
		return 0, err
	}

	permissions, err = GetGuildTotalRolePermissions(ctx, shard, guildId, userId, permissions)
	if err != nil {
		return 0, err
	}

	return permissions, nil
}

func GetChannelMemberPermissions(ctx context.Context, shard *gateway.Shard, userId, channelId uint64, initialPermissions uint64) (uint64, error) {
	ch, err := shard.GetChannel(ctx, channelId)
	if err != nil {
		return 0, err
	}

	for _, overwrite := range ch.PermissionOverwrites {
		if overwrite.Type == channel.PermissionTypeMember && overwrite.Id == userId {
			initialPermissions &= ^overwrite.Deny
			initialPermissions |= overwrite.Allow
		}
	}

	return initialPermissions, nil
}

func GetChannelTotalRolePermissions(ctx context.Context, shard *gateway.Shard, guildId, userId, channelId uint64, initialPermissions uint64) (uint64, error) {
	member, err := shard.GetGuildMember(ctx, guildId, userId)
	if err != nil {
		return 0, err
	}

	roles, err := shard.GetGuildRoles(ctx, guildId)
	if err != nil {
		return 0, err
	}

	ch, err := shard.GetChannel(ctx, channelId)
	if err != nil {
		return 0, err
	}

	var allow, deny uint64

	for _, memberRole := range member.Roles {
		for _, role := range roles {
			if memberRole == role.Id {
				for _, overwrite := range ch.PermissionOverwrites {
					if overwrite.Type == channel.PermissionTypeRole && overwrite.Id == role.Id {
						allow |= overwrite.Allow
						deny |= overwrite.Deny
						break
					}
				}
			}
		}
	}

	initialPermissions &= ^deny
	initialPermissions |= allow

	return initialPermissions, nil
}

func GetChannelBasePermissions(ctx context.Context, shard *gateway.Shard, guildId, channelId uint64, initialPermissions uint64) (uint64, error) {
	roles, err := shard.GetGuildRoles(ctx, guildId)
	if err != nil {
		return 0, err
	}

	var publicRole *guild.Role
	for _, role := range roles {
		if role.Id == guildId {
			publicRole = &role
			break
		}
	}

	if publicRole == nil {
		return 0, errors.New("couldn't find public role")
	}

	ch, err := shard.GetChannel(ctx, channelId)
	if err != nil {
		return 0, err
	}

	for _, overwrite := range ch.PermissionOverwrites {
		if overwrite.Type == channel.PermissionTypeRole && overwrite.Id == publicRole.Id {
			initialPermissions &= ^overwrite.Deny
			initialPermissions |= overwrite.Allow
			break
		}
	}

	return initialPermissions, nil
}

func GetGuildTotalRolePermissions(ctx context.Context, shard *gateway.Shard, guildId, userId uint64, initialPermissions uint64) (uint64, error) {
	member, err := shard.GetGuildMember(ctx, guildId, userId)
	if err != nil {
		return 0, err
	}

	roles, err := shard.GetGuildRoles(ctx, guildId)
	if err != nil {
		return 0, err
	}

	for _, memberRole := range member.Roles {
		for _, role := range roles {
			if memberRole == role.Id {
				initialPermissions |= role.Permissions
			}
		}
	}

	return initialPermissions, nil
}

func GetBasePermissions(ctx context.Context, shard *gateway.Shard, guildId uint64) (uint64, error) {
	roles, err := shard.GetGuildRoles(ctx, guildId)
	if err != nil {
		return 0, err
	}

	var publicRole *guild.Role
	for _, role := range roles {
		if role.Id == guildId {
			publicRole = &role
			break
		}
	}

	if publicRole == nil {
		return 0, errors.New("couldn't find public role")
	}

	return publicRole.Permissions, nil
}
