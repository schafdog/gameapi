use gameapi;
create table user (
  id UUID,
  name text,
  highscore int,
  friends text,
  version int,
  PRIMARY KEY(id)
);
