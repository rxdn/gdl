INSERT INTO roles("role_id", "guild_id", "data")
VALUES($1, $2, $3)
ON CONFLICT("role_id", "guild_id")
DO UPDATE SET "data" = $3;