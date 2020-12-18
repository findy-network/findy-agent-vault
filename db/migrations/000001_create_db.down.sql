DROP TABLE IF EXISTS "event";

DROP INDEX IF EXISTS "event_cursor_index";

DROP TABLE IF EXISTS "job";

DROP INDEX IF EXISTS "job_cursor_index";

DROP TYPE IF EXISTS "job_result";

DROP TYPE IF EXISTS "job_status";

DROP TYPE IF EXISTS "protocol_type";

DROP INDEX IF EXISTS "message_cursor_index";

DROP TABLE IF EXISTS "message";

DROP TABLE IF EXISTS "proof_attribute";

DROP INDEX IF EXISTS "proof_cursor_index";

DROP TABLE IF EXISTS "proof";

DROP TYPE IF EXISTS "proof_role";

DROP TABLE IF EXISTS "credential_attribute";

DROP INDEX IF EXISTS "credential_cursor_index";

DROP TABLE IF EXISTS "credential";

DROP TYPE IF EXISTS "credential_role";

DROP INDEX IF EXISTS "connection_cursor_index";

DROP TABLE IF EXISTS "connection";

DROP TABLE IF EXISTS "agent";

DROP EXTENSION IF EXISTS "uuid-ossp";