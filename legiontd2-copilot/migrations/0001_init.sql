-- Runtime: локальная история собственных матчей (см. ТЗ, DR3/DR5, FR4)
CREATE TABLE IF NOT EXISTS matches (
    id              TEXT PRIMARY KEY,          -- локальный uuid, НЕ id из API v2
    started_at      DATETIME NOT NULL,
    ended_at        DATETIME,
    patch_version   TEXT,
    result          TEXT                       -- 'won' | 'lost' | 'unknown' пока матч не завершён
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
    confidence          REAL,                  -- из EconomyState.confidence (Perception Service)
    source              TEXT NOT NULL DEFAULT 'ocr' -- 'ocr' | 'overwolf' | 'manual' — см. ТЗ 8.4/Phase 0.5
);

CREATE TABLE IF NOT EXISTS recommendations (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    wave_snapshot_id INTEGER NOT NULL REFERENCES wave_snapshots(id),
    kind            TEXT NOT NULL,             -- 'spend' | 'save' | 'mercenary' | 'placement'
    explanation     TEXT NOT NULL,             -- человекочитаемое объяснение (AI3 в ТЗ)
    accepted        BOOLEAN                    -- заполняется постфактум для будущей оценки качества советов
);

-- Офлайн: локальный кэш справочников из официального API v2 (/units, /info)
-- Обновляется отдельным батч-джобом (ТЗ, Phase 1.2), не в рантайме матча.
CREATE TABLE IF NOT EXISTS units_reference (
    unit_id         TEXT NOT NULL,
    patch_version   TEXT NOT NULL,
    name            TEXT,
    mythium_cost    INTEGER,
    gold_cost       INTEGER,
    hp              INTEGER,
    unit_class      TEXT,                      -- Creature | Fighter | King | Mercenary | Worker
    raw_json        TEXT NOT NULL,              -- полный ответ API на случай новых полей
    fetched_at      DATETIME NOT NULL,
    PRIMARY KEY (unit_id, patch_version)
);
