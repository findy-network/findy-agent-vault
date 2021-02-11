ALTER TABLE "connection" ADD COLUMN archived timestamptz;
ALTER TABLE "message" ADD COLUMN archived timestamptz;
ALTER TABLE "proof" ADD COLUMN archived timestamptz;
ALTER TABLE "credential" ADD COLUMN archived timestamptz;