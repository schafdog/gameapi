curl -s -X POST -H "Content-Type: application/json" -d '{ "missing_name": "John" }'    "http://localhost:8000/user" > fail_user.log
curl -s -X POST -H "Content-Type: application/json" -d '{ "name": "" }'    "http://localhost:8000/user" > fail_user.log
curl -s -X POST -H "Content-Type: application/json" -d '{ "name": "John" }'    "http://localhost:8000/user" > new_user.log
USER=`cat new_user.log | grep "id" | cut -d '"' -f 4`
echo USER $USER
curl -s -X GET -H "Content-Type: application/json"                            "http://localhost:8000/user" > get_users.log
curl -s -X GET -H "Content-Type: application/json"                            "http://localhost:8000/user/$USER/stats" > get_new_user_stats.log

curl -s -X GET -H "Content-Type: application/json"                                       "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/stats" > get_first_stats.log
curl -s -X PUT -H "Content-Type: application/json" -d '{ "gamesPlayed": 0, score: 0}'    "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/stats" > clear_stats.log
curl -s -X PUT -H "Content-Type: application/json"                                       "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/stats" > get_clear_stats.log
curl -s -X PUT -H "Content-Type: application/json" -d '{ "gamesPlayed": 10, score: 100}' "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/stats" > put_stats.log
curl -s -X PUT -H "Content-Type: application/json" -d '{ "Friends": [] }'                "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/friends" > get_stats.log
curl -s -X GET -H "Content-Type: application/json"                                       "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/friends" > get_empty_friends.log
curl -s -X PUT -H "Content-Type: application/json" -d '{ "Friends": ["171a1963-7f9f-11e8-90df-6a0001660200", "edd52600-7f9e-11e8-90da-6a0001660200"] }' \
                                                                                      "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/friends" > put_friends.log
curl -s -X GET -H "Content-Type: application/json"                                       "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/friends" > get_friends.log
curl -s -X PUT -H "Content-Type: application/json" -d '{ "Friends": [] }' \
                                                                                      "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/friends" > reset_friends.log
curl -s -X GET -H "Content-Type: application/json"                                       "http://localhost:8000/user/c21039b3-7f8b-11e8-8a31-6a0001660200/friends" > get_empty_friends.log
