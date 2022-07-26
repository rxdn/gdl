package cache

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"strconv"
	"sync"
)

type PgCache struct {
	*pgxpool.Pool
	Options CacheOptions

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
	// create schema
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS guilds("guild_id" int8 NOT NULL UNIQUE, "data" jsonb NOT NULL, PRIMARY KEY("guild_id"));`)
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS channels("channel_id" int8 NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("channel_id", "guild_id"));`)
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS users("user_id" int8 NOT NULL UNIQUE, "data" jsonb NOT NULL, PRIMARY KEY("user_id"));`)
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS members("guild_id" int8 NOT NULL, "user_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("guild_id", "user_id"));`)
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS roles("role_id" int8 NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("role_id", "guild_id"));`)
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS emojis("emoji_id" int8 NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("emoji_id", "guild_id"));`)
	pgMustRun(db, `CREATE TABLE IF NOT EXISTS voice_states("guild_id" int8 NOT NULL, "user_id" INT8 NOT NULL, "data" jsonb NOT NULL, PRIMARY KEY("guild_id", "user_id"));`) // we may not have a cached user

	// create indexes
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS channels_guild_id ON channels("guild_id");`)
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS members_guild_id ON members("guild_id");`)
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS member_user_id ON members("user_id");`)
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS roles_guild_id ON roles("guild_id");`)
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS emojis_guild_id ON emojis("guild_id");`)
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS voice_states_guild_id ON voice_states("guild_id");`)
	pgMustRun(db, `CREATE INDEX CONCURRENTLY IF NOT EXISTS voice_states_user_id ON voice_states("user_id");`)

	return PgCache{
		Pool:    db,
		Options: options,
	}
}

func pgMustRun(db *pgxpool.Pool, query string) {
	if _, err := db.Exec(context.Background(), query); err != nil {
		panic(err)
	}
}

func (c *PgCache) GetOptions() CacheOptions {
	return c.Options
}

func (c *PgCache) StoreUser(user user.User) {
	if c.Options.Users {
		if encoded, err := json.Marshal(user.ToCachedUser()); err == nil {
			_, _ = c.Exec(context.Background(), `INSERT INTO users("user_id", "data") VALUES($1, $2) ON CONFLICT("user_id") DO UPDATE SET "data" = $2;`, user.Id, string(encoded))
		}
	}
}

