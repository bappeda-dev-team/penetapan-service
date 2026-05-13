-- =========================================================
-- SNAPSHOT METADATA
-- =========================================================

CREATE TABLE sync_penetapan_metadata_opd (
    id BIGSERIAL PRIMARY KEY,

    kode_opd VARCHAR(50) NOT NULL,

    tahun INTEGER NOT NULL,

    jenis_penetapan VARCHAR(30) NOT NULL,

    status VARCHAR(30) NOT NULL,

    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    finished_at TIMESTAMPTZ,

    sync_by VARCHAR(100),

    error_message TEXT,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sync_penetapan_metadata
ON sync_penetapan_metadata_opd(
    kode_opd,
    tahun,
    jenis_penetapan
);

-- =========================================================
-- ROOT SNAPSHOT
-- =========================================================

CREATE TABLE penetapan_opd (
    id BIGSERIAL PRIMARY KEY,

    kode_opd VARCHAR(50) NOT NULL,

    tahun INTEGER NOT NULL,

    jenis_penetapan VARCHAR(30) NOT NULL,

    versi INTEGER NOT NULL DEFAULT 1,

    snapshot_status VARCHAR(30) NOT NULL,

    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    generated_by VARCHAR(100),

    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- hanya boleh ada satu snapshot aktif
CREATE UNIQUE INDEX uq_penetapan_active
ON penetapan_opd(
    kode_opd,
    tahun,
    jenis_penetapan
)
WHERE is_active = TRUE;

-- versi tidak boleh duplicate
CREATE UNIQUE INDEX uq_penetapan_version
ON penetapan_opd(
    kode_opd,
    tahun,
    jenis_penetapan,
    versi
);

CREATE INDEX idx_penetapan_lookup
ON penetapan_opd(
    kode_opd,
    tahun,
    jenis_penetapan
);
