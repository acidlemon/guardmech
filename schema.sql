-- authorization is reserved word...
CREATE TABLE IF NOT EXISTS auth (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    principal_id BIGINT NOT NULL,
    issuer VARCHAR(255) CHARACTER SET utf8 NOT NULL,
    subject VARCHAR(255) CHARACTER SET utf8 NOT NULL,
    email VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
    PRIMARY KEY (seq_id),
    UNIQUE uniq_issuer_subject (issuer,subject)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS api_key (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    principal_id BIGINT NOT NULL,
    token VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS principal (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

-- group is reserved word...
CREATE TABLE IF NOT EXISTS group_info (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS group_rule (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    rule_type INTEGER NOT NULL, /* 1=specific domain, 2=whole domain, 3=member of  */
    condition TEXT(4096),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

-- role is reserved word...
CREATE TABLE IF NOT EXISTS role_info (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS permission (
    seq_id BIGINT NOT NULL auto_increment,
    unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
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


