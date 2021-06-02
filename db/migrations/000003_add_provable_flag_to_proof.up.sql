ALTER TABLE "proof" ADD COLUMN provable timestamptz;
ALTER TYPE job_status ADD VALUE 'BLOCKED'; 