INSERT INTO users("user_id", "data", last_seen)
VALUES($1, $2, NOW())
ON CONFLICT("user_id")
DO UPDATE SET "data" = excluded.data, "last_seen" = excluded.last_seen;