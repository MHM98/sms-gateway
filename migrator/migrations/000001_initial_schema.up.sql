CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS wallets (
    user_id BIGINT UNSIGNED NOT NULL,
    balance BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_at DATETIME(6) NOT NULL
        DEFAULT CURRENT_TIMESTAMP(6)
        ON UPDATE CURRENT_TIMESTAMP(6),

    PRIMARY KEY (user_id),

    CONSTRAINT fk_wallets_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
);


CREATE TABLE IF NOT EXISTS message_requests (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    idempotency_key VARCHAR(64) NOT NULL,
    payload_hash BINARY(32) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    PRIMARY KEY (id),

    UNIQUE KEY uq_message_requests_idempotency (
        user_id,
        idempotency_key
    ),

    CONSTRAINT fk_message_requests_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
);


CREATE TABLE IF NOT EXISTS messages (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    request_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,

    recipient VARCHAR(20) NOT NULL,
    body VARCHAR(255) NOT NULL,

    service_type ENUM('normal', 'express')
        NOT NULL DEFAULT 'normal',

    status ENUM('pending', 'processing', 'submitted', 'failed')
        NOT NULL DEFAULT 'pending',

    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    submitted_to_operator_at DATETIME(6) NULL,

    PRIMARY KEY (id, created_at),
    
    INDEX idx_messages_request (
        request_id
    ),

    INDEX idx_messages_user_report (
        user_id,
        created_at DESC,
        id DESC
    ),

    INDEX idx_messages_pending (
        status,
        created_at,
        id
    )
)
PARTITION BY RANGE COLUMNS (created_at) (
    PARTITION p_before_20260715
        VALUES LESS THAN ('2026-07-15'),

    PARTITION p20260715
        VALUES LESS THAN ('2026-07-16'),

    PARTITION p20260716
        VALUES LESS THAN ('2026-07-17'),

    PARTITION p20260717
        VALUES LESS THAN ('2026-07-18'),

    PARTITION p20260718
        VALUES LESS THAN ('2026-07-19'),

    PARTITION p20260719
        VALUES LESS THAN ('2026-07-20'),

    PARTITION p20260720
        VALUES LESS THAN ('2026-07-21'),

    PARTITION p20260721
        VALUES LESS THAN ('2026-07-22'),

    PARTITION p_future
        VALUES LESS THAN (MAXVALUE)
);