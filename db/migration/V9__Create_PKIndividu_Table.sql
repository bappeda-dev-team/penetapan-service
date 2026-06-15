-- =========================================================
-- SNAPSHOT PK INDIVIDU
-- =========================================================

CREATE TABLE pk_individu (
    id BIGSERIAL PRIMARY KEY,

    penetapan_individu_id BIGINT NOT NULL,

    pegawai_id VARCHAR(100) NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    level_pk INTEGER NOT NULL,

    kode_pk VARCHAR(255) NOT NULL,

    nama_pk VARCHAR(255) NOT NULL,

    keterangan_pk TEXT,

    nama_pemilik_pk VARCHAR(255) NOT NULL,

    versi INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_pk_individu_penetapan
        FOREIGN KEY (penetapan_individu_id)
        REFERENCES penetapan_individu(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pk_individu_snapshot
        UNIQUE (
            penetapan_individu_id,
            kode_pk
        )
);

CREATE INDEX idx_pk_individu_snapshot
    ON pk_individu (penetapan_individu_id);

CREATE INDEX idx_pk_individu_lookup
    ON pk_individu (
        pegawai_id,
        kode_opd,
        tahun_aktif
    );


-- =========================================================
-- INDIKATOR PK
-- =========================================================

CREATE TABLE indikator_pk (
    id BIGSERIAL PRIMARY KEY,

    pk_individu_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    kode_indikator_pk VARCHAR(255) NOT NULL,

    nama_indikator_pk TEXT NOT NULL,

    rumus_perhitungan TEXT,

    sumber_data TEXT,

    definisi_operasional TEXT,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_pk
        FOREIGN KEY (pk_individu_id)
        REFERENCES pk_individu(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_indikator_pk
        UNIQUE (
            pk_individu_id,
            kode_indikator_pk
        )
);

CREATE INDEX idx_indikator_pk_parent
    ON indikator_pk (pk_individu_id);

CREATE INDEX idx_indikator_pk_lookup
    ON indikator_pk (
        kode_opd,
        tahun_aktif
    );


-- =========================================================
-- TARGET PK
-- =========================================================

CREATE TABLE target_pk (
    id BIGSERIAL PRIMARY KEY,

    indikator_pk_id BIGINT NOT NULL,

    kode_target_pk VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,4) NOT NULL,

    satuan VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_pk
        FOREIGN KEY (indikator_pk_id)
        REFERENCES indikator_pk(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_pk
        UNIQUE (
            indikator_pk_id,
            kode_target_pk
        )
);

CREATE INDEX idx_target_pk_parent
    ON target_pk (indikator_pk_id);
