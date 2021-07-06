create table url
(
    id         serial       not null
        constraint url_pk primary key,
    name       varchar(255),
    link       varchar(510) not null,
    hash       varchar(32)  not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);

alter table url
    owner to url_shortener;

create unique index url_hash_uindex on url (hash);

create table stats
(
    url_id integer           not null
        constraint stats_pk
            primary key
        constraint stats_url_id_fk
            references url
            on update cascade on delete cascade,
    hits   integer default 0 not null
);

alter table stats
    owner to url_shortener;

CREATE FUNCTION url_after_inset() RETURNS trigger AS $$
BEGIN
    INSERT INTO stats (url_id, hits) VALUES (NEW.id, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER url_after_inset AFTER INSERT ON url
    FOR EACH ROW EXECUTE PROCEDURE url_after_inset();

