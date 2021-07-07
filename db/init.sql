CREATE DATABASE url_shortener
    WITH
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

\c url_shortener;

CREATE USER url_shortener WITH password 'url_shortener';

GRANT ALL ON DATABASE url_shortener TO url_shortener;
