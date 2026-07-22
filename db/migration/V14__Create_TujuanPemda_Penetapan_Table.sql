-- =========================================================
-- SNAPSHOT TUJUAN PEMDA
-- =========================================================

CREATE TABLE tujuan_pemda_penetapan (
    id BIGSERIAL PRIMARY KEY,

    penetapan_pemda_id BIGINT NOT NULL,

    kode_tujuan_pemda VARCHAR(255) NOT NULL,

    tujuan_pemda TEXT NOT NULL,

    periode VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_tujuan_pemda_penetapan
        FOREIGN KEY (penetapan_pemda_id)
        REFERENCES penetapan_pemda(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_tujuan_pemda_penetapan
        UNIQUE (
            penetapan_pemda_id,
            kode_tujuan_pemda
        )
);

CREATE INDEX idx_tujuan_pemda_penetapan_parent
    ON tujuan_pemda_penetapan (penetapan_pemda_id);

CREATE INDEX idx_tujuan_pemda_penetapan_lookup
    ON tujuan_pemda_penetapan (tahun_aktif);


-- =========================================================
-- INDIKATOR TUJUAN PEMDA
-- =========================================================

CREATE TABLE indikator_tujuan_pemda_penetapan (
    id BIGSERIAL PRIMARY KEY,

    tujuan_pemda_id BIGINT NOT NULL,

    kode_indikator VARCHAR(255) NOT NULL,

    indikator TEXT NOT NULL,

    rumus_perhitungan TEXT,

    sumber_data TEXT,

    definisi_operasional TEXT,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_tujuan_pemda_penetapan
        FOREIGN KEY (tujuan_pemda_id)
        REFERENCES tujuan_pemda_penetapan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_indikator_tujuan_pemda_penetapan
        UNIQUE (
            tujuan_pemda_id,
            kode_indikator
        )
);

CREATE INDEX idx_indikator_tujuan_pemda_penetapan_parent
    ON indikator_tujuan_pemda_penetapan (tujuan_pemda_id);


-- =========================================================
-- TARGET INDIKATOR TUJUAN PEMDA
-- =========================================================

CREATE TABLE target_indikator_tujuan_pemda_penetapan (
    id BIGSERIAL PRIMARY KEY,

    indikator_tujuan_pemda_id BIGINT NOT NULL,

    kode_target VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,4) NOT NULL,

    satuan VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_indikator_tujuan_pemda_penetapan
        FOREIGN KEY (indikator_tujuan_pemda_id)
        REFERENCES indikator_tujuan_pemda_penetapan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_indikator_tujuan_pemda_penetapan
        UNIQUE (
            indikator_tujuan_pemda_id,
            kode_target
        )
);

CREATE INDEX idx_target_indikator_tujuan_pemda_penetapan_parent
    ON target_indikator_tujuan_pemda_penetapan (indikator_tujuan_pemda_id);
