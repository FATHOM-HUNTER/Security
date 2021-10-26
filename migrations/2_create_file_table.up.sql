CREATE TABLE file_meta (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    file_name TEXT NOT NULL,
    uid VARCHAR(36),
    last_modified BIGINT NOT NULL,
    md5_hash TEXT NOT NULL,
    file_contents BIGINT NOT NULL,
    version INT NOT NULL,
    CONSTRAINT FK_uid_uid FOREIGN KEY (uid) REFERENCES auth (uid)
);