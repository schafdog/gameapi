OPTS="-s -H \"Content-Type: application/json\" "

curl $OPTS -X POST -d '{ "missing_name": "John" }'    "http://localhost:8000/user"      > 1_fail_user.log
curl $OPTS -X POST -d '{ "name": "" }'    "http://localhost:8000/user"                  > 2_fail_user.log
curl $OPTS -X POST -H "X-UUID: deaddead-dead-dead-dead-deaddeaddead" -d '{ "name": "Dennis" }'    \
                                                          "http://localhost:8000/user"  > 3_new_user.log
USERID=`cat 3_new_user.log | grep "id" | cut -d '"' -f 4`

curl $OPTS -X POST -H "X-UUID: 18dd75e9-3d4a-48e2-bafc-3c8f95a8f0d1" -d '{ "name": "John" }'    \
     "http://localhost:8000/user"  > 4_john.log
FRIEND1=`cat 4_john.log | grep "id" | cut -d '"' -f 4`
curl $OPTS -X POST -H "X-UUID: f9a9af78-6681-4d7d-8ae7-fc41e7a24d08" -d '{ "name": "Bob" }'    \
                                                          "http://localhost:8000/user"  > 4_bob.log
FRIEND2=`cat 4_bob.log | grep "id" | cut -d '"' -f 4`
curl $OPTS -X POST -H "X-UUID: 2d18862b-b9c3-40f5-803e-5e100a520249" -d '{ "name": "Alice" }'    \
                                                          "http://localhost:8000/user"  > 4_alice.log
FRIEND3=`cat 4_alice.log | grep "id" | cut -d '"' -f 4`
echo "Friends $FRIEND1, $FRIEND2, $FRIEND3" 
curl $OPTS -X PUT -d '{ "gamesPlayed": 12, "score": 322}'   "http://localhost:8000/user/$FRIEND1/state"  > 10_put_state_f1.log
curl $OPTS -X PUT -d '{ "gamesPlayed": 9,  "score": 21}'    "http://localhost:8000/user/$FRIEND2/state"  > 10_put_state_f2.log
curl $OPTS -X PUT -d '{ "gamesPlayed": 99, "score": 99332}' "http://localhost:8000/user/$FRIEND3/state"  > 10_put_state_f3.log

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
curl $OPTS -X PUT -d "{ \"friends\": [\"$FRIEND1\", \"$FRIEND2\", \"$FRIEND3\"] }" \
                                                          "http://localhost:8000/user/$USERID/friends" > 14_put_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/friends" > 15_get_friends.log
curl $OPTS -X PUT -d '{ "friends": [] }'                  "http://localhost:8000/user/$USERID/friends" > 16_reset_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID/friends" > 17_get_empty_friends.log
curl $OPTS -X GET                                         "http://localhost:8000/user/$USERID"         > 18_get_new_user_again.log
#curl $OPTS -X DELETE                                      "http://localhost:8000/user/$USERID"         > 19_delete_user_again.log
#curl $OPTS -X DELETE                                      "http://localhost:8000/user/$USERID"         > 19_delete_Friend_1.log
#curl $OPTS -X DELETE                                      "http://localhost:8000/user/$USERID"         > 19_delete_Friend_2.log
