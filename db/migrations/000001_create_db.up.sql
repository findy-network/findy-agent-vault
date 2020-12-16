CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "agent" (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  agent_id VARCHAR(256) UNIQUE NOT NULL,
  label VARCHAR(1024),
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  last_accessed timestamptz NOT NULL DEFAULT (now() at time zone 'UTC')
);

CREATE INDEX "agent_id_index" ON agent (agent_id);

CREATE INDEX "agent_cursor_index" ON agent (cursor);

CREATE TABLE "connection"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tenant_id uuid NOT NULL,
  our_did VARCHAR(256) NOT NULL,
  their_did VARCHAR(256) NOT NULL,
  their_endpoint VARCHAR(4096) NOT NULL,
  their_label VARCHAR(1024),
  invited BOOLEAN NOT NULL DEFAULT FALSE,
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  approved timestamptz,
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  CONSTRAINT fk_connection_agent
    FOREIGN KEY(tenant_id) REFERENCES agent(id)
);

CREATE INDEX "connection_cursor_index" ON connection (cursor);

CREATE TYPE "credential_role" AS ENUM ('ISSUER', 'HOLDER');

CREATE TABLE "credential"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tenant_id uuid NOT NULL,
  connection_id uuid NOT NULL,
  role credential_role NOT NULL,
  schema_id VARCHAR(4096) NOT NULL,
  cred_def_id VARCHAR(4096) NOT NULL,
  initiated_by_us BOOLEAN NOT NULL DEFAULT FALSE,
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  approved timestamptz,
  issued timestamptz,
  failed timestamptz,
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  CONSTRAINT fk_credential_agent
    FOREIGN KEY(tenant_id) REFERENCES agent(id),
  CONSTRAINT fk_credential_connection
    FOREIGN KEY(connection_id) REFERENCES connection(id)
);

CREATE INDEX "credential_cursor_index" ON credential (cursor);

CREATE TABLE "credential_attribute"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  credential_id uuid NOT NULL,
  "name" VARCHAR(1024) NOT NULL,
  "value" VARCHAR(4096) NOT NULL,
  index SMALLINT NOT NULL,
  CONSTRAINT fk_credential_attribute_credential
    FOREIGN KEY(credential_id) REFERENCES "credential"(id)
);

CREATE TYPE "proof_role" AS ENUM ('VERIFIER', 'PROVER');

CREATE TABLE "proof"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tenant_id uuid NOT NULL,
  connection_id uuid NOT NULL,
  role proof_role NOT NULL,
  initiated_by_us BOOLEAN NOT NULL DEFAULT FALSE,
  result BOOLEAN NOT NULL DEFAULT FALSE,
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  approved timestamptz,
  verified timestamptz,
  failed timestamptz,
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  CONSTRAINT fk_proof_agent
    FOREIGN KEY(tenant_id) REFERENCES agent(id),
  CONSTRAINT fk_proof_connection
    FOREIGN KEY(connection_id) REFERENCES connection(id)
);

CREATE INDEX "proof_cursor_index" ON proof (cursor);

CREATE TABLE "proof_attribute"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  proof_id uuid NOT NULL,
  "name" VARCHAR(1024) NOT NULL,
  "value" VARCHAR(4096),
  "cred_def_id" VARCHAR(4096) NOT NULL,
  index SMALLINT NOT NULL,
  CONSTRAINT fk_proof_attribute_proof
    FOREIGN KEY(proof_id) REFERENCES "proof"(id)
);

CREATE TABLE "message"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tenant_id uuid NOT NULL,
  connection_id uuid NOT NULL,
  message VARCHAR(4096) NOT NULL,
  sent_by_me BOOLEAN NOT NULL DEFAULT FALSE,
  delivered BOOLEAN DEFAULT NULL,
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  CONSTRAINT fk_message_agent
    FOREIGN KEY(tenant_id) REFERENCES agent(id),
  CONSTRAINT fk_message_connection
    FOREIGN KEY(connection_id) REFERENCES connection(id)
);

CREATE INDEX "message_cursor_index" ON message (cursor);

CREATE TYPE "protocol_type" AS ENUM ('NONE', 'CONNECTION', 'CREDENTIAL', 'PROOF', 'BASIC_MESSAGE');

CREATE TYPE "job_status" AS ENUM ('WAITING', 'PENDING', 'COMPLETE');

CREATE TYPE "job_result" AS ENUM ('NONE', 'SUCCESS', 'FAILURE');

CREATE TABLE "job"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tenant_id uuid NOT NULL,
  connection_id uuid,
  protocol_id uuid,
  protocol_type protocol_type NOT NULL DEFAULT 'NONE',
  "status" job_status NOT NULL DEFAULT 'WAITING',
  result job_result NOT NULL DEFAULT 'NONE',
  initiated_by_us BOOLEAN NOT NULL DEFAULT FALSE,
  updated timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  CONSTRAINT fk_job_agent
    FOREIGN KEY(tenant_id) REFERENCES agent(id),
  CONSTRAINT fk_job_connection
    FOREIGN KEY(connection_id) REFERENCES connection(id)
);

CREATE INDEX "job_cursor_index" ON job (cursor);

CREATE TABLE "event"(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tenant_id uuid NOT NULL,
  connection_id uuid,
  job_id uuid,
  "description" VARCHAR(4096) NOT NULL,
  "read" BOOLEAN NOT NULL DEFAULT FALSE,
  created timestamptz NOT NULL DEFAULT (now() at time zone 'UTC'),
  cursor BIGINT NOT NULL GENERATED ALWAYS AS (extract(epoch from created at time zone 'UTC') * 1000) STORED,
  CONSTRAINT fk_event_agent
    FOREIGN KEY(tenant_id) REFERENCES agent(id),
  CONSTRAINT fk_event_connection
    FOREIGN KEY(connection_id) REFERENCES connection(id),
  CONSTRAINT fk_event_job
    FOREIGN KEY(job_id) REFERENCES job(id)
);

CREATE INDEX "event_cursor_index" ON event (cursor);
