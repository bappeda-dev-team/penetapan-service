-- =========================================================
-- SNAPSHOT METADATA
-- =========================================================
CREATE TABLE sync_penetapan_metadata_individu(

    id BIGSERIAL PRIMARY KEY,

    pegawai_id VARCHAR(100) NOT NULL,

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

CREATE INDEX idx_sync_penetapan_individu_metadata
ON sync_penetapan_metadata_individu(
    pegawai_id,
    kode_opd,
    tahun,
    jenis_penetapan
);

-- =========================================================
-- ROOT SNAPSHOT INDIVIDU
-- =========================================================
CREATE TABLE penetapan_individu (
    id BIGSERIAL PRIMARY KEY,

    pegawai_id VARCHAR(100) NOT NULL,

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
CREATE UNIQUE INDEX uq_penetapan_individu_active
ON penetapan_individu(
    pegawai_id,
    kode_opd,
    tahun,
    jenis_penetapan
)
WHERE is_active = TRUE;

-- versi tidak boleh duplicate
CREATE UNIQUE INDEX uq_penetapan_individu_version
ON penetapan_individu(
    pegawai_id,
    kode_opd,
    tahun,
    jenis_penetapan,
    versi
);

CREATE INDEX idx_penetapan_individu_lookup
ON penetapan_individu(
    pegawai_id,
    kode_opd,
    tahun,
    jenis_penetapan
);
