OPTS="-s -H \"Content-Type: application/json\" "

curl $OPTS -X POST -d '{ "missing_name": "John" }'    "http://localhost:8000/user"      > 1_fail_user.log
curl $OPTS -X POST -d '{ "name": "" }'    "http://localhost:8000/user"                  > 2_fail_user.log
curl $OPTS -X POST -H "X-UUID: deaddead-dead-dead-dead-deaddeaddead" -d '{ "name": "Dennis" }'    \
                                                          "http://localhost:8000/user"  > 3_new_user.log
USERID=`cat 3_new_user.log | grep "id" | cut -d '"' -f 4`
curl $OPTS -X POST -H "X-UUID: 171a1963-7f9f-11e8-90df-6a0001660200" -d '{ "name": "Doe" }'    \
                                                          "http://localhost:8000/user"  > 4_new_friend.log
FRIEND1=`cat 4_new_friend.log | grep "id" | cut -d '"' -f 4`
curl $OPTS -X POST -H "X-UUID: edd52600-7f9e-11e8-90da-6a0001660200" -d '{ "name": "John" }'    \
                                                          "http://localhost:8000/user"  > 5_new_friend.log
FRIEND2=`cat 5_new_friend.log | grep "id" | cut -d '"' -f 4`
curl $OPTS -X GET                                         "http://localhost:8000/user"  > 6_get_users.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID"         > 5_get_new_user.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/state"   > 6_get_new_user_state.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/state"   > 7_get_first_state.log
curl $OPTS -X PUT -d '{ "gamesPlayed": 0, "score": 0 }'   "http://localhost:8000/user/$USERID/state"   > 8_clear_state.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/state"   > 9_get_clear_state.log
curl $OPTS -X PUT -d '{ "gamesPlayed": 10, "score": 100}' "http://localhost:8000/user/$USERID/state"   > 10_put_state.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/state"   > 11_get_state.log
curl $OPTS -X PUT -d '{ "friends": [] }'                  "http://localhost:8000/user/$USERID/friends" > 12_clear_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/friends" > 13_get_empty_friends.log
curl $OPTS -X PUT -d '{ "friends": ["171a1963-7f9f-11e8-90df-6a0001660200", "edd52600-7f9e-11e8-90da-6a0001660200"] }' \
                                                          "http://localhost:8000/user/$USERID/friends" > 14_put_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/friends" > 15_get_friends.log
curl $OPTS -X PUT -d '{ "friends": [] }'                  "http://localhost:8000/user/$USERID/friends" > 16_reset_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/friends" > 17_get_empty_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID"         > 18_get_new_user_again.log
curl $OPTS -X DELETE                                      "http://localhost:8000/user/$USERID"         > 19_delete_user_again.log
curl $OPTS -X DELETE                                      "http://localhost:8000/user/$USERID"         > 19_delete_Friend_1.log
curl $OPTS -X DELETE                                      "http://localhost:8000/user/$USERID"         > 19_delete_Friend_2.log
