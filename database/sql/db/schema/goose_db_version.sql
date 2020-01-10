CREATE TABLE goose_db_version
(
    id integer NOT NULL DEFAULT nextval('goose_db_version_id_seq'::regclass),
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now(),
    CONSTRAINT goose_db_version_pkey PRIMARY KEY (id)
);
