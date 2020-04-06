package permission

import (
	"errors"
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/sirupsen/logrus"
)

func HasPermissionsChannel(shard *gateway.Shard, guildId, userId, channelId uint64, permissions ...Permission) bool {
	sum, err := GetEffectivePermissionsChannel(shard, guildId, userId, channelId)
	if err != nil {
		return false
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

func HasPermissions(shard *gateway.Shard, guildId, userId uint64, permissions ...Permission) bool {
	sum, err := GetEffectivePermissions(shard, guildId, userId)
	if err != nil {
		return false
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

func GetAllPermissionsChannel(shard *gateway.Shard, guildId, userId, channelId uint64) []Permission {
	permissions := make([]Permission, 0)

	sum, err := GetEffectivePermissionsChannel(shard, guildId, userId, channelId)
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

func GetAllPermissions(shard *gateway.Shard, guildId, userId uint64) []Permission {
	permissions := make([]Permission, 0)

	sum, err := GetEffectivePermissions(shard, guildId, userId)
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

func GetEffectivePermissionsChannel(shard *gateway.Shard, guildId, userId, channelId uint64) (int, error) {
	permissions, err := GetBasePermissions(shard, guildId)
	if err != nil {
		return 0, err
	}

	permissions, err = GetGuildTotalRolePermissions(shard, guildId, userId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelBasePermissions(shard, guildId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelTotalRolePermissions(shard, guildId, userId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelMemberPermissions(shard, userId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	return permissions, nil
}

func GetEffectivePermissions(shard *gateway.Shard, guildId, userId uint64) (int, error) {
	permissions, err := GetBasePermissions(shard, guildId)
	if err != nil {
		return 0, err
	}

	permissions, err = GetGuildTotalRolePermissions(shard, guildId, userId, permissions)
	if err != nil {
		return 0, err
	}

	return permissions, nil
}

func GetChannelMemberPermissions(shard *gateway.Shard, userId, channelId uint64, initialPermissions int) (int, error) {
	ch, err := shard.GetChannel(channelId)
	if err != nil {
		return 0, err
	}

	for _, overwrite := range ch.PermissionOverwrites {
		if overwrite.Type == channel.PermissionTypeMember && overwrite.Id == userId {
			initialPermissions &= twosComplement(overwrite.Deny)
			initialPermissions |= overwrite.Allow
		}
	}

	return initialPermissions, nil
}

func GetChannelTotalRolePermissions(shard *gateway.Shard, guildId, userId, channelId uint64, initialPermissions int) (int, error) {
	member, err := shard.GetGuildMember(guildId, userId)
	if err != nil {
		return 0, err
	}

	roles, err := shard.GetGuildRoles(guildId)
	if err != nil {
		return 0, err
	}

	ch, err := shard.GetChannel(channelId)
	if err != nil {
		return 0, err
	}

	allow, deny := 0, 0

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

	initialPermissions &= twosComplement(deny)
	initialPermissions |= allow

	return initialPermissions, nil
}

func GetChannelBasePermissions(shard *gateway.Shard, guildId, channelId uint64, initialPermissions int) (int, error) {
	roles, err := shard.GetGuildRoles(guildId)
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

	ch, err := shard.GetChannel(channelId)
	if err != nil {
		return 0, err
	}

	for _, overwrite := range ch.PermissionOverwrites {
		if overwrite.Type == channel.PermissionTypeRole && overwrite.Id == publicRole.Id {
			initialPermissions &= twosComplement(overwrite.Deny)
			initialPermissions |= overwrite.Allow
			break
		}
	}

	return initialPermissions, nil
}

func GetGuildTotalRolePermissions(shard *gateway.Shard, guildId, userId uint64, initialPermissions int) (int, error) {
	member, err := shard.GetGuildMember(guildId, userId)
	if err != nil {
		return 0, err
	}

	roles, err := shard.GetGuildRoles(guildId)
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

func GetBasePermissions(shard *gateway.Shard, guildId uint64) (int, error) {
	roles, err := shard.GetGuildRoles(guildId)
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
