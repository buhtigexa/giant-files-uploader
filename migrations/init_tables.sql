create table `data`
(
    int      bigint       not null,
    filename varchar(255) not null,
    time     bigint

)


CREATE TABLE data (
                      id SERIAL PRIMARY KEY,
                      filename TEXT NOT NULL,
                      total REAL NOT NULL DEFAULT 0,
                      time BIGINT NOT NULL
);



