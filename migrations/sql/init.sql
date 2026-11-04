--liquibase formatted sql

--changeset kbar:init.sql

CREATE SCHEMA users;

CREATE TABLE users.users (
    id  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username text NOT NULL UNIQUE,
    created_at timestamp DEFAULT now()
);

CREATE TABLE users.groups (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    created_by uuid,
    created_at timestamp DEFAULT now()
);

CREATE TABLE users.users_groups (
    user_id UUID NOT NULL REFERENCES users.users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES users.users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT now(),
    
    -- Составной первичный ключ гарантирует, что пользователь 
    -- не сможет вступить в одну и ту же группу дважды
    PRIMARY KEY (user_id, group_id)
);