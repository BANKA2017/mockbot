# 🤖Mock Bot \<DEV\>
---

## Flags

| flag       | default                    | description              |
| :--------- | :------------------------- | :----------------------- |
| `path`     |                            | Path for assets/database |
| `webhook`  | `:1323`                    | webhook address          |
| `endpoint` | `https://api.telegram.org` | Bot api endpoint         |
| `test`     | `false`                    | Test mode                |

## Bot commands

### Standard

| command          | format                                                    | description |
| :--------------- | :-------------------------------------------------------- | :---------- |
| `/hey`           |                                                           | TODO        |
| `/me`            |                                                           | TODO        |
| `/get`           | `/get chat_setting_name`                                  | TODO        |
| `/set`           | `/set chat_setting_name value [value2 [value3...]]`       | TODO        |
| `/chat_settings` |                                                           | TODO        |
| `/rank`          |                                                           | TODO        |
| `/system_set`    | `/system_set bot_setting_name value [value2 [value3...]]` | TODO        |
| `/system_get`    | `/system_get bot_setting_name`                            | TODO        |
| `/bot_settings`  |                                                           | TODO        |

### AT

- reply:/aaa [bbb ccc ddd...]
  - "甲 aaa 乙"
  - "甲 aaa了 乙 [bbb ccc ddd...]"

### 🐱meow~

- AT the bot or reply to the bot with prefix or suffix **喵一个** or words **喵**
  - "喵"

## User

- staff -> bot
- administrator -> chat

## Chat type

- private, group, supergroup

Not yet supported `channel` at all 

## Init

- Create a SQLite file named `mockbot.db` by `mockbot.db.sql`
- Download a font file and rename it to `font.ttf`
- `insert into bot_settings (bot_id,key,value) values ('<bot_id>','token','<bot_token>');`
- `insert into staff (bot_id,role,user_id) values ('<bot_id>','owner','<user_id>');`
- `go run main.go --path=.`
- ...

## TODO

- [ ] Search `// TODO` or `//TODO` form source code
- [ ] i18n, Chinese simplify only now
- [ ] Modify staff by command
- [ ] Add bot by command
- [ ] Web API
- [ ] ...