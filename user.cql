use gameapi;
drop table user;
create table user (
  id UUID,
  name text,
  score int,
  gamesplayed int,
  friends set<UUID>,
  version int,
  PRIMARY KEY(id)
);
