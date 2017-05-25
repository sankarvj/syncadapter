# syncadapter
Adapter helps keep data sync between server and client

# why
Implementation of offline data storage in mobile apps is little complicated task compared to other frameworks.
Syncadapter will solve that problem by exposing interfaces which handles common sync problems. (i.e finding new/updated items
from the list received from server, take the entry from db and convert that entry to server scope or vice versa and etc.)

# goals
Handles common sync logic during get,put,post and delete in a single place.
Find conflicts and handle it gracefully.

# breakup

This library consists of three packages : core, technique, performer.  

* core - base components
* technique - various techniques user can choose to execute at one particular time
* performer - helps you find new/updated/deletd items, helps you convert data from one scope to the other.

# implementation

Existing models in the client system should inherit the basemodel from the syncadapter

                      --- model from syncadapter ---

                        type BaseModel struct {
                            Id      int64    //local id
                            Key     int64    //server id
                            Updated int64    //last updated time
                            Synced  bool     //synced or not
                        }
                          
                          --- your model ---
                       
                       type Ticket struct {
                            Subject   string
                            Desc      string
                            requester string
                            agent     string
                            created   time.Time
                            core.BaseModel //Embed BaseModel
                        }

So that all the methods declared under BaseModel is promoted to be accessed via other models inherit it

                        ticket.PrepareLocal()
                        ticket.SetLocalId(id int64)
                        etc...

Invidual models can override promoted methods in case if any modification needed

                        func (obj *BaseModel) SetLocalId(id int64) {
                            obj.Id = id
                        }

                        func (obj *Ticket) SetLocalId(id int64) {
                            obj.Id = id * 10
                        }

If a column in a table references a id from the other table than that column must be tagged like the below structure 

                        type Note struct {
                            Ticketid int64 `rt:"tickets" rk:"id"`
                            Name     string
                            Desc     string
                            created  time.Time
                            core.BaseModel
                        }

    Note : ticketid column of notes table references id column of Ticket table. 

Once your model extends the basemodel. All your models are now free to call performer methods to make use of the sync
helpers.

For in depth implementation details refer this [example project](https://github.com/sankarvj/sample_syncadapter_client)


# edge cases
 * What if the reference of one table points to other table and local db don't have that table ? (Say ticket table has assignee id as a column and in the server side assignee id is the userid of users table)
 
 Solution 1 : Local table needs to be modified in accordance with the data it needs to show.

 If an API updates two model in the server. (Create Ticket API create a ticket and adds the current user to assignee table implicitly)
   Is that format is good REST API format ? Or is it not ? 

 Solution 2 : Your implementation has to handle the response differently

 * What if it really updates the server end and fails to update the changes in the localdb ?
 
 Solution 1 : Send a client id and store it in the server db
 Solution 2 : Take action based on server error code

 * What if user changes his time in the local device ?

    Solution 1) Server will update the local time on each response.

 * Conflict of two users needs to be shown in UI. How ? What is the purpose ?

 * Tags vs Interface (Tags is too brutal. Use Interface in each class to convert server to db item. But I think tags are best)

 * Periodic needs to include GET first and then only it should post,put.




