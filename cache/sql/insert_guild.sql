INSERT INTO guilds("guild_id", "data")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "data" = $2;