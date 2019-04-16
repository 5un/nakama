/*
 * Copyright 2019 Signal Video Game
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

-- +migrate Up
CREATE TABLE IF NOT EXISTS password_reset_token (
    PRIMARY KEY (reset_token),
    reset_token    BYTEA         NOT NULL,
    user_id        BYTEA         NOT NULL,
    expiration     BIGINT        NOT NULL,
    created_at     BIGINT        CHECK (created_at > 0) NOT NULL,
    updated_at     BIGINT        CHECK (updated_at > 0) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS password_reset_tokens;