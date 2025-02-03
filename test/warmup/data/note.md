## Note

If using [bbolt](https://github.com/etcd-io/bbolt), every day notes can be a bucket. The name can be something like `20250101`. 

There's built-in [backup](https://github.com/etcd-io/bbolt?tab=readme-ov-file#database-backups) mechanism and it should in some case but it's a single file, would be easy too.

### Architect trail #1

For example: Bucket `20250101` will have notes like this.

| Log_Timestamp | Event_Tags | Event_DateTime | Event_Content |
|---------------|------------|----------------|---------------|
|20250101_230000| Workout,Note | 20241231, this is for yesterday| I did some dumb stuff |
|20250103_002900| BJJ, Cardio | 20250102, this is for yesterday| Went to BJJ class and also did some cardio |
|...|...|...|...|
|Key column | NOT nullable | NOT nullable | nullable|

When rendering (pick the record) to display, what this the event key? *Thinking to have use `event tag and date time`*, because you would only have one event type on a day. E.g., you did reading yesterday, even you read multiple times, it simply just checks that. 


### Architect trail #2

To build something on an `event` based structure. It might makes sense to do things like below

- Event will be unique everyday, i.e. you won't have *two* reading entry on a same day
- You can have different event every day, i.e., reading and running
- Every event can have content, but can be optional
- When inputs data, it's possible to have something in bulk, meaning to have *running, lifting, reading* all loged at the same time.
    - With that, the unique keys in each bucket needs to be `log timetamp` + `event type`
- What about changes? Thinking to NOT enable changes but rather keep appending, then render the most recent entry for that `event date` and `event type` combination.


For example: Bucket `20250101` will have notes like this.

| Log_Timestamp_Event_Type | Event_Date | Event_Content |
|---------------|------------|----------------|
|20250101_230000_read|  2024-1231 | reading for an hour |
|20250101_230000_run|  2024-1231 | half a mile run, feels good |
|...|...|...|
|Key column | NOT nullable | nullable|

If two days later another Bucket `20250103` created, the event date will duplicated on the prior Bucket `20250101`

| Log_Timestamp_Event_Type | Event_Date | Event_Content |
|---------------|------------|----------------|
|20250103_080000_read| 2024-1231 | reading for an hour |
|20250103_090000_run| 2025-0103 | half a mile run, feels good |
|...|...|...|
|Key column | NOT nullable | nullable|


### Architect trail #3

Thinking it a little more after a walk. Each `Bucket` still is named as an event date (when the event occurs), like a calendar date. 

Each `Bucket` will have a structure like below. 

| Event-Type+Log-Timestamp | Event-Time | Event_Content |
|---------------|------------|----------------|
|read + 20250103_080000| 3:00 | reading for an hour |
|run + 20250103_090000| 14:00 | half a mile run, feels good |
|bjj + 20250103_090000| afternoon | half a mile run, feels good |
|...|...|...|
|Key column | nullable | nullable|

Potentailly, it can have additional columns the key column splitted. The `sequence-number` will start from 0 by each `event-type`. It will helps rendering the results if the goal is to have the lastest view of that `event-type` within a day. 

In this case, there's no `update` to record but rather than keep appending. 

| Event-Type+Log-Timestamp | Event-Time | Event-Content | Event-Type |Seq-nbr|
|---------------|------------|----------------|--|--|
|read + 20250103_080000| 00:00 | Some reading |read| 0|
|read + 20250103_080001| 00:00 | Math reading|read|1|

The table above demonstrate that the second record of `event-type = 'read'` will be returned.

## Back up

Need to build in some logics to backup the file periodically.