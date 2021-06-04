UPDATE agent SET raw_jwt = '' WHERE raw_jwt IS NULL;
ALTER TABLE agent ALTER COLUMN raw_jwt SET NOT NULL;
UPDATE agent SET label = '' WHERE label IS NULL;
ALTER TABLE agent ALTER COLUMN label SET NOT NULL;
