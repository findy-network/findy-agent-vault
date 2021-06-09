UPDATE agent SET raw_jwt = '' WHERE raw_jwt IS NULL;
ALTER TABLE agent ALTER COLUMN raw_jwt SET NOT NULL;
UPDATE agent SET label = '' WHERE label IS NULL;
ALTER TABLE agent ALTER COLUMN label SET NOT NULL;


UPDATE connection SET approved = timestamp '0001-01-01' WHERE approved IS NULL;
ALTER TABLE connection ALTER COLUMN approved SET NOT NULL;
ALTER TABLE connection ALTER COLUMN approved SET DEFAULT timestamp '0001-01-01';

UPDATE connection SET archived = timestamp '0001-01-01' WHERE archived IS NULL;
ALTER TABLE connection ALTER COLUMN archived SET NOT NULL;
ALTER TABLE connection ALTER COLUMN archived SET DEFAULT timestamp '0001-01-01';


UPDATE credential SET approved = timestamp '0001-01-01' WHERE approved IS NULL;
ALTER TABLE credential ALTER COLUMN approved SET NOT NULL;
ALTER TABLE credential ALTER COLUMN approved SET DEFAULT timestamp '0001-01-01';

UPDATE credential SET issued = timestamp '0001-01-01' WHERE issued IS NULL;
ALTER TABLE credential ALTER COLUMN issued SET NOT NULL;
ALTER TABLE credential ALTER COLUMN issued SET DEFAULT timestamp '0001-01-01';

UPDATE credential SET failed = timestamp '0001-01-01' WHERE failed IS NULL;
ALTER TABLE credential ALTER COLUMN failed SET NOT NULL;
ALTER TABLE credential ALTER COLUMN failed SET DEFAULT timestamp '0001-01-01';

UPDATE credential SET archived = timestamp '0001-01-01' WHERE archived IS NULL;
ALTER TABLE credential ALTER COLUMN archived SET NOT NULL;
ALTER TABLE credential ALTER COLUMN archived SET DEFAULT timestamp '0001-01-01';


UPDATE proof SET provable = timestamp '0001-01-01' WHERE provable IS NULL;
ALTER TABLE proof ALTER COLUMN provable SET NOT NULL;
ALTER TABLE proof ALTER COLUMN provable SET DEFAULT timestamp '0001-01-01';

UPDATE proof SET approved = timestamp '0001-01-01' WHERE approved IS NULL;
ALTER TABLE proof ALTER COLUMN approved SET NOT NULL;
ALTER TABLE proof ALTER COLUMN approved SET DEFAULT timestamp '0001-01-01';

UPDATE proof SET verified = timestamp '0001-01-01' WHERE verified IS NULL;
ALTER TABLE proof ALTER COLUMN verified SET NOT NULL;
ALTER TABLE proof ALTER COLUMN verified SET DEFAULT timestamp '0001-01-01';

UPDATE proof SET failed = timestamp '0001-01-01' WHERE failed IS NULL;
ALTER TABLE proof ALTER COLUMN failed SET NOT NULL;
ALTER TABLE proof ALTER COLUMN failed SET DEFAULT timestamp '0001-01-01';

UPDATE proof SET archived = timestamp '0001-01-01' WHERE archived IS NULL;
ALTER TABLE proof ALTER COLUMN archived SET NOT NULL;
ALTER TABLE proof ALTER COLUMN archived SET DEFAULT timestamp '0001-01-01';


UPDATE message SET archived = timestamp '0001-01-01' WHERE archived IS NULL;
ALTER TABLE message ALTER COLUMN archived SET NOT NULL;
ALTER TABLE message ALTER COLUMN archived SET DEFAULT timestamp '0001-01-01';
