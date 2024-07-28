BEGIN TRANSACTION;
DROP TABLE IF EXISTS "settings";
CREATE TABLE IF NOT EXISTS "settings" (
	"key"	TEXT,
	"value"	TEXT,
	PRIMARY KEY("key")
);
DROP TABLE IF EXISTS "checkin";
CREATE TABLE IF NOT EXISTS "checkin" (
	"id"	INTEGER UNIQUE,
	"user_id"	TEXT,
	"date"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "bot_settings";
CREATE TABLE IF NOT EXISTS "bot_settings" (
	"bot_id"	TEXT,
	"key"	TEXT,
	"value"	TEXT,
	PRIMARY KEY("bot_id","key")
) WITHOUT ROWID;
DROP TABLE IF EXISTS "chat_settings";
CREATE TABLE IF NOT EXISTS "chat_settings" (
	"chat_id"	TEXT,
	"key"	TEXT,
	"value"	TEXT,
	PRIMARY KEY("chat_id","key")
) WITHOUT ROWID;
DROP TABLE IF EXISTS "message";
CREATE TABLE IF NOT EXISTS "message" (
	"message_id"	TEXT UNIQUE,
	"bot_id"	TEXT,
	"chat_id"	TEXT,
	"date"	INTEGER,
	"content"	TEXT,
	"raw_content"	TEXT,
	"auto_delete"	INTEGER,
	PRIMARY KEY("message_id","chat_id")
);
DROP TABLE IF EXISTS "staff";
CREATE TABLE IF NOT EXISTS "staff" (
	"user_id"	INTEGER NOT NULL,
	"role"	TEXT NOT NULL,
	"bot_id"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("user_id","bot_id")
);
DROP TABLE IF EXISTS "group_message";
CREATE TABLE IF NOT EXISTS "group_message" (
	"message_id"	TEXT UNIQUE,
	"chat_id"	TEXT,
	"user_id"	TEXT,
	"full_name"	TEXT,
	"date"	INTEGER,
	"text"	TEXT,
	"raw_content"	TEXT,
	PRIMARY KEY("message_id","chat_id")
);
DROP INDEX IF EXISTS "idx_bot_settings_primary";
CREATE UNIQUE INDEX IF NOT EXISTS "idx_bot_settings_primary" ON "bot_settings" (
	"bot_id",
	"key"
);
DROP INDEX IF EXISTS "idx_chat_settings_primary";
CREATE UNIQUE INDEX IF NOT EXISTS "idx_chat_settings_primary" ON "chat_settings" (
	"chat_id",
	"key"
);
DROP INDEX IF EXISTS "idx_checkin_user_id_date";
CREATE INDEX IF NOT EXISTS "idx_checkin_user_id_date" ON "checkin" (
	"user_id",
	"date"
);
DROP INDEX IF EXISTS "idx_message_message_id_chat_id";
CREATE INDEX IF NOT EXISTS "idx_message_message_id_chat_id" ON "message" (
	"message_id",
	"chat_id"
);
DROP INDEX IF EXISTS "idx_message_bot_id_auto_delete_date";
CREATE INDEX IF NOT EXISTS "idx_message_bot_id_auto_delete_date" ON "message" (
	"bot_id",
	"auto_delete",
	"date"
);
DROP INDEX IF EXISTS "idx_checkin_date_user_id";
CREATE INDEX IF NOT EXISTS "idx_checkin_date_user_id" ON "checkin" (
	"date"	DESC,
	"user_id"
);
DROP INDEX IF EXISTS "idx_group_message_chat_id_date";
CREATE INDEX IF NOT EXISTS "idx_group_message_chat_id_date" ON "group_message" (
	"chat_id",
	"date"
);
COMMIT;
