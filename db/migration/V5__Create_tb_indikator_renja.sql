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
            kode_indikator,
            tahun
        )
);

CREATE INDEX idx_indikator_program_parent
ON indikator_renja_program(program_id);

CREATE INDEX idx_indikator_program_tahun
ON indikator_renja_program(tahun);



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
            tahun
        )
);

CREATE INDEX idx_target_program_parent
ON target_indikator_renja_program(
    indikator_program_id
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
            kode_indikator,
            tahun
        )
);

CREATE INDEX idx_indikator_kegiatan_parent
ON indikator_renja_kegiatan(
    kegiatan_id
);

CREATE INDEX idx_indikator_kegiatan_tahun
ON indikator_renja_kegiatan(
    tahun
);



CREATE TABLE target_indikator_renja_kegiatan (
    id BIGSERIAL PRIMARY KEY,

    indikator_kegiatan_id BIGINT NOT NULL,

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
            tahun
        )
);

CREATE INDEX idx_target_kegiatan_parent
ON target_indikator_renja_kegiatan(
    indikator_kegiatan_id
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
            kode_indikator,
            tahun
        )
);

CREATE INDEX idx_indikator_subkegiatan_parent
ON indikator_renja_subkegiatan(
    subkegiatan_id
);

CREATE INDEX idx_indikator_subkegiatan_tahun
ON indikator_renja_subkegiatan(
    tahun
);



CREATE TABLE target_indikator_renja_subkegiatan (
    id BIGSERIAL PRIMARY KEY,

    indikator_subkegiatan_id BIGINT NOT NULL,

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
            tahun
        )
);

CREATE INDEX idx_target_subkegiatan_parent
ON target_indikator_renja_subkegiatan(
    indikator_subkegiatan_id
);
