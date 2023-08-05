INSERT INTO voice_states("guild_id", "user_id", "data")
VALUES($1, $2, $3)
ON CONFLICT("guild_id", "user_id")
DO UPDATE SET "data" = $3;