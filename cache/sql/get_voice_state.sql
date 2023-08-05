SELECT voice_states.data, members.data, users.data
FROM voice_states
LEFT JOIN members ON members.user_id=voice_states.user_id
LEFT JOIN users ON users.user_id=voice_states.user_id
WHERE voice_states.guild_id = $1 AND voice_states.user_id=$2;