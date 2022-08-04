#### go_gin_api.sys_config 
配置表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | id | int(11) | PRI | NO | auto_increment |  |
| 2 | note | 备注 | varchar(255) |  | YES |  |  |
| 3 | params_name | 参数名称 | varchar(255) |  | YES |  |  |
| 4 | params_key | 参数键名 | varchar(255) |  | YES |  |  |
| 5 | params_value | 参数键值 | varchar(255) |  | YES |  |  |
| 6 | channel_tag | 渠道标签 | varchar(255) |  | YES |  |  |
| 7 | created_time | 创建时间 | datetime |  | YES |  | CURRENT_TIMESTAMP |
| 8 | updated_time | 更新时间 | datetime |  | YES | on update CURRENT_TIMESTAMP | CURRENT_TIMESTAMP |
