package cache

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
	"strconv"
	"sync"
)

type PgCache struct {
	*pgxpool.Pool
	options CacheOptions

	// TODO: Should we store self in the DB? Seems kinda redundant
	selfLock sync.RWMutex
	self     user.User
}

var json = jsoniter.Config{
	MarshalFloatWith6Digits: false,
	EscapeHTML:              false,
	SortMapKeys:             false,
}.Froze()

func NewPgCache(db *pgxpool.Pool, options CacheOptions) PgCache {
	return PgCache{
		Pool:    db,
		options: options,
	}
}

var (
	//go:embed sql/insert_user.sql
	queryInsertUser string
	//go:embed sql/get_user.sql
	queryGetUser string
	//go:embed sql/get_users.sql
	queryGetUsers string

	//go:embed sql/insert_guild.sql
	queryInsertGuild string
	//go:embed sql/get_guild.sql
	queryGetGuild string
	//go:embed sql/get_guild_count.sql
	queryGetGuildCount string
	//go:embed sql/get_guild_owner.sql
	queryGetGuildOwner string
	//go:embed sql/delete_guild.sql
	queryDeleteGuild string

	//go:embed sql/get_channel.sql
	queryGetChannel string
	//go:embed sql/get_guild_channels.sql
	queryGetGuildChannels string
	//go:embed sql/insert_channel.sql
	queryInsertChannel string
	//go:embed sql/delete_channel.sql
	queryDeleteChannel string
	//go:embed sql/delete_guild_channels.sql
	queryDeleteGuildChannels string

	//go:embed sql/get_role.sql
	queryGetRole string
	//go:embed sql/get_roles.sql
	queryGetRoles string
	//go:embed sql/get_guild_role.sql
	queryGetGuildRole string
	//go:embed sql/get_guild_roles.sql
	queryGetGuildRoles string
	//go:embed sql/insert_role.sql
	queryInsertRole string
	//go:embed sql/delete_role.sql
	queryDeleteRole string
	//go:embed sql/delete_guild_roles.sql
	queryDeleteGuildRoles string

	//go:embed sql/get_member.sql
	queryGetMember string
	//go:embed sql/get_guild_members.sql
	queryGetGuildMembers string
	//go:embed sql/get_guild_members_with_users.sql
	queryGetGuildMembersWithUsers string
	//go:embed sql/insert_member.sql
	queryInsertMember string
	//go:embed sql/delete_member.sql
	queryDeleteMember string

	//go:embed sql/get_emoji.sql
	queryGetEmoji string
	//go:embed sql/get_guild_emojis.sql
	queryGetGuildEmojis string
	//go:embed sql/insert_emoji.sql
	queryInsertEmoji string
	//go:embed sql/delete_emoji.sql
	queryDeleteEmoji string

	//go:embed sql/get_voice_state.sql
	queryGetVoiceState string
	//go:embed sql/get_guild_voice_states.sql
	queryGetGuildVoiceStates string
	//go:embed sql/insert_voice_state.sql
	queryInsertVoiceState string
	//go:embed sql/delete_voice_state.sql
	queryDeleteVoiceState string
)

