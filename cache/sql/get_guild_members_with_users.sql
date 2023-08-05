SELECT members.user_id, members.data, users.data
FROM members
LEFT JOIN users ON members.user_id = users.user_id
WHERE "guild_id" = $1;