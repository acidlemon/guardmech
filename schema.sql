-- authorization is reserved word...
CREATE TABLE IF NOT EXISTS auth_oidc (
    seq_id BIGINT NOT NULL auto_increment,
    auth_oidc_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    principal_seq_id BIGINT NOT NULL,
    issuer VARCHAR(255) CHARACTER SET utf8 NOT NULL,
    subject VARCHAR(255) CHARACTER SET utf8 NOT NULL,
    email VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL,
    PRIMARY KEY (seq_id),
    UNIQUE uniq_issuer_subject (issuer,subject)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS auth_apikey (
    seq_id BIGINT NOT NULL auto_increment,
    auth_apikey_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    principal_seq_id BIGINT NOT NULL,
    name VARCHAR(191) NOT NULL,
    masked_token VARCHAR(255) CHARACTER SET utf8 NOT NULL,
    salt VARCHAR(255) CHARACTER SET utf8 NOT NULL,
    hashed_token VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS principal (
    seq_id BIGINT NOT NULL auto_increment,
    principal_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

-- group is reserved word...
CREATE TABLE IF NOT EXISTS group_info (
    seq_id BIGINT NOT NULL auto_increment,
    group_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS group_rule (
    seq_id BIGINT NOT NULL auto_increment,
    group_rule_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    rule_type INTEGER NOT NULL, /* 1=specific domain, 2=whole domain, 3=member of 4=specific */
    condition TEXT(4096),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

-- role is reserved word...
CREATE TABLE IF NOT EXISTS role_info (
    seq_id BIGINT NOT NULL auto_increment,
    role_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS permission (
    seq_id BIGINT NOT NULL auto_increment,
    permission_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS mapping_rule (
    seq_id BIGINT NOT NULL auto_increment,
    mapping_rule_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
    rule_type INT NOT NULL, -- 1=specific domain, 2=whole domain, 3=member of 4=specific 
    detail VARCHAR(255) NOT NULL,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT(2048),
    priority INT NOT NULL,
    association_type INT NOT NULL, -- 1=group, 2=role
    association_id VARCHAR(40) CHARACTER SET latin1 NOT NULL,
    PRIMARY KEY (seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;


-- N:M relations

CREATE TABLE IF NOT EXISTS principal_group_map (
    principal_seq_id BIGINT NOT NULL,
    group_seq_id BIGINT NOT NULL,
    manager BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (principal_seq_id, group_seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS principal_role_map (
    principal_seq_id BIGINT NOT NULL,
    role_seq_id BIGINT NOT NULL,
    PRIMARY KEY (principal_seq_id, role_seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS group_role_map (
    group_seq_id BIGINT NOT NULL,
    role_seq_id BIGINT NOT NULL,
    PRIMARY KEY (group_seq_id, role_seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS role_permission_map (
    role_seq_id BIGINT NOT NULL,
    permission_seq_id BIGINT NOT NULL,
    PRIMARY KEY (role_seq_id, permission_seq_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;


