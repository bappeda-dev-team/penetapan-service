-- =========================================================
-- SNAPSHOT INDIKATOR PROGRAM
-- =========================================================

CREATE TABLE indikator_renja_program (
    id BIGSERIAL PRIMARY KEY,

    program_id BIGINT NOT NULL,

    kode_indikator VARCHAR(255) NOT NULL,

    indikator TEXT NOT NULL,

    tahun INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_program
        FOREIGN KEY(program_id)
        REFERENCES renja_program(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_indikator_program
        UNIQUE(
            program_id,
            kode_indikator
        )
);

CREATE INDEX idx_indikator_program_tahun
ON indikator_renja_program(tahun);

CREATE INDEX idx_indikator_program_source
ON indikator_renja_program(
    kode_indikator
);

CREATE TABLE target_indikator_renja_program (
    id BIGSERIAL PRIMARY KEY,

    indikator_program_id BIGINT NOT NULL,

    kode_target VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(100) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_program
        FOREIGN KEY(indikator_program_id)
        REFERENCES indikator_renja_program(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_program
        UNIQUE(
            indikator_program_id,
            kode_target
        )
);

CREATE INDEX idx_target_program_source
ON target_indikator_renja_program(
    kode_target
);

CREATE INDEX idx_target_program_tahun
ON target_indikator_renja_program(
    tahun
);

-- =========================================================
-- SNAPSHOT INDIKATOR KEGIATAN
-- =========================================================

CREATE TABLE indikator_renja_kegiatan (
    id BIGSERIAL PRIMARY KEY,

    kegiatan_id BIGINT NOT NULL,

    kode_indikator VARCHAR(255) NOT NULL,

    indikator TEXT NOT NULL,

    tahun INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_kegiatan
        FOREIGN KEY(kegiatan_id)
        REFERENCES renja_kegiatan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_indikator_kegiatan
        UNIQUE(
            kegiatan_id,
            kode_indikator
        )
);

CREATE INDEX idx_indikator_kegiatan_tahun
ON indikator_renja_kegiatan(tahun);

CREATE INDEX idx_indikator_kegiatan_source
ON indikator_renja_kegiatan(
    kode_indikator
);

CREATE TABLE target_indikator_renja_kegiatan (
    id BIGSERIAL PRIMARY KEY,

    indikator_kegiatan_id BIGINT NOT NULL,

    kode_target VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(100) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_kegiatan
        FOREIGN KEY(indikator_kegiatan_id)
        REFERENCES indikator_renja_kegiatan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_kegiatan
        UNIQUE(
            indikator_kegiatan_id,
            kode_target
        )
);

CREATE INDEX idx_target_kegiatan_source
ON target_indikator_renja_kegiatan(
    kode_target
);

CREATE INDEX idx_target_kegiatan_tahun
ON target_indikator_renja_kegiatan(
    tahun
);

-- =========================================================
-- SNAPSHOT INDIKATOR SUBKEGIATAN
-- =========================================================

CREATE TABLE indikator_renja_subkegiatan (
    id BIGSERIAL PRIMARY KEY,

    subkegiatan_id BIGINT NOT NULL,

    kode_indikator VARCHAR(255) NOT NULL,

    indikator TEXT NOT NULL,

    tahun INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_subkegiatan
        FOREIGN KEY(subkegiatan_id)
        REFERENCES renja_subkegiatan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_indikator_subkegiatan
        UNIQUE(
            subkegiatan_id,
            kode_indikator
        )
);

CREATE INDEX idx_indikator_subkegiatan_tahun
ON indikator_renja_subkegiatan(tahun);

CREATE INDEX idx_indikator_subkegiatan_source
ON indikator_renja_subkegiatan(
    kode_indikator
);


CREATE TABLE target_indikator_renja_subkegiatan (
    id BIGSERIAL PRIMARY KEY,

    indikator_subkegiatan_id BIGINT NOT NULL,

    kode_target VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(100) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_subkegiatan
        FOREIGN KEY(indikator_subkegiatan_id)
        REFERENCES indikator_renja_subkegiatan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_subkegiatan
        UNIQUE(
            indikator_subkegiatan_id,
            kode_target
        )
);

CREATE INDEX idx_target_subkegiatan_source
ON target_indikator_renja_subkegiatan(
    kode_target
);

CREATE INDEX idx_target_subkegiatan_tahun
ON target_indikator_renja_subkegiatan(
    tahun
);
