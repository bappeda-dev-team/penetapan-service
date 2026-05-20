-- =========================================================
-- SNAPSHOT PAGU PROGRAM
-- =========================================================

CREATE TABLE pagu_renja_program (
    id BIGSERIAL PRIMARY KEY,

    program_id BIGINT NOT NULL,

    kode_pagu VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    pagu BIGINT NOT NULL,

    jenis_pagu VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_pagu_program
        FOREIGN KEY(program_id)
        REFERENCES renja_program(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pagu_program
        UNIQUE(
            program_id,
            kode_pagu
        )
);

CREATE INDEX idx_pagu_program_source
ON pagu_renja_program(
    kode_pagu
);

CREATE INDEX idx_pagu_program_tahun
ON pagu_renja_program(
    tahun
);


-- =========================================================
-- SNAPSHOT PAGU KEGIATAN
-- =========================================================

CREATE TABLE pagu_renja_kegiatan (
    id BIGSERIAL PRIMARY KEY,

    kegiatan_id BIGINT NOT NULL,

    kode_pagu VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    pagu BIGINT NOT NULL,

    jenis_pagu VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_pagu_kegiatan
        FOREIGN KEY(kegiatan_id)
        REFERENCES renja_kegiatan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pagu_kegiatan
        UNIQUE(
            kegiatan_id,
            kode_pagu
        )
);

CREATE INDEX idx_pagu_kegiatan_source
ON pagu_renja_kegiatan(
    kode_pagu
);

CREATE INDEX idx_pagu_kegiatan_tahun
ON pagu_renja_kegiatan(
    tahun
);


-- =========================================================
-- SNAPSHOT PAGU SUBKEGIATAN
-- =========================================================

CREATE TABLE pagu_renja_subkegiatan (
    id BIGSERIAL PRIMARY KEY,

    subkegiatan_id BIGINT NOT NULL,

    kode_pagu VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    pagu BIGINT NOT NULL,

    jenis_pagu VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_pagu_subkegiatan
        FOREIGN KEY(subkegiatan_id)
        REFERENCES renja_subkegiatan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pagu_subkegiatan
        UNIQUE(
            subkegiatan_id,
            kode_pagu
        )
);

CREATE INDEX idx_pagu_subkegiatan_source
ON pagu_renja_subkegiatan(
    kode_pagu
);

CREATE INDEX idx_pagu_subkegiatan_tahun
ON pagu_renja_subkegiatan(
    tahun
);
