CREATE
VIEW `view_user_login`AS
SELECT a.user_id AS user_id, b.record_time AS user_login_time, a.cmbi_channel_id AS
user_game_channel FROM realm_account.t_user a, realm_log.t_login b WHERE a.user_id = b.user_id ORDER BY a.user_id DESC