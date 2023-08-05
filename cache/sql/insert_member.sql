INSERT INTO members("guild_id", "user_id", "data", "last_seen")
VALUES($1, $2, $3, NOW())
ON CONFLICT("guild_id", "user_id")
DO UPDATE SET "data" = excluded.data, "last_seen" = excluded.last_seen;