func (c *PgCache) StoreUsers(users []user.User) {
	if c.Options.Users {
		batch := &pgx.Batch{}

		for _, u := range users {
			if encoded, err := json.Marshal(u.ToCachedUser()); err == nil {
				batch.Queue(`INSERT INTO users("user_id", "data") VALUES($1, $2) ON CONFLICT("user_id") DO UPDATE SET "data" = $2;`, u.Id, string(encoded))
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()

		_, _ = br.Exec()
	}
}

func (c *PgCache) GetUser(id uint64) (u user.User, ok bool) {
	var raw string
	if err := c.QueryRow(context.Background(), `SELECT "data" FROM users WHERE "user_id" = $1;`, id).Scan(&raw); err != nil {
		return
	}

	var cached user.CachedUser
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return
	}

	return cached.ToUser(id), true
}

func (c *PgCache) GetUsers(ids []uint64) (map[uint64]user.User, error) {
	idArray := &pgtype.Int8Array{}
	if err := idArray.Set(ids); err != nil {
		return nil, err
	}

	rows, err := c.Query(context.Background(), `SELECT "user_id", "data" FROM users WHERE "user_id" = ANY($1);`, idArray)
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

// TODO: The "data" field just has null values. Find the cause and solution.
func (c *PgCache) StoreGuilds(guilds []guild.Guild) {
	if c.Options.Guilds {
		// store guilds
		batch := &pgx.Batch{}

		for _, guild := range guilds {
			// append guild
			if encoded, err := json.Marshal(guild.ToCachedGuild()); err == nil {
				batch.Queue(`INSERT INTO guilds("guild_id", "data") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "data" = $2;`, guild.Id, string(encoded))
			}

			// append channels
			if c.Options.Channels {
				for _, channel := range guild.Channels {
					if encoded, err := json.Marshal(channel.ToCachedChannel()); err == nil {
						batch.Queue(`INSERT INTO channels("channel_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("channel_id") DO UPDATE SET "data" = $3;`, channel.Id, channel.GuildId, string(encoded))
					}
				}
			}

			// append roles
			if c.Options.Roles {
				for _, role := range guild.Roles {
					if encoded, err := json.Marshal(role.ToCachedRole(guild.Id)); err == nil {
						batch.Queue(`INSERT INTO roles("role_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("role_id", "guild_id") DO UPDATE SET "data" = $3;`, role.Id, guild.Id, string(encoded))
					}
				}
			}

			// append members
			if c.Options.Members {
				for _, member := range guild.Members {
					if encoded, err := json.Marshal(member.ToCachedMember()); err == nil {
						batch.Queue(`INSERT INTO members("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, guild.Id, member.User.Id, string(encoded))
					}
				}
			}

			// append users
			if c.Options.Users {
				for _, member := range guild.Members {
					if encoded, err := json.Marshal(member.User.ToCachedUser()); err == nil {
						batch.Queue(`INSERT INTO users("user_id", "data") VALUES($1, $2) ON CONFLICT("user_id") DO UPDATE SET "data" = $2;`, member.User.Id, string(encoded))
					}
				}
			}

			// append emojis
			if c.Options.Emojis {
				for _, emoji := range guild.Emojis {
					if encoded, err := json.Marshal(emoji.ToCachedEmoji(guild.Id)); err == nil {
						batch.Queue(`INSERT INTO emojis("emoji_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("emoji_id") DO UPDATE SET "data" = $3;`, emoji.Id, guild.Id, string(encoded))
					}
				}
			}

			// append voice states
			if c.Options.VoiceStates {
				for _, state := range guild.VoiceStates {
					if encoded, err := json.Marshal(state.ToCachedVoiceState()); err == nil {
						batch.Queue(`INSERT INTO voice_states("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, state.GuildId, state.UserId, string(encoded))
					}
				}
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()
		_, _ = br.Exec()
	}
}

func (c *PgCache) StoreGuild(g guild.Guild) {
	if c.Options.Guilds {
		if encoded, err := json.Marshal(g.ToCachedGuild()); err == nil {
			_, _ = c.Exec(context.Background(), `INSERT INTO guilds("guild_id", "data") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "data" = $2;`, g.Id, string(encoded))
		}
	}
	for i, channel := range g.Channels {
		channel.GuildId = g.Id
		g.Channels[i] = channel
	}

	c.StoreChannels(g.Channels)
	c.StoreRoles(g.Roles, g.Id)
	c.StoreMembers(g.Members, g.Id)
	c.StoreEmojis(g.Emojis, g.Id)
	c.StoreVoiceStates(g.VoiceStates)

	var users []user.User
	for _, m := range g.Members {
		users = append(users, m.User)
	}
	c.StoreUsers(users)
}

// use withMembers with extreme caution!
func (c *PgCache) GetGuild(id uint64, withMembers bool) (g guild.Guild, ok bool) {
	var raw string
	err := c.QueryRow(context.Background(), `SELECT "data" FROM guilds WHERE "guild_id" = $1;`, id).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			return g, false
		}

		fmt.Println(err.Error())
		return
	}

	var cachedGuild guild.CachedGuild
	if err := json.Unmarshal([]byte(raw), &cachedGuild); err != nil {
		fmt.Println(err.Error())
		return
	}

	g = cachedGuild.ToGuild(id)

	g.Channels = c.GetGuildChannels(id)
	g.Roles = c.GetGuildRoles(id)

	if withMembers {
		g.Members = c.GetGuildMembers(id, false)
	}

	g.Emojis = c.GetGuildEmojis(id)
	g.VoiceStates = c.GetGuildVoiceStates(id)

	return g, true
}

func (c *PgCache) GetGuildChannels(guildId uint64) (channels []channel.Channel) {
	if !c.Options.Channels {
		return
	}

	rows, err := c.Query(context.Background(), `SELECT "channel_id", "data" FROM channels WHERE "guild_id" = $1;`, guildId)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var channelId uint64
		var raw string

		if err := rows.Scan(&channelId, &raw); err != nil {
			continue
		}

		var cached channel.CachedChannel
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return
		}

		channels = append(channels, cached.ToChannel(channelId, guildId))
	}

	return
}

func (c *PgCache) GetGuildRoles(guildId uint64) (roles []guild.Role) {
	if !c.Options.Roles {
		return
	}

	rows, err := c.Query(context.Background(), `SELECT "role_id", "data" FROM roles WHERE "guild_id" = $1;`, guildId)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var roleId uint64
		var raw string

		if err := rows.Scan(&roleId, &raw); err != nil {
			continue
		}

		var cached guild.CachedRole
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return
		}

		roles = append(roles, cached.ToRole(roleId))
	}

	return
}

func (c *PgCache) GetGuildMembers(guildId uint64, withUserData bool) (members []member.Member) {
	if !c.Options.Members {
		return
	}

	var query string
	if withUserData {
		query = `
SELECT "members.user_id", "members.data", "users.data"
FROM members
LEFT JOIN users ON "members.user_id"="users.user_id"
WHERE "guild_id" = $1;
`
	} else {
		query = `SELECT "user_id", "data" FROM members WHERE "guild_id" = $1;`
	}

	rows, err := c.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var userId uint64
		var memberRaw, userRaw string

		var err error
		if withUserData {
			err = rows.Scan(&userId, &memberRaw, &userRaw)
		} else {
			err = rows.Scan(&userId, &memberRaw)
		}

		if err != nil {
			continue
		}

		var cachedMember member.CachedMember
		if err := json.Unmarshal([]byte(memberRaw), &cachedMember); err != nil {
			return
		}

		var cachedUser user.CachedUser
		if withUserData && userRaw != "" {
			if err := json.Unmarshal([]byte(userRaw), &cachedUser); err != nil {
				return
			}
		}

		member := cachedMember.ToMember(cachedUser.ToUser(userId))
		members = append(members, member)
	}

	return
}

func (c *PgCache) GetGuildEmojis(guildId uint64) (emojis []emoji.Emoji) {
	if !c.Options.Emojis {
		return
	}

	rows, err := c.Query(context.Background(), `SELECT "emoji_id", "data" FROM emojis WHERE "guild_id" = $1;`, guildId)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var emojiId uint64
		var raw string

		if err := rows.Scan(&emojiId, &raw); err != nil {
			continue
		}

		var cached emoji.CachedEmoji
		if err := json.Unmarshal([]byte(raw), &cached); err != nil {
			return
		}

		emojis = append(emojis, cached.ToEmoji(emojiId, user.User{}))
	}

	return
}

// TODO: FIX
func (c *PgCache) GetGuilds() []guild.Guild {
	return nil
}

func (c *PgCache) DeleteGuild(id uint64) {
	if c.Options.Guilds {
		_, _ = c.Exec(context.Background(), `DELETE FROM guilds WHERE "guild_id" = $1;`, id)
	}
}

func (c *PgCache) GetGuildCount() int {
	var count int
	_ = c.QueryRow(context.Background(), "SELECT COUNT(*) FROM guilds;").Scan(&count)
	return count
}

func (c *PgCache) GetGuildOwner(guildId uint64) (uint64, bool) {
	query := `SELECT data->'owner_id' FROM guilds WHERE "guild_id" = $1;`

	var ownerId string
	if err := c.QueryRow(context.Background(), query, guildId).Scan(&ownerId); err != nil { // Includes pgx.ErrNoRows
		return 0, false
	}

	parsed, err := strconv.ParseUint(ownerId, 10, 64)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

func (c *PgCache) StoreMember(m member.Member, guildId uint64) {
	if c.Options.Members {
		if encoded, err := json.Marshal(m.ToCachedMember()); err == nil {
			_, _ = c.Exec(context.Background(), `INSERT INTO members("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, guildId, m.User.Id, string(encoded))
		}
	}
}

func (c *PgCache) StoreMembers(members []member.Member, guildId uint64) {
	if c.Options.Members {
		batch := &pgx.Batch{}

		for _, m := range members {
			if encoded, err := json.Marshal(m.ToCachedMember()); err == nil {
				batch.Queue(`INSERT INTO members("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, guildId, m.User.Id, string(encoded))
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()
		_, _ = br.Exec()
	}
}

func (c *PgCache) GetMember(guildId, userId uint64) (member.Member, bool) {
	if !c.Options.Members {
		return member.Member{}, false
	}

	query := `
SELECT members.data, users.data FROM members
LEFT JOIN users ON members.user_id=users.user_id
WHERE "guild_id" = $1 AND members.user_id = $2;
`

	var memberRaw, userRaw sql.NullString
	if err := c.QueryRow(context.Background(), query, guildId, userId).Scan(&memberRaw, &userRaw); err != nil {
		return member.Member{}, false
	}

	// we need to cache either member or user
	if !memberRaw.Valid || !userRaw.Valid {
		return member.Member{}, false
	}

	var cachedMember member.CachedMember
	if err := json.Unmarshal([]byte(memberRaw.String), &cachedMember); err != nil {
		return member.Member{}, false
	}

	var cachedUser user.CachedUser
	if err := json.Unmarshal([]byte(userRaw.String), &cachedUser); err != nil {
		return member.Member{}, false
	}

	return cachedMember.ToMember(cachedUser.ToUser(userId)), true
}

func (c *PgCache) DeleteMember(userId, guildId uint64) {
	if c.Options.Members {
		_, _ = c.Exec(context.Background(), `DELETE FROM members WHERE "guild_id" = $1 AND "user_id" = $2;`, guildId, userId)
	}
}

func (c *PgCache) StoreChannel(ch channel.Channel) {
	if c.Options.Channels {
		if encoded, err := json.Marshal(ch.ToCachedChannel()); err == nil {
			_, err = c.Exec(context.Background(), `INSERT INTO channels("channel_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("channel_id") DO UPDATE SET "data" = $3;`, ch.Id, ch.GuildId, string(encoded))
		}
	}
}

func (c *PgCache) StoreChannels(channels []channel.Channel) {
	if c.Options.Channels {
		batch := &pgx.Batch{}

		for _, ch := range channels {
			if encoded, err := json.Marshal(ch.ToCachedChannel()); err == nil {
				batch.Queue(`INSERT INTO channels("channel_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("channel_id") DO UPDATE SET "data" = $3;`, ch.Id, ch.GuildId, string(encoded))
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()
		_, _ = br.Exec()
	}
}

func (c *PgCache) GetChannel(channelId uint64) (channel.Channel, bool) {
	if !c.Options.Channels {
		return channel.Channel{}, false
	}

	var guildId uint64
	var raw string
	if err := c.QueryRow(context.Background(), `SELECT "guild_id", "data" FROM channels WHERE "channel_id" = $1;`, channelId).Scan(&guildId, &raw); err != nil {
		return channel.Channel{}, false
	}

	var cached channel.CachedChannel
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return channel.Channel{}, false
	}

	return cached.ToChannel(channelId, guildId), true
}

func (c *PgCache) DeleteChannel(channelId uint64) {
	if c.Options.Channels {
		_, _ = c.Exec(context.Background(), `DELETE FROM channels WHERE "channel_id" = $1;`, channelId)
	}
}

func (c *PgCache) DeleteGuildChannels(guildId uint64) {
	if c.Options.Channels {
		_, _ = c.Exec(context.Background(), `DELETE FROM channels WHERE "guild_id" = $1;`, guildId)
	}
}

func (c *PgCache) StoreRole(role guild.Role, guildId uint64) {
	if c.Options.Roles {
		if encoded, err := json.Marshal(role.ToCachedRole(guildId)); err == nil {
			_, _ = c.Exec(context.Background(), `INSERT INTO roles("role_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("role_id", "guild_id") DO UPDATE SET "data" = $3;`, role.Id, guildId, string(encoded))
		}
	}
}

func (c *PgCache) StoreRoles(roles []guild.Role, guildId uint64) {
	if c.Options.Roles {
		batch := &pgx.Batch{}

		for _, role := range roles {
			if encoded, err := json.Marshal(role.ToCachedRole(guildId)); err == nil {
				batch.Queue(`INSERT INTO roles("role_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("role_id", "guild_id") DO UPDATE SET "data" = $3;`, role.Id, guildId, string(encoded))
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()
		_, _ = br.Exec()
	}
}

func (c *PgCache) GetRole(id uint64) (guild.Role, bool) {
	if !c.Options.Roles {
		return guild.Role{}, false
	}

	var raw string
	if err := c.QueryRow(context.Background(), `SELECT "data" FROM roles WHERE "role_id" = $1;`, id).Scan(&raw); err != nil {
		return guild.Role{}, false
	}

	var cached guild.CachedRole
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return guild.Role{}, false
	}

	return cached.ToRole(id), true
}

func (c *PgCache) GetRoles(guildId uint64, ids []uint64) (map[uint64]guild.Role, error) {
	if !c.Options.Roles {
		return nil, nil
	}

	idArray := &pgtype.Int8Array{}
	if err := idArray.Set(ids); err != nil {
		return nil, err
	}

	query := `SELECT "role_id", "data" FROM roles WHERE "role_id" = ANY($1) AND "guild_id" = $2;`
	rows, err := c.Query(context.Background(), query, idArray, guildId)
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

func (c *PgCache) GetRoleInGuild(id, guildId uint64) (guild.Role, bool) {
	if !c.Options.Roles {
		return guild.Role{}, false
	}

	var raw string
	if err := c.QueryRow(context.Background(), `SELECT "data" FROM roles WHERE "role_id" = $1 AND "guild_id" = $2;`, id, guildId).Scan(&raw); err != nil {
		return guild.Role{}, false
	}

	var cached guild.CachedRole
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return guild.Role{}, false
	}

	return cached.ToRole(id), true
}

func (c *PgCache) DeleteRole(roleId uint64) {
	if c.Options.Roles {
		_, _ = c.Exec(context.Background(), `DELETE FROM roles WHERE "role_id" = $1;`, roleId)
	}
}

func (c *PgCache) DeleteGuildRoles(guildId uint64) {
	if c.Options.Channels {
		_, _ = c.Exec(context.Background(), `DELETE FROM roles WHERE "guild_id" = $1;`, guildId)
	}
}

func (c *PgCache) StoreEmoji(emoji emoji.Emoji, guildId uint64) {
	if c.Options.Emojis {
		if encoded, err := json.Marshal(emoji.ToCachedEmoji(guildId)); err == nil {
			_, _ = c.Exec(context.Background(), `INSERT INTO emojis("emoji_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("emoji_id") DO UPDATE SET "data" = $3;`, emoji.Id, guildId, string(encoded))
		}
	}
}

func (c *PgCache) StoreEmojis(emojis []emoji.Emoji, guildId uint64) {
	if c.Options.Emojis {
		batch := &pgx.Batch{}

		for _, e := range emojis {
			if encoded, err := json.Marshal(e.ToCachedEmoji(guildId)); err == nil {
				batch.Queue(`INSERT INTO emojis("emoji_id", "guild_id", "data") VALUES($1, $2, $3) ON CONFLICT("emoji_id") DO UPDATE SET "data" = $3;`, e.Id, guildId, string(encoded))
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()
		_, _ = br.Exec()
	}
}

func (c *PgCache) GetEmoji(id uint64) (emoji.Emoji, bool) {
	if !c.Options.Emojis {
		return emoji.Emoji{}, false
	}

	var raw string
	if err := c.QueryRow(context.Background(), `SELECT "data" FROM emojis WHERE "emoji_id" = $1;`, id).Scan(&raw); err != nil {
		return emoji.Emoji{}, false
	}

	var cached emoji.CachedEmoji
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return emoji.Emoji{}, false
	}

	return cached.ToEmoji(id, user.User{}), true
}

func (c *PgCache) DeleteEmoji(emojiId uint64) {
	if c.Options.Emojis {
		_, _ = c.Exec(context.Background(), `DELETE FROM emojis WHERE "emoji_id" = $1;`, emojiId)
	}
}

func (c *PgCache) StoreVoiceState(state guild.VoiceState) {
	if c.Options.VoiceStates {
		if encoded, err := json.Marshal(state.ToCachedVoiceState()); err == nil {
			_, _ = c.Exec(context.Background(), `INSERT INTO voice_states("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, state.GuildId, state.UserId, string(encoded))
		}
	}
}

func (c *PgCache) StoreVoiceStates(states []guild.VoiceState) {
	if c.Options.VoiceStates {
		batch := &pgx.Batch{}

		for _, state := range states {
			if encoded, err := json.Marshal(state.ToCachedVoiceState()); err == nil {
				batch.Queue(`INSERT INTO voice_states("guild_id", "user_id", "data") VALUES($1, $2, $3) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "data" = $3;`, state.GuildId, state.UserId, string(encoded))
			}
		}

		br := c.SendBatch(context.Background(), batch)
		defer br.Close()
		_, _ = br.Exec()
	}
}

func (c *PgCache) GetVoiceState(userId, guildId uint64) (guild.VoiceState, bool) {
	if !c.Options.VoiceStates {
		return guild.VoiceState{}, false
	}

	query := `
SELECT voice_states.data, members.data, users.data
FROM voice_states
LEFT JOIN members ON members.user_id=voice_states.user_id
LEFT JOIN users ON users.user_id=voice_states.user_id
WHERE voice_states.guild_id = $1 AND voice_states.user_id=$2;
`

	var voiceStateRaw, memberRaw, userRaw string
	if err := c.QueryRow(context.Background(), query, guildId, userId).Scan(&voiceStateRaw, &memberRaw, &userRaw); err != nil {
		return guild.VoiceState{}, false
	}

	var cachedVoiceState guild.CachedVoiceState
	if err := json.Unmarshal([]byte(voiceStateRaw), &cachedVoiceState); err != nil {
		return guild.VoiceState{}, false
	}

	var cachedMember member.CachedMember
	if len(memberRaw) > 0 {
		if err := json.Unmarshal([]byte(memberRaw), &cachedMember); err != nil {
			return guild.VoiceState{}, false
		}
	}

	var cachedUser user.CachedUser
	if len(memberRaw) > 0 {
		if err := json.Unmarshal([]byte(userRaw), &cachedUser); err != nil {
			return guild.VoiceState{}, false
		}
	}

	return cachedVoiceState.ToVoiceState(guildId, cachedMember.ToMember(cachedUser.ToUser(userId))), true
}

func (c *PgCache) GetGuildVoiceStates(guildId uint64) (states []guild.VoiceState) {
	if !c.Options.VoiceStates {
		return
	}

	query := `
SELECT voice_states.user_id, voice_states.data, members.data, users.data
FROM voice_states
LEFT JOIN members ON members.user_id=voice_states.user_id
LEFT JOIN users ON users.user_id=voice_states.user_id
WHERE voice_states.guild_id = $1;
`

	rows, err := c.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var userId uint64
		var voiceStateRaw, memberRaw, userRaw string

		if err := rows.Scan(&userId, &voiceStateRaw, &memberRaw, &userRaw); err != nil {
			continue
		}

		var cachedVoiceState guild.CachedVoiceState
		if err := json.Unmarshal([]byte(voiceStateRaw), &cachedVoiceState); err != nil {
			continue
		}

		var cachedMember member.CachedMember
		if len(memberRaw) > 0 {
			if err := json.Unmarshal([]byte(memberRaw), &cachedMember); err != nil {
				continue
			}
		}

		var cachedUser user.CachedUser
		if len(memberRaw) > 0 {
			if err := json.Unmarshal([]byte(userRaw), &cachedUser); err != nil {
				continue
			}
		}

		states = append(states, cachedVoiceState.ToVoiceState(userId, cachedMember.ToMember(cachedUser.ToUser(userId))))
	}

	return
}

func (c *PgCache) DeleteVoiceState(userId, guildId uint64) {
	if c.Options.Emojis {
		_, _ = c.Exec(context.Background(), `DELETE FROM voice_states WHERE "user_id" = $1 AND "guild_id" = $2;`, userId, guildId)
	}
}

func (c *PgCache) StoreSelf(self user.User) {
	c.selfLock.Lock()
	c.self = self
	c.selfLock.Unlock()
}

func (c *PgCache) GetSelf() (user.User, bool) {
	c.selfLock.RLock()
	self := c.self
	c.selfLock.RUnlock()

	return self, self.Id != 0
}
