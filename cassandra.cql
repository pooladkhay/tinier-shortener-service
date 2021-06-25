CREATE TABLE urls (
   url TEXT,
   hash TEXT,
   user TEXT,
   is_private BOOLEAN,
   secret TEXT,
   created_at TIMESTAMP,
   expires_at TIMESTAMP,
   PRIMARY KEY (hash, user)
);
