CREATE DATABASE IF NOT EXISTS `auth-service`;
USE `auth-service`;
CREATE TABLE IF NOT EXISTS gender_lang (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    gender TINYINT(1) NOT NULL,
    lang VARCHAR(20) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL
);
CREATE TABLE IF NOT EXISTS auth_user (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    full_name VARCHAR(255) NOT null,
    email VARCHAR(100) NOT null,
    username VARCHAR(100) not null,
    `password` VARCHAR(255) NOT null,
    gender TINYINT(1) NOT null,
    telephone VARCHAR(15) null,
    has_deleted TINYINT(1) NOT NULL DEFAULT 0,
    picture VARCHAR(255) null,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL
);

CREATE TABLE IF NOT EXISTS auth_group (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    `name` VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL
);

CREATE TABLE IF NOT EXISTS auth_user_group (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    user_id UUID NOT null,
    group_id UUID NOT null,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL,
    FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES auth_group(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS auth_portal (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    `order` INT(3) NOT NULL,
    path VARCHAR(255) NOT NULL,
    icon VARCHAR(255) null,
    font_icon VARCHAR(50) null,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL
);
CREATE TABLE IF NOT EXISTS auth_portal_lang (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    portal_id UUID NOT null,
    `name` VARCHAR(255) NOT null,
    description TEXT null,
    lang VARCHAR(20) NOT NULL,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL,
    FOREIGN KEY (portal_id) REFERENCES auth_portal(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS auth_function (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    portal_id UUID NOT null,
    parent_id UUID null,
    method ENUM('GET', 'POST', 'PUT', 'DELETE', 'PATCH') NOT NULL DEFAULT 'GET',
    `position` VARCHAR(20) NOT null,
    icon VARCHAR(255) null,
    font_icon VARCHAR(50) null,
    is_show TINYINT(1) NOT NULL DEFAULT 1,
    shortcut_key VARCHAR(100) null,
    path VARCHAR(255) NOT null,
    `path_param` JSON NULL DEFAULT NULL,
    `order` INT(3) NOT NULL,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL,
    FOREIGN KEY (portal_id) REFERENCES auth_portal(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS auth_function_lang (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    function_id UUID NOT null,
    `name` VARCHAR(255) NOT null,
    lang VARCHAR(20) NOT null,
    description TEXT null,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL,
    FOREIGN KEY (function_id) REFERENCES auth_function(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS auth_permission (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    group_id UUID NOT null,
    function_id UUID NOT null,
    grant_create TINYINT(1) NOT NULL DEFAULT 0,
    grant_read TINYINT(1) NOT NULL DEFAULT 0,
    grant_update TINYINT(1) NOT NULL DEFAULT 0,
    grant_delete TINYINT(1) NOT NULL DEFAULT 0,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP(),
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL,
    FOREIGN KEY (group_id) REFERENCES auth_group(id) ON DELETE CASCADE,
    FOREIGN KEY (function_id) REFERENCES auth_function(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS auth_refresh_tokens (
    id UUID PRIMARY KEY NOT NULL DEFAULT (UUID()),
    user_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at DATETIME not NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    modified_at DATETIME NULL,
    modified_by VARCHAR(255) NULL,
    FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE
);
-- default group
INSERT INTO auth_group (id, `name`, description, created_by)
VALUE
(('28cfc72c-fbc5-11ef-95a5-1c697a672693'), 'Default', 'Default Group', 'System');

-- default portal

INSERT INTO auth_portal (id, `order`, path, created_at, created_by)
VALUE
(UUID(), 0, '/v1/portal', NOW(), 'MIGRATION');

INSERT INTO auth_portal_lang
(id, portal_id, `name`, description, lang, created_at, created_by)
VALUE
((UUID()), (SELECT id FROM auth_portal WHERE path = '/v1/portal'), 'Default Portal', 'Default Portal', 'en', NOW(), 'MIGRATION'),
((UUID()), (SELECT id FROM auth_portal WHERE path = '/v1/portal'), 'Portal Utama', 'Portal utama', 'id', NOW(), 'MIGRATION');
