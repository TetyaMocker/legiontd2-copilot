package storage

const migrationInit = `
CREATE TABLE IF NOT EXISTS matches (
    id              TEXT PRIMARY KEY,
    started_at      DATETIME NOT NULL,
    ended_at        DATETIME,
    patch_version   TEXT,
    result          TEXT
);

CREATE TABLE IF NOT EXISTS wave_snapshots (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    match_id            TEXT NOT NULL REFERENCES matches(id),
    wave_number         INTEGER NOT NULL,
    captured_at         DATETIME NOT NULL,
    mythium             INTEGER,
    income              INTEGER,
    king_hp_percent     INTEGER,
    ally_king_hp_percent INTEGER,
    confidence          REAL,
    source              TEXT NOT NULL DEFAULT 'hudapi'
);

CREATE TABLE IF NOT EXISTS recommendations (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    wave_snapshot_id INTEGER NOT NULL REFERENCES wave_snapshots(id),
    kind            TEXT NOT NULL,
    explanation     TEXT NOT NULL,
    accepted        BOOLEAN
);

CREATE TABLE IF NOT EXISTS units_reference (
    unit_id         TEXT NOT NULL,
    patch_version   TEXT NOT NULL,
    name            TEXT,
    mythium_cost    INTEGER,
    gold_cost       INTEGER,
    hp              INTEGER,
    unit_class      TEXT,
    raw_json        TEXT NOT NULL,
    fetched_at      DATETIME NOT NULL,
    PRIMARY KEY (unit_id, patch_version)
);
`
