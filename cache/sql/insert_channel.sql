INSERT INTO channels("channel_id", "guild_id", "data")
VALUES($1, $2, $3)
ON CONFLICT("channel_id")
DO UPDATE SET "data" = $3;