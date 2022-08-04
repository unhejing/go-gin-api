#### go_gin_api.sys_user 
用户表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 主键 | int(11) unsigned | PRI | NO | auto_increment |  |
| 2 | username | 用户名 | varchar(32) | UNI | NO |  |  |
| 3 | password | 密码 | varchar(100) |  | NO |  |  |
| 4 | nickname | 昵称 | varchar(60) |  | NO |  |  |
| 5 | mobile | 手机号 | varchar(20) |  | NO |  |  |
| 6 | is_used | 是否启用 1:是  -1:否 | tinyint(1) |  | NO |  | 1 |
| 7 | is_deleted | 是否删除 1:是  -1:否 | tinyint(1) |  | NO |  | -1 |
| 8 | created_time | 创建时间 | datetime |  | NO |  | CURRENT_TIMESTAMP |
| 9 | updated_time | 更新时间 | datetime |  | NO | on update CURRENT_TIMESTAMP | CURRENT_TIMESTAMP |
