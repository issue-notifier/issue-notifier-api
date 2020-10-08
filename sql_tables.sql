create table github_user(
    user_id uuid primary key default uuid_generate_v4(),
    username varchar(40),
    email varchar(50)
);

create table global_repository(
    repo_id uuid primary key default uuid_generate_v4(),
    repo_name varchar(255)
);

create table user_subscription(
    sub_id uuid primary key default uuid_generate_v4(),
    user_id uuid,
    repo_id uuid,
    labels jsonb not null,
    last_notification_sent timestamptz,
    unique (user_id, repo_id),
    CONSTRAINT fk_user_id
      FOREIGN KEY(user_id) 
	  REFERENCES github_user(user_id),
    CONSTRAINT fk_repo_id
      FOREIGN KEY(repo_id) 
	  REFERENCES global_repository(repo_id)
);

--                 sub_id                |               user_id                |               repo_id                |                                                                                                                    labels                                                                                                                    | last_notification_sent 
-- --------------------------------------+--------------------------------------+--------------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+------------------------
--  701807c6-40df-47ee-bdfd-0077f80c22f1 | 67d3571d-03dc-495d-8b03-985892406a7f | b4cdedf9-8b50-47f8-a44d-9d8af7c7e3cc | [{"name": "Type: Enhancement", "color": "#84b6eb"}, {"name": "Type: Feature Request", "color": "#c7def8"}, {"name": "Type: Release", "color": "#00D8EA"}]                                                                                    | 
--  2b786636-abb0-4c86-80ff-ec7d5e216043 | 67d3571d-03dc-495d-8b03-985892406a7f | 22783798-ae4f-425b-95e0-74bb978056f6 | [{"name": "Proposal-Accepted", "color": "#009800"}, {"name": "Soon", "color": "#b60205"}, {"name": "UX", "color": "#5d53ed"}]                                                                                                                | 
--  2b5ce45e-744b-4ddf-9239-a08158f9dd9b | 67d3571d-03dc-495d-8b03-985892406a7f | b2a476b6-629c-455a-9ab4-6c6182134177 | [{"name": "Bug", "color": "#e11d21"}, {"name": "Domain: Index Types", "color": "#f7b7e9"}, {"name": "Domain: Quick Fixes", "color": "#d4c5f9"}, {"name": "High Priority", "color": "#e11d21"}, {"name": "Out of Scope", "color": "#556677"}] | 
--  73e24975-8a84-404c-b243-f158dbdd746a | 67d3571d-03dc-495d-8b03-985892406a7f | 57dbc88e-cda8-4951-bd7e-fdc76f0fb2dd | [{"name": "bug", "color": "#fc2929"}, {"name": "question", "color": "#cc317c"}]                                                                                                                                                              | 


update user_subscription 
set labels = labels || '[{"name":"testupdate", "color":"#fff"},{"name":"testupdate2","color":"#ccc"}]'::jsonb
where user_id = '67d3571d-03dc-495d-8b03-985892406a7f' and repo_id = '57dbc88e-cda8-4951-bd7e-fdc76f0fb2dd';

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