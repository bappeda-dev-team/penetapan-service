CREATE TABLE rencana_aksi_opd (
    id                    BIGSERIAL PRIMARY KEY,
    kode_rencana_aksi_opd VARCHAR(255) NOT NULL,
    penetapan_id          BIGINT NOT NULL,

    kode_opd              VARCHAR(100) NOT NULL,
    kode_sasaran_opd      VARCHAR(100) NOT NULL,
    kode_pk               VARCHAR(100) NOT NULL,
    nama_pk               VARCHAR(500) NOT NULL,
    pegawai_id            VARCHAR(100) NOT NULL,
    kode_subkegiatan      VARCHAR(100) NOT NULL,
    nama_subkegiatan      TEXT NOT NULL,

    anggaran_renaksi      BIGINT NOT NULL DEFAULT 0,
    tahun                 INTEGER NOT NULL,

    tw1                   INTEGER NOT NULL DEFAULT 0,
    tw2                   INTEGER NOT NULL DEFAULT 0,
    tw3                   INTEGER NOT NULL DEFAULT 0,
    tw4                   INTEGER NOT NULL DEFAULT 0,

    created_date          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by            VARCHAR(255),
    versi                 INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT uk_rencana_aksi_opd_kode
        UNIQUE (kode_rencana_aksi_opd),

    CONSTRAINT fk_rencana_aksi_opd_penetapan
        FOREIGN KEY (penetapan_id)
        REFERENCES penetapan_opd(id)
);

CREATE INDEX idx_rencana_aksi_opd_penetapan_id
    ON rencana_aksi_opd (penetapan_id);

CREATE INDEX idx_rencana_aksi_opd_penetapan_sasaran
    ON rencana_aksi_opd (penetapan_id, kode_sasaran_opd);

CREATE INDEX idx_rencana_aksi_opd_penetapan_pk
    ON rencana_aksi_opd (penetapan_id, kode_pk);

CREATE INDEX idx_rencana_aksi_opd_penetapan_opd_tahun
    ON rencana_aksi_opd (penetapan_id, kode_opd, tahun);
