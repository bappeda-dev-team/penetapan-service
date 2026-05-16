CREATE TABLE tujuan_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_tujuan_opd VARCHAR(255) NOT NULL,

    tujuan_opd TEXT NOT NULL,

    periode VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100)

);

CREATE INDEX idx_tujuan_opd_kode_opd_tahun_aktif
ON tujuan_opd(kode_opd, tahun_aktif);



CREATE TABLE indikator_tujuan_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    tujuan_opd_id BIGINT NOT NULL,

    kode_indikator VARCHAR(255) NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    indikator TEXT NOT NULL,

    rumus_perhitungan TEXT,

    sumber_data TEXT,

    definisi_operasional TEXT,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_tujuan
        FOREIGN KEY(tujuan_opd_id)
        REFERENCES tujuan_opd(id)
        ON DELETE CASCADE

);

CREATE INDEX idx_indikator_tujuan_opd_parent
ON indikator_tujuan_opd(tujuan_opd_id);

CREATE INDEX idx_indikator_tujuan_opd_kode_opd_tahun_aktif
ON indikator_tujuan_opd(kode_opd, tahun_aktif);


CREATE TABLE target_indikator_tujuan_opd (
    id BIGSERIAL PRIMARY KEY,

    indikator_tujuan_id BIGINT NOT NULL,

    kode_target VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_indikator_tujuan_opd
        FOREIGN KEY(indikator_tujuan_id)
        REFERENCES indikator_tujuan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_tujuan_tahun
        UNIQUE(indikator_tujuan_id, tahun)

);

CREATE INDEX idx_target_indikator_tujuan_parent
ON target_indikator_tujuan_opd(indikator_tujuan_id);



CREATE TABLE sasaran_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_sasaran_opd VARCHAR(255) NOT NULL,

    sasaran_opd TEXT NOT NULL,

    periode VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100)

);

CREATE INDEX idx_sasaran_opd_kode_opd_tahun_aktif
ON sasaran_opd(kode_opd, tahun_aktif);



CREATE TABLE indikator_sasaran_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    sasaran_opd_id BIGINT NOT NULL,

    kode_indikator VARCHAR(255) NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    indikator TEXT NOT NULL,

    rumus_perhitungan TEXT,

    sumber_data TEXT,

    definisi_operasional TEXT,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_indikator_sasaran
        FOREIGN KEY(sasaran_opd_id)
        REFERENCES sasaran_opd(id)
        ON DELETE CASCADE

);

CREATE INDEX idx_indikator_sasaran_opd_parent
ON indikator_sasaran_opd(sasaran_opd_id);

CREATE INDEX idx_indikator_sasaran_opd_kode_opd_tahun_aktif
ON indikator_sasaran_opd(kode_opd, tahun_aktif);


CREATE TABLE target_indikator_sasaran_opd (
    id BIGSERIAL PRIMARY KEY,

    indikator_sasaran_id BIGINT NOT NULL,

    kode_target VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_indikator_sasaran_opd
        FOREIGN KEY(indikator_sasaran_id)
        REFERENCES indikator_sasaran_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_sasaran_tahun
        UNIQUE(indikator_sasaran_id, tahun)
);

CREATE INDEX idx_target_indikator_sasaran_opd_parent
ON target_indikator_sasaran_opd(indikator_sasaran_id);
