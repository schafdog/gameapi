use gameapi;
create table user (
  id UUID,
  name text,
  highscore int,
  friends set<UUID>,
  version int,
  PRIMARY KEY(id)
);
