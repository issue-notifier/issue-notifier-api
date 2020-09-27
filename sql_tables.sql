create table github_user(
    user_id uuid primary key default uuid_generate_v4(),
    username varchar(40),
    email varchar(50)
);

create table global_repository(
    repo_id uuid primary key default uuid_generate_v4(),
    repo_name varchar(255),
    api_url varchar(255),
    html_url varchar(255)
);

create table user_subscription(
    sub_id uuid primary key default uuid_generate_v4(),
    user_id uuid,
    repo_id uuid,
    label varchar(255),
    last_notification_sent timestamptz,
    unique (user_id, repo_id, label),
    CONSTRAINT fk_user_id
      FOREIGN KEY(user_id) 
	  REFERENCES github_user(user_id),
    CONSTRAINT fk_repo_id
      FOREIGN KEY(repo_id) 
	  REFERENCES global_repository(repo_id)
);