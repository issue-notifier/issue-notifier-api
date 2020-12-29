sudo -i -u postgres -- get into postgres user

psql -- start postgres

\connect issue_notifer -- connect to the database

\d -- show tables

Ctrl + L -- clear screen

create table github_user(
    user_id uuid primary key default uuid_generate_v4(),
    username varchar(40),
    email varchar(50),
    last_notification_sent timestamptz
);

create table global_repository(
    repo_id uuid primary key default uuid_generate_v4(),
    repo_name varchar(255),
    last_event_time timestamptz,
    last_event_fetched_time timestamptz
);

create table user_subscription(
    sub_id uuid primary key default uuid_generate_v4(),
    user_id uuid,
    repo_id uuid,
    labels jsonb not null,
    last_notification_time timestamptz,
    unique (user_id, repo_id),
    CONSTRAINT fk_user_id
      FOREIGN KEY(user_id) 
	  REFERENCES github_user(user_id),
    CONSTRAINT fk_repo_id
      FOREIGN KEY(repo_id) 
	  REFERENCES global_repository(repo_id)
);

-- this table will be cleaned up regularly, only notifications which are left to be sent will be stored in this table
create table notifications(
    notif_id uuid primary key default uuid_generate_v4(),
    sent bool default 'false',

);

create table notification_data (
    notif_data_id uuid primary key default uuid_generate_v4(),
    user_id uuid,
    repo_id uuid,
    issues jsonb not NULL,
    sent bool default 'false',
    -- freq INTEGER, -- immediate - 0, daily - 1, weekly - 2, monthly - 3, lets make these codes -- this will come under user data, also keep data of when last notif was sent to that user
)

--                 sub_id                |               user_id                |               repo_id                |                                                                                                                    labels                                                                                                                    | last_notification_sent 
-- --------------------------------------+--------------------------------------+--------------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+------------------------
--  701807c6-40df-47ee-bdfd-0077f80c22f1 | 67d3571d-03dc-495d-8b03-985892406a7f | b4cdedf9-8b50-47f8-a44d-9d8af7c7e3cc | [{"name": "Type: Enhancement", "color": "#84b6eb"}, {"name": "Type: Feature Request", "color": "#c7def8"}, {"name": "Type: Release", "color": "#00D8EA"}]                                                                                    | 
--  2b786636-abb0-4c86-80ff-ec7d5e216043 | 67d3571d-03dc-495d-8b03-985892406a7f | 22783798-ae4f-425b-95e0-74bb978056f6 | [{"name": "Proposal-Accepted", "color": "#009800"}, {"name": "Soon", "color": "#b60205"}, {"name": "UX", "color": "#5d53ed"}]                                                                                                                | 
--  2b5ce45e-744b-4ddf-9239-a08158f9dd9b | 67d3571d-03dc-495d-8b03-985892406a7f | b2a476b6-629c-455a-9ab4-6c6182134177 | [{"name": "Bug", "color": "#e11d21"}, {"name": "Domain: Index Types", "color": "#f7b7e9"}, {"name": "Domain: Quick Fixes", "color": "#d4c5f9"}, {"name": "High Priority", "color": "#e11d21"}, {"name": "Out of Scope", "color": "#556677"}] | 
--  73e24975-8a84-404c-b243-f158dbdd746a | 67d3571d-03dc-495d-8b03-985892406a7f | 57dbc88e-cda8-4951-bd7e-fdc76f0fb2dd | [{"name": "bug", "color": "#fc2929"}, {"name": "question", "color": "#cc317c"}]                                                                                                                                                              | 


update user_subscription 
set labels = labels || '[{"name":"Type: Bug", "color":"#b60205"},{"name":"Resolution: Needs More Information","color":"#fffde7"}]'::jsonb
where user_id = '56cfc307-d721-42b1-bd67-31dade89678e' and repo_id = '890f3a1e-ec1b-4ad1-af53-229ea8acdfcc';

update user_subscription 
set labels = labels - (
    select i from generate_series(0, jsonb_array_length(labels) - 1) as i
    where labels->i->>'name' = 'UX'
)::integer - (
    select i from generate_series(0, jsonb_array_length(labels) - 1) as i
    where labels->i->>'name' = 'Soon'
)::integer
where sub_id = '2b786636-abb0-4c86-80ff-ec7d5e216043';

UPDATE USER_SUBSCRIPTION 
SET LABELS = LABELS - (
    SELECT I FROM GENERATE_SERIES(0, JSONB_ARRAY_LENGTH(LABELS) - 1) AS I 
    WHERE LABELS->I->>'name' = 'Type: Release'
)::INTEGER - (
    SELECT I FROM GENERATE_SERIES(0, JSONB_ARRAY_LENGTH(LABELS) - 1) AS I 
    WHERE LABELS->I->>'name' = 'Type: Enhancement'
)::INTEGER 
WHERE USER_ID = '67d3571d-03dc-495d-8b03-985892406a7f' AND REPO_ID = 'b4cdedf9-8b50-47f8-a44d-9d8af7c7e3cc';

