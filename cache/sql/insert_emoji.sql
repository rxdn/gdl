INSERT INTO emojis("emoji_id", "guild_id", "data")
VALUES($1, $2, $3)
ON CONFLICT("emoji_id")
DO UPDATE SET "data" = $3;