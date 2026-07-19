-- =========================================================
-- RENAKSI INDIVIDU
-- =========================================================

CREATE TABLE renaksi_individu (
    id BIGSERIAL PRIMARY KEY,

    pk_individu_id BIGINT NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    kode_renaksi VARCHAR(255) NOT NULL,

    urutan INTEGER NOT NULL,

    nama_rencana_aksi TEXT NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renaksi_individu_pk
        FOREIGN KEY (pk_individu_id)
        REFERENCES pk_individu(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renaksi_individu
        UNIQUE (
            pk_individu_id,
            kode_renaksi
        )
);

CREATE INDEX idx_renaksi_individu_parent
    ON renaksi_individu (pk_individu_id);

CREATE INDEX idx_renaksi_individu_lookup
    ON renaksi_individu (
        kode_opd,
        tahun_aktif
    );


-- =========================================================
-- PELAKSANAAN RENAKSI INDIVIDU
-- =========================================================

CREATE TABLE renaksi_individu_pelaksanaan (
    id BIGSERIAL PRIMARY KEY,

    renaksi_individu_id BIGINT NOT NULL,

    bulan INTEGER NOT NULL,

    bobot INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renaksi_individu_pelaksanaan
        FOREIGN KEY (renaksi_individu_id)
        REFERENCES renaksi_individu(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_renaksi_individu_pelaksanaan
        UNIQUE (
            renaksi_individu_id,
            bulan
        )
);

CREATE INDEX idx_renaksi_individu_pelaksanaan_parent
    ON renaksi_individu_pelaksanaan (renaksi_individu_id);
