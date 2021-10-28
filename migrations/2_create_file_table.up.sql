CREATE TABLE file_meta (
    file_name VARCHAR(256) NOT NULL,
    uid VARCHAR(36),
    last_modified BIGINT NOT NULL,
    md5_hash TEXT NOT NULL,
    file_contents BIGINT NOT NULL,
    version VARCHAR(28) NOT NULL,
    CONSTRAINT FK_uid_uid FOREIGN KEY (uid) REFERENCES auth (uid),
    CONSTRAINT pk_filename_version PRIMARY KEY (file_name, version, uid)
);