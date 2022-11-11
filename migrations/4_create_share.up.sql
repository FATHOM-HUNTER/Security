CREATE TABLE shares (
    uid VARCHAR(36) NOT NULL PRIMARY KEY,
    file_name VARCHAR(36) NOT NULL,
    owner_uid VARCHAR(36) NOT NULL,
    shared_with_email text NOT NULL,
    CONSTRAINT FK_file_name FOREIGN KEY (file_name) REFERENCES file_meta (file_name),
    CONSTRAINT FK_owner_uid FOREIGN KEY (owner_uid) REFERENCES file_meta (uid)
);