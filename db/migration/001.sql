create table url
(
    id         serial       not null,
    name       varchar(255),
    link       varchar(510) not null,
    hash       varchar(32)  not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);

alter table url owner to url_shortener;

create unique index url_hash_uindex on url (hash);