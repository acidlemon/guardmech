-- authorization is reserved word...
CREATE TABLE IF NOT EXISTS auth (
    id BIGINT NOT NULL auto_increment,
    principal_id BIGINT NOT NULL,
    account VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
    PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS api_key (
    id BIGINT NOT NULL auto_increment,
    principal_id BIGINT NOT NULL,
    token VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
    PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS principal (
    id BIGINT NOT NULL auto_increment,
    name VARCHAR(191) NOT NULL,
    description TEXT(2048),
    PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

-- group is reserved word...
CREATE TABLE IF NOT EXISTS group_info (
    id BIGINT NOT NULL auto_increment,
    name VARCHAR(191) NOT NULL,
    description TEXT(2048),
    PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

-- role is reserved word...
CREATE TABLE IF NOT EXISTS role_info (
    id BIGINT NOT NULL auto_increment,
    name VARCHAR(191) NOT NULL,
    description TEXT(2048),
    PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS permission (
    id BIGINT NOT NULL auto_increment,
    name VARCHAR(191) NOT NULL,
    description TEXT(2048),
    PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;


-- N:M relations

CREATE TABLE IF NOT EXISTS principal_group_map (
    principal_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    manager BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (principal_id, group_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS principal_role_map (
    principal_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (principal_id, role_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS group_role_map (
    group_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (group_id, role_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS role_permission_map (
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, permission_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;


