CREATE
VIEW `view_user_create`AS
SELECT c.user_id AS user_id, c.create_time AS user_create_time, c.cmbi_channel_id AS user_game_channel,
a.role_id AS role_id, a.role_name AS role_name, b.fCreateTime AS role_create_time FROM realm_adapter.t_role a,
realm_role.tchar_new b, realm_account.t_user c WHERE a.role_id = b.fCharID AND a.user_id = c.user_id
ORDER BY a.user_id DESC