func (c *PgCache) CreateSchema(ctx context.Context) error {
	batch := &pgx.Batch{}

	// create schema
	batch.Queue(`CREATE TABLE IF NOT EXISTS guilds("guild_id" int8 NOT NULL UNIQUE, "data" jsonb NOT NULL, PRIMARY KEY("guild_id"));`)

	batch.Queue(`CREATE TABLE IF NOT EXISTS channels("channel_id" int8 NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("channel_id", "guild_id"));`)
	batch.Queue(`CREATE TABLE IF NOT EXISTS users("user_id" int8 NOT NULL UNIQUE, "data" jsonb NOT NULL, PRIMARY KEY("user_id"));`)
	batch.Queue(`CREATE TABLE IF NOT EXISTS members("guild_id" int8 NOT NULL, "user_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("guild_id", "user_id"));`)
	batch.Queue(`CREATE TABLE IF NOT EXISTS roles("role_id" int8 NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("role_id", "guild_id"));`)
	batch.Queue(`CREATE TABLE IF NOT EXISTS emojis("emoji_id" int8 NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("emoji_id", "guild_id"));`)
	batch.Queue(`CREATE TABLE IF NOT EXISTS voice_states("guild_id" int8 NOT NULL, "user_id" INT8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("guild_id", "user_id"));`) // we may not have a cached user

	// create indexes
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS channels_guild_id ON channels("guild_id");`)
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS members_guild_id ON members("guild_id");`)
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS member_user_id ON members("user_id");`)
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS roles_guild_id ON roles("guild_id");`)
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS emojis_guild_id ON emojis("guild_id");`)
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS voice_states_guild_id ON voice_states("guild_id");`)
	batch.Queue(`CREATE INDEX CONCURRENTLY IF NOT EXISTS voice_states_user_id ON voice_states("user_id");`)

	_, err := c.SendBatch(ctx, batch).Exec()
	return err
}

func (c *PgCache) Options() CacheOptions {
	return c.options
}

func (c *PgCache) StoreUser(ctx context.Context, user user.User) error {
	if !c.options.Users {
		return nil
	}

	encoded, err := json.Marshal(user.ToCachedUser())
	if err != nil {
		return err
	}

	_, err = c.Exec(ctx, queryInsertUser, user.Id, string(encoded))
	return err
}

func (c *PgCache) StoreUsers(ctx context.Context, users []user.User) error {
	if !c.options.Users {
		return nil
	}

	batch := &pgx.Batch{}

	for _, u := range users {
		encoded, err := json.Marshal(u.ToCachedUser())
		if err != nil {
			return err
		}

		batch.Queue(queryInsertUser, u.Id, string(encoded))
	}

	br := c.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	return err
}

func (c *PgCache) GetUser(ctx context.Context, id uint64) (user.User, error) {
	var raw string
	if err := c.QueryRow(ctx, queryGetUser, id).Scan(&raw); err == pgx.ErrNoRows {
		return user.User{}, ErrNotFound
	} else if err != nil {
		return user.User{}, err
	}

	var cached user.CachedUser
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return user.User{}, err
	}

	return cached.ToUser(id), nil
}

func (c *PgCache) GetUsers(ctx context.Context, ids []uint64) (map[uint64]user.User, error) {
	idArray := &pgtype.Int8Array{}
	if err := idArray.Set(ids); err != nil {
		return nil, err
	}

	rows, err := c.Query(ctx, queryGetUsers, idArray)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := make(map[uint64]user.User)
	for rows.Next() {
		var id uint64
		var raw string
		if err := rows.Scan(&id, &raw); err != nil {
			return nil, err
		}

		var cached user.CachedUser
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return nil, err
		}

		users[id] = cached.ToUser(id)
	}

	return users, nil
}