--2b786636-abb0-4c86-80ff-ec7d5e216043 | 67d3571d-03dc-495d-8b03-985892406a7f | 22783798-ae4f-425b-95e0-74bb978056f6 | [{"name": "Proposal-Accepted", "color": "#009800"}, {"name": "Soon", "color": "#b60205"}, {"name": "UX", "color": "#5d53ed"}] 


select * 
from user_subscription,jsonb_to_recordset(user_subscription.labels) 
as labels(name text);

--                sub_id                |               user_id                |               repo_id                |                                                                                      labels                                                                                      | last_notification_sent |           name           
----------------------------------------+--------------------------------------+--------------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+------------------------+--------------------------
-- 73ad048c-eb39-4543-beea-fb6df5980ac0 | 67d3571d-03dc-495d-8b03-985892406a7f | 5fd8702f-51e3-45b6-8a53-8de45028e2cc | [{"name": "Documentation", "color": "#aaffaa"}, {"name": "Soon", "color": "#b60205"}, {"name": "cla: yes", "color": "#0e8a16"}, {"name": "release-blocker", "color": "#b60205"}] |                        | Documentation            
-- 73ad048c-eb39-4543-beea-fb6df5980ac0 | 67d3571d-03dc-495d-8b03-985892406a7f | 5fd8702f-51e3-45b6-8a53-8de45028e2cc | [{"name": "Documentation", "color": "#aaffaa"}, {"name": "Soon", "color": "#b60205"}, {"name": "cla: yes", "color": "#0e8a16"}, {"name": "release-blocker", "color": "#b60205"}] |                        | Soon                      
-- 73ad048c-eb39-4543-beea-fb6df5980ac0 | 67d3571d-03dc-495d-8b03-985892406a7f | 5fd8702f-51e3-45b6-8a53-8de45028e2cc | [{"name": "Documentation", "color": "#aaffaa"}, {"name": "Soon", "color": "#b60205"}, {"name": "cla: yes", "color": "#0e8a16"}, {"name": "release-blocker", "color": "#b60205"}] |                        | cla: yes                  
-- 73ad048c-eb39-4543-beea-fb6df5980ac0 | 67d3571d-03dc-495d-8b03-985892406a7f | 5fd8702f-51e3-45b6-8a53-8de45028e2cc | [{"name": "Documentation", "color": "#aaffaa"}, {"name": "Soon", "color": "#b60205"}, {"name": "cla: yes", "color": "#0e8a16"}, {"name": "release-blocker", "color": "#b60205"}] |                        | release-blocker          
-- 0a36ebec-ccd3-476b-9a8a-605c5ff6100b | 67d3571d-03dc-495d-8b03-985892406a7f | 890f3a1e-ec1b-4ad1-af53-229ea8acdfcc | [{"name": "Difficulty: starter", "color": "#94ce52"}, {"name": "good first issue", "color": "#6ce26a"}, {"name": "good first issue (taken)", "color": "#b60205"}]                |                        | Difficulty: starter       
-- 0a36ebec-ccd3-476b-9a8a-605c5ff6100b | 67d3571d-03dc-495d-8b03-985892406a7f | 890f3a1e-ec1b-4ad1-af53-229ea8acdfcc | [{"name": "Difficulty: starter", "color": "#94ce52"}, {"name": "good first issue", "color": "#6ce26a"}, {"name": "good first issue (taken)", "color": "#b60205"}]                |                        | good first issue          
-- 0a36ebec-ccd3-476b-9a8a-605c5ff6100b | 67d3571d-03dc-495d-8b03-985892406a7f | 890f3a1e-ec1b-4ad1-af53-229ea8acdfcc | [{"name": "Difficulty: starter", "color": "#94ce52"}, {"name": "good first issue", "color": "#6ce26a"}, {"name": "good first issue (taken)", "color": "#b60205"}]                |                        | good first issue (taken)  


 select labels.name 
 from user_subscription,jsonb_to_recordset(user_subscription.labels) 
 as labels(name text) 
 where repo_id='5fd8702f-51e3-45b6-8a53-8de45028e2cc';

--       name       
-----------------
-- Documentation
-- Soon
-- cla: yes
-- release-blocker


select distinct labels.name, labels.color 
from user_subscription,jsonb_to_recordset(user_subscription.labels) 
as labels(name text, color text) 
where repo_id='5fd8702f-51e3-45b6-8a53-8de45028e2cc';

--      name       |  color  
-------------------+---------
-- Documentation   | #aaffaa
-- Soon            | #b60205
-- cla: yes        | #0e8a16
-- release-blocker | #b60205
--(4 rows)
