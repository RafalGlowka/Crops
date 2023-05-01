ALTER TABLE users RENAME COLUMN seller TO verified;
UPDATE users SET verified = True;
ALTER TABLE users ADD verificationcode text NOT NULL DEFAULT md5(random()::text);
ALTER TABLE users ADD CONSTRAINT uniqueVerificationCode UNIQUE (verificationcode);
ALTER TABLE users ADD balance BIGINT NOT NULL;
ALTER TABLE fields ADD verified BOOL NOT NULL DEFAULT False;