func (c *PgCache) StoreGuilds(ctx context.Context, guilds []guild.Guild) error {
	if !c.options.Guilds {
		return nil
	}

	// store guilds
	batch := &pgx.Batch{}

	for _, guild := range guilds {
		// append guild
		encoded, err := json.Marshal(guild.ToCachedGuild())
		if err != nil {
			return err
		}

		batch.Queue(`INSERT INTO guilds("guild_id", "data") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "data" = $2;`, guild.Id, string(encoded))

		// append channels
		if c.options.Channels {
			for _, channel := range guild.Channels {
				encoded, err := json.Marshal(channel.ToCachedChannel())
				if err != nil {
					return err
				}

				batch.Queue(`INSERT INTO channels("channel_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("channel_id") DO UPDATE SET "data" = $3;`, channel.Id, channel.GuildId, string(encoded))
			}
		}

		// append roles
		if c.options.Roles {
			for _, role := range guild.Roles {
				encoded, err := json.Marshal(role.ToCachedRole(guild.Id))
				if err != nil {
					return err
				}

				batch.Queue(`INSERT INTO roles("role_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("role_id", "guild_id") DO UPDATE SET "data" = $3;`, role.Id, guild.Id, string(encoded))
			}
		}

		// append emojis
		if c.options.Emojis {
			for _, emoji := range guild.Emojis {
				encoded, err := json.Marshal(emoji.ToCachedEmoji(guild.Id))
				if err != nil {
					return err
				}

				batch.Queue(`INSERT INTO emojis("emoji_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("emoji_id") DO UPDATE SET "data" = $3;`, emoji.Id, guild.Id, string(encoded))
			}
		}

		// append voice states
		if c.options.VoiceStates {
			for _, state := range guild.VoiceStates {
				encoded, err := json.Marshal(state.ToCachedVoiceState())
				if err != nil {
					return err
				}

				batch.Queue(`INSERT INTO voice_states("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, state.GuildId, state.UserId, string(encoded))
			}
		}
	}

	br := c.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	return err
}

func (c *PgCache) StoreGuild(ctx context.Context, g guild.Guild) error {
	if c.options.Guilds {
		encoded, err := json.Marshal(g.ToCachedGuild())
		if err != nil {
			return err
		}

		if _, err := c.Exec(context.Background(), queryInsertGuild, g.Id, string(encoded)); err != nil {
			return err
		}
	}

	for i, channel := range g.Channels {
		channel.GuildId = g.Id
		g.Channels[i] = channel
	}

	if err := c.StoreChannels(ctx, g.Channels); err != nil {
		return err
	}

	if err := c.StoreRoles(ctx, g.Roles, g.Id); err != nil {
		return err
	}

	if err := c.StoreEmojis(ctx, g.Emojis, g.Id); err != nil {
		return err
	}

	if err := c.StoreVoiceStates(ctx, g.VoiceStates); err != nil {
		return err
	}

	return nil
}

func (c *PgCache) GetGuild(ctx context.Context, id uint64) (guild.Guild, error) {
	var raw string
	err := c.QueryRow(context.Background(), queryGetGuild, id).Scan(&raw)
	if err == pgx.ErrNoRows {
		return guild.Guild{}, ErrNotFound
	} else if err != nil {
		return guild.Guild{}, err
	}

	var cachedGuild guild.CachedGuild
	if err := json.Unmarshal([]byte(raw), &cachedGuild); err != nil {
		return guild.Guild{}, err
	}

	return cachedGuild.ToGuild(id), nil
}

func (c *PgCache) GetGuildChannels(ctx context.Context, guildId uint64) ([]channel.Channel, error) {
	if !c.options.Channels {
		return nil, nil
	}

	rows, err := c.Query(ctx, queryGetGuildChannels, guildId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var channels []channel.Channel
	for rows.Next() {
		var channelId uint64
		var raw string

		if err := rows.Scan(&channelId, &raw); err != nil {
			return nil, err
		}

		var cached channel.CachedChannel
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return nil, err
		}

		channels = append(channels, cached.ToChannel(channelId, guildId))
	}

	return channels, nil
}

func (c *PgCache) GetGuildRoles(ctx context.Context, guildId uint64) ([]guild.Role, error) {
	if !c.options.Roles {
		return nil, nil
	}

	rows, err := c.Query(ctx, queryGetGuildRoles, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var roles []guild.Role
	for rows.Next() {
		var roleId uint64
		var raw string

		if err := rows.Scan(&roleId, &raw); err != nil {
			return nil, err
		}

		var cached guild.CachedRole
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return nil, err
		}

		roles = append(roles, cached.ToRole(roleId))
	}

	return roles, nil
}

func (c *PgCache) GetGuildMembers(ctx context.Context, guildId uint64, withUserData bool) ([]member.Member, error) {
	if !c.options.Members {
		return nil, nil
	}

	var query string
	if withUserData {
		query = queryGetGuildMembersWithUsers
	} else {
		query = queryGetGuildMembers
	}

	rows, err := c.Query(ctx, query, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var members []member.Member
	for rows.Next() {
		var userId uint64
		var memberRaw, userRaw string

		if withUserData {
			if err := rows.Scan(&userId, &memberRaw, &userRaw); err != nil {
				return nil, err
			}
		} else {
			if err := rows.Scan(&userId, &memberRaw); err != nil {
				return nil, err
			}
		}

		var cachedMember member.CachedMember
		if err := json.Unmarshal([]byte(memberRaw), &cachedMember); err != nil {
			return nil, err
		}

		var cachedUser user.CachedUser
		if withUserData && userRaw != "" {
			if err := json.Unmarshal([]byte(userRaw), &cachedUser); err != nil {
				return nil, err
			}
		}

		members = append(members, cachedMember.ToMember(cachedUser.ToUser(userId)))
	}

	return members, nil
}

func (c *PgCache) GetGuildEmojis(ctx context.Context, guildId uint64) ([]emoji.Emoji, error) {
	if !c.options.Emojis {
		return nil, nil
	}

	rows, err := c.Query(ctx, queryGetGuildEmojis, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var emojis []emoji.Emoji
	for rows.Next() {
		var emojiId uint64
		var raw string
		if err := rows.Scan(&emojiId, &raw); err != nil {
			return nil, err
		}

		var cached emoji.CachedEmoji
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return nil, err
		}

		emojis = append(emojis, cached.ToEmoji(emojiId, user.User{}))
	}

	return emojis, nil
}

func (c *PgCache) DeleteGuild(ctx context.Context, id uint64) error {
	_, err := c.Exec(ctx, queryDeleteGuild, id)
	return err
}

func (c *PgCache) GetGuildCount(ctx context.Context) (int, error) {
	var count int
	if err := c.QueryRow(ctx, queryGetGuildCount).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (c *PgCache) GetGuildOwner(ctx context.Context, guildId uint64) (uint64, error) {
	var ownerId string
	if err := c.QueryRow(ctx, queryGetGuildOwner, guildId).Scan(&ownerId); err == pgx.ErrNoRows {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	parsed, err := strconv.ParseUint(ownerId, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func (c *PgCache) StoreMember(ctx context.Context, m member.Member, guildId uint64) error {
	if !c.options.Members {
		return nil
	}

	encoded, err := json.Marshal(m.ToCachedMember())
	if err != nil {
		return err
	}

	_, err = c.Exec(ctx, queryInsertMember, guildId, m.User.Id, string(encoded))
	return err
}

func (c *PgCache) StoreMembers(ctx context.Context, members []member.Member, guildId uint64) error {
	if c.options.Members {
		return nil
	}

	batch := &pgx.Batch{}

	for _, m := range members {
		encoded, err := json.Marshal(m.ToCachedMember())
		if err != nil {
			return err
		}

		batch.Queue(queryInsertMember, guildId, m.User.Id, string(encoded))
	}

	br := c.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	return err
}

func (c *PgCache) GetMember(ctx context.Context, guildId, userId uint64) (member.Member, error) {
	if !c.options.Members {
		return member.Member{}, ErrNotFound
	}

	var memberRaw, userRaw sql.NullString
	if err := c.QueryRow(ctx, queryGetMember, guildId, userId).Scan(&memberRaw, &userRaw); err == pgx.ErrNoRows {
		return member.Member{}, ErrNotFound
	} else if err != nil {
		return member.Member{}, err
	}

	if !memberRaw.Valid {
		return member.Member{}, ErrNotFound
	}

	var cachedMember member.CachedMember
	if err := json.Unmarshal([]byte(memberRaw.String), &cachedMember); err != nil {
		return member.Member{}, err
	}

	var cachedUser user.CachedUser
	if userRaw.Valid {
		if err := json.Unmarshal([]byte(userRaw.String), &cachedUser); err != nil {
			return member.Member{}, err
		}
	}

	return cachedMember.ToMember(cachedUser.ToUser(userId)), nil
}

func (c *PgCache) DeleteMember(ctx context.Context, userId, guildId uint64) error {
	_, err := c.Exec(ctx, queryDeleteMember, guildId, userId)
	return err
}

func (c *PgCache) StoreChannel(ctx context.Context, ch channel.Channel) error {
	if !c.options.Channels {
		return nil
	}

	encoded, err := json.Marshal(ch.ToCachedChannel())
	if err != nil {
		return err
	}

	_, err = c.Exec(ctx, queryInsertChannel, ch.Id, ch.GuildId, string(encoded))
	return err
}

func (c *PgCache) StoreChannels(ctx context.Context, channels []channel.Channel) error {
	if !c.options.Channels {
		return nil
	}

	batch := &pgx.Batch{}

	for _, ch := range channels {
		encoded, err := json.Marshal(ch.ToCachedChannel())
		if err != nil {
			return err
		}

		batch.Queue(queryInsertChannel, ch.Id, ch.GuildId, string(encoded))
	}

	br := c.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	return err
}

func (c *PgCache) ReplaceChannels(ctx context.Context, guildId uint64, channels []channel.Channel) error {
	if !c.options.Channels {
		return nil
	}

	tx, err := c.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queryDeleteGuildChannels, guildId); err != nil {
		return err
	}

	for _, ch := range channels {
		encoded, err := json.Marshal(ch.ToCachedChannel())
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, queryInsertChannel, ch.Id, ch.GuildId, string(encoded)); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (c *PgCache) GetChannel(ctx context.Context, channelId uint64) (channel.Channel, error) {
	if !c.options.Channels {
		return channel.Channel{}, ErrNotFound
	}

	var guildId uint64
	var raw string
	if err := c.QueryRow(ctx, queryGetChannel, channelId).Scan(&guildId, &raw); err == pgx.ErrNoRows {
		return channel.Channel{}, ErrNotFound
	} else if err != nil {
		return channel.Channel{}, err
	}

	var cached channel.CachedChannel
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return channel.Channel{}, err
	}

	return cached.ToChannel(channelId, guildId), nil
}

func (c *PgCache) DeleteChannel(ctx context.Context, channelId uint64) error {
	return utils.Second(c.Exec(ctx, queryDeleteChannel, channelId))
}

func (c *PgCache) DeleteGuildChannels(ctx context.Context, guildId uint64) error {
	return utils.Second(c.Exec(ctx, queryDeleteGuildChannels, guildId))
}

func (c *PgCache) StoreRole(ctx context.Context, role guild.Role, guildId uint64) error {
	if !c.options.Roles {
		return nil
	}

	encoded, err := json.Marshal(role.ToCachedRole(guildId))
	if err != nil {
		return err
	}

	_, err = c.Exec(ctx, queryInsertRole, role.Id, guildId, string(encoded))
	return err
}

func (c *PgCache) StoreRoles(ctx context.Context, roles []guild.Role, guildId uint64) error {
	if !c.options.Roles {
		return nil
	}

	batch := &pgx.Batch{}

	for _, role := range roles {
		encoded, err := json.Marshal(role.ToCachedRole(guildId))
		if err != nil {
			return err
		}

		batch.Queue(queryInsertRole, role.Id, guildId, string(encoded))
	}

	br := c.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	return err
}

func (c *PgCache) GetRole(ctx context.Context, id uint64) (guild.Role, error) {
	if !c.options.Roles {
		return guild.Role{}, ErrNotFound
	}

	var raw string
	if err := c.QueryRow(ctx, queryGetRole, id).Scan(&raw); err == pgx.ErrNoRows {
		return guild.Role{}, ErrNotFound
	} else if err != nil {
		return guild.Role{}, err
	}

	var cached guild.CachedRole
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return guild.Role{}, err
	}

	return cached.ToRole(id), nil
}

func (c *PgCache) GetRoles(ctx context.Context, guildId uint64, ids []uint64) (map[uint64]guild.Role, error) {
	if !c.options.Roles {
		return nil, ErrNotFound
	}

	idArray := &pgtype.Int8Array{}
	if err := idArray.Set(ids); err != nil {
		return nil, err
	}

	rows, err := c.Query(ctx, queryGetRoles, idArray, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := make(map[uint64]guild.Role)
	for rows.Next() {
		var id uint64
		var raw string
		if err := rows.Scan(&id, &raw); err != nil {
			return nil, err
		}

		var cached guild.CachedRole
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return nil, err
		}

		roles[id] = cached.ToRole(id)
	}

	return roles, nil
}

func (c *PgCache) GetGuildRole(ctx context.Context, id, guildId uint64) (guild.Role, error) {
	if !c.options.Roles {
		return guild.Role{}, ErrNotFound
	}

	var raw string
	if err := c.QueryRow(ctx, queryGetGuildRole, id, guildId).Scan(&raw); err == pgx.ErrNoRows {
		return guild.Role{}, ErrNotFound
	} else if err != nil {
		return guild.Role{}, err
	}

	var cached guild.CachedRole
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return guild.Role{}, err
	}

	return cached.ToRole(id), nil
}

func (c *PgCache) DeleteRole(ctx context.Context, roleId uint64) error {
	_, err := c.Exec(ctx, queryDeleteRole, roleId)
	return err
}

func (c *PgCache) DeleteGuildRoles(ctx context.Context, guildId uint64) error {
	_, err := c.Exec(ctx, queryDeleteGuildRoles, guildId)
	return err
}

func (c *PgCache) StoreEmoji(ctx context.Context, emoji emoji.Emoji, guildId uint64) error {
	if !c.options.Emojis {
		return nil
	}

	encoded, err := json.Marshal(emoji.ToCachedEmoji(guildId))
	if err != nil {
		return err
	}

	_, err = c.Exec(ctx, queryInsertEmoji, emoji.Id, guildId, string(encoded))
	return err
}

func (c *PgCache) StoreEmojis(ctx context.Context, emojis []emoji.Emoji, guildId uint64) error {
	if !c.options.Emojis {
		return nil
	}

	conversionFunc := func(item emoji.Emoji) emoji.CachedEmoji {
		return item.ToCachedEmoji(guildId)
	}

	return batchStore(ctx, c, queryInsertEmoji, emojis, conversionFunc, func(item emoji.Emoji, encoded string) []interface{} {
		return []interface{}{item.Id, guildId, encoded}
	})
}

func (c *PgCache) GetEmoji(ctx context.Context, id uint64) (emoji.Emoji, error) {
	if !c.options.Emojis {
		return emoji.Emoji{}, ErrNotFound
	}

	var raw string
	if err := c.QueryRow(ctx, queryGetEmoji, id).Scan(&raw); err == pgx.ErrNoRows {
		return emoji.Emoji{}, ErrNotFound
	} else if err != nil {
		return emoji.Emoji{}, err
	}

	var cached emoji.CachedEmoji
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return emoji.Emoji{}, err
	}

	return cached.ToEmoji(id, user.User{}), nil
}

func (c *PgCache) DeleteEmoji(ctx context.Context, emojiId uint64) error {
	_, err := c.Exec(ctx, queryDeleteEmoji, emojiId)
	return err
}

func (c *PgCache) StoreVoiceState(ctx context.Context, state guild.VoiceState) error {
	if !c.options.VoiceStates {
		return nil
	}

	encoded, err := json.Marshal(state.ToCachedVoiceState())
	if err != nil {
		return err
	}

	_, err = c.Exec(ctx, queryInsertVoiceState, state.GuildId, state.UserId, string(encoded))
	return err
}

func (c *PgCache) StoreVoiceStates(ctx context.Context, states []guild.VoiceState) error {
	if !c.options.VoiceStates {
		return nil
	}

	conversionFunc := func(state guild.VoiceState) guild.CachedVoiceState {
		return state.ToCachedVoiceState()
	}

	return batchStore(ctx, c, queryInsertVoiceState, states, conversionFunc, func(state guild.VoiceState, encoded string) []interface{} {
		return []interface{}{state.GuildId, state.UserId, encoded}
	})
}

func (c *PgCache) GetVoiceState(ctx context.Context, userId, guildId uint64) (guild.VoiceState, error) {
	if !c.options.VoiceStates {
		return guild.VoiceState{}, ErrNotFound
	}

	var voiceStateRaw, memberRaw, userRaw string
	if err := c.QueryRow(ctx, queryGetVoiceState, guildId, userId).Scan(&voiceStateRaw, &memberRaw, &userRaw); err == pgx.ErrNoRows {
		return guild.VoiceState{}, ErrNotFound
	} else if err != nil {
		return guild.VoiceState{}, err
	}

	var cachedVoiceState guild.CachedVoiceState
	if err := json.Unmarshal([]byte(voiceStateRaw), &cachedVoiceState); err != nil {
		return guild.VoiceState{}, err
	}

	var cachedMember member.CachedMember
	if len(memberRaw) > 0 {
		if err := json.Unmarshal([]byte(memberRaw), &cachedMember); err != nil {
			return guild.VoiceState{}, err
		}
	}

	var cachedUser user.CachedUser
	if len(memberRaw) > 0 {
		if err := json.Unmarshal([]byte(userRaw), &cachedUser); err != nil {
			return guild.VoiceState{}, err
		}
	}

	return cachedVoiceState.ToVoiceState(guildId, cachedMember.ToMember(cachedUser.ToUser(userId))), nil
}

func (c *PgCache) GetGuildVoiceStates(ctx context.Context, guildId uint64) ([]guild.VoiceState, error) {
	if !c.options.VoiceStates {
		return nil, nil
	}

	rows, err := c.Query(ctx, queryGetGuildVoiceStates, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var states []guild.VoiceState
	for rows.Next() {
		var userId uint64
		var voiceStateRaw, memberRaw, userRaw string

		if err := rows.Scan(&userId, &voiceStateRaw, &memberRaw, &userRaw); err != nil {
			return nil, err
		}

		var cachedVoiceState guild.CachedVoiceState
		if err := json.Unmarshal([]byte(voiceStateRaw), &cachedVoiceState); err != nil {
			return nil, err
		}

		var cachedMember member.CachedMember
		if len(memberRaw) > 0 {
			if err := json.Unmarshal([]byte(memberRaw), &cachedMember); err != nil {
				return nil, err
			}
		}

		var cachedUser user.CachedUser
		if len(memberRaw) > 0 {
			if err := json.Unmarshal([]byte(userRaw), &cachedUser); err != nil {
				return nil, err
			}
		}

		states = append(states, cachedVoiceState.ToVoiceState(userId, cachedMember.ToMember(cachedUser.ToUser(userId))))
	}

	return states, nil
}

func (c *PgCache) DeleteVoiceState(ctx context.Context, userId, guildId uint64) error {
	_, err := c.Exec(ctx, queryDeleteVoiceState, userId, guildId)
	return err
}

func (c *PgCache) StoreSelf(ctx context.Context, self user.User) error {
	c.selfLock.Lock()
	c.self = self
	c.selfLock.Unlock()
	return nil
}

func (c *PgCache) GetSelf(ctx context.Context) (user.User, error) {
	c.selfLock.RLock()
	self := c.self
	c.selfLock.RUnlock()

	if self.Id == 0 {
		return user.User{}, ErrNotFound
	}

	return self, nil
}

func batchStore[T any, U any](
	ctx context.Context,
	c *PgCache,
	query string,
	items []T,
	convertFunc func(T) U,
	argFunc func(item T, encoded string) []interface{},
) error {
	batch := &pgx.Batch{}

	for _, item := range items {
		encoded, err := json.Marshal(convertFunc(item))
		if err != nil {
			return err
		}

		batch.Queue(query, argFunc(item, string(encoded))...)
	}

	br := c.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	return err
}
