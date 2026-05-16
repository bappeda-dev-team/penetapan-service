-- =========================================================
-- SNAPSHOT URUSAN
-- =========================================================

CREATE TABLE renja_urusan (
    id BIGSERIAL PRIMARY KEY,

    penetapan_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_urusan VARCHAR(255) NOT NULL,

    urusan TEXT NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renja_urusan_snapshot
        FOREIGN KEY (penetapan_id)
        REFERENCES penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renja_urusan
        UNIQUE(
            penetapan_id,
            kode_urusan
        )
);

CREATE INDEX idx_renja_urusan_snapshot
ON renja_urusan(
    penetapan_id
);

CREATE INDEX idx_renja_urusan_lookup
ON renja_urusan(
    kode_opd,
    tahun_aktif
);



-- =========================================================
-- SNAPSHOT BIDANG URUSAN
-- =========================================================

CREATE TABLE renja_bidang_urusan (
    id BIGSERIAL PRIMARY KEY,

    penetapan_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_urusan VARCHAR(255) NOT NULL,

    kode_bidang_urusan VARCHAR(255) NOT NULL,

    bidang_urusan TEXT NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renja_bidang_snapshot
        FOREIGN KEY (penetapan_id)
        REFERENCES penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renja_bidang
        UNIQUE(
            penetapan_id,
            kode_bidang_urusan
        )
);

CREATE INDEX idx_renja_bidang_snapshot
ON renja_bidang_urusan(
    penetapan_id
);

CREATE INDEX idx_renja_bidang_lookup
ON renja_bidang_urusan(
    kode_opd,
    tahun_aktif
);

CREATE INDEX idx_renja_bidang_parent
ON renja_bidang_urusan(
    kode_urusan
);



-- =========================================================
-- SNAPSHOT PROGRAM
-- =========================================================

CREATE TABLE renja_program (
    id BIGSERIAL PRIMARY KEY,

    penetapan_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_bidang_urusan VARCHAR(255) NOT NULL,

    kode_program VARCHAR(255) NOT NULL,

    program TEXT NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renja_program_snapshot
        FOREIGN KEY (penetapan_id)
        REFERENCES penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renja_program
        UNIQUE(
            penetapan_id,
            kode_program
        )
);

CREATE INDEX idx_renja_program_snapshot
ON renja_program(
    penetapan_id
);

CREATE INDEX idx_renja_program_lookup
ON renja_program(
    kode_opd,
    tahun_aktif
);

CREATE INDEX idx_renja_program_parent
ON renja_program(
    kode_bidang_urusan
);



-- =========================================================
-- SNAPSHOT KEGIATAN
-- =========================================================

CREATE TABLE renja_kegiatan (
    id BIGSERIAL PRIMARY KEY,

    penetapan_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_program VARCHAR(255) NOT NULL,

    kode_kegiatan VARCHAR(255) NOT NULL,

    kegiatan TEXT NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renja_kegiatan_snapshot
        FOREIGN KEY (penetapan_id)
        REFERENCES penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renja_kegiatan
        UNIQUE(
            penetapan_id,
            kode_kegiatan
        )
);

CREATE INDEX idx_renja_kegiatan_snapshot
ON renja_kegiatan(
    penetapan_id
);

CREATE INDEX idx_renja_kegiatan_lookup
ON renja_kegiatan(
    kode_opd,
    tahun_aktif
);

CREATE INDEX idx_renja_kegiatan_parent
ON renja_kegiatan(
    kode_program
);



-- =========================================================
-- SNAPSHOT SUBKEGIATAN
-- =========================================================

CREATE TABLE renja_subkegiatan (
    id BIGSERIAL PRIMARY KEY,

    penetapan_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_kegiatan VARCHAR(255) NOT NULL,

    kode_subkegiatan VARCHAR(255) NOT NULL,

    subkegiatan TEXT NOT NULL,

    pegawai_id VARCHAR(100),

    nama_pegawai TEXT,

    total_anggaran BIGINT NOT NULL DEFAULT 0,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renja_subkegiatan_snapshot
        FOREIGN KEY (penetapan_id)
        REFERENCES penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renja_subkegiatan
        UNIQUE(
            penetapan_id,
            kode_subkegiatan
        )
);

CREATE INDEX idx_renja_subkegiatan_snapshot
ON renja_subkegiatan(
    penetapan_id
);

CREATE INDEX idx_renja_subkegiatan_lookup
ON renja_subkegiatan(
    kode_opd,
    tahun_aktif
);

CREATE INDEX idx_renja_subkegiatan_parent
ON renja_subkegiatan(
    kode_kegiatan
);
