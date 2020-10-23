CREATE
VIEW `view_user_pay` AS
SELECT b.user_id as user_id, b.server_id AS user_game_server, a.role_id AS role_id,
 a.role_name AS role_name, b.pay_type AS role_pay_type, b.order_id AS role_order_id, b.product_id AS role_product_id, b.coin AS role_pay_coin,
b.status AS role_pay_status, b.purchase_time AS role_purchase_time FROM realm_adapter.t_role a, realm_account.t_purchase_record b WHERE a.role_id = b.role_id  and b.pay_type =2 ORDER BY a.user_id ASC
