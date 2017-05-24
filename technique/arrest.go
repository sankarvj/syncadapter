package technique

//#Core goals
//* minimize server load
//* really useful if periodic technique implemented using network changes and if network changes often occurs due to other conditions

//#Needs server side implementation
//* yes

//#Logic
//* server should send a incremental_random_number,timestamp,table over impulse/periodic/basic technique.
//* if (server_inc_id - local_inc_id) != 0 this is called as event of conflict in this case arrester will let API query go hit the server.
//* otherwise it stops
