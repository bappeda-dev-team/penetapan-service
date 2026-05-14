CREATE TABLE tb_tujuan_penetapan_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_tujuan_opd VARCHAR(255),

    tujuan_opd TEXT NOT NULL,

    periode VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100)

);

CREATE INDEX idx_tujuan_penetapan_opd_kode_tahun
ON tb_tujuan_penetapan_opd(kode_opd, tahun_aktif);



CREATE TABLE tb_indikator_tujuan_penetapan_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    id_tujuan_opd BIGINT NOT NULL,

    kode_indikator VARCHAR(255),

    kode_opd VARCHAR(50) NOT NULL,

    indikator TEXT NOT NULL,

    rumus_perhitungan TEXT,

    sumber_data TEXT,

    definisi_operasional TEXT,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_tujuan_penetapan
        FOREIGN KEY(id_tujuan_opd)
        REFERENCES tb_tujuan_penetapan_opd(id)
        ON DELETE CASCADE

);

CREATE INDEX idx_indikator_tujuan_penetapan_parent
ON tb_indikator_tujuan_penetapan_opd(id_tujuan_opd);

CREATE INDEX idx_indikator_tujuan_penetapan_kode_tahun
ON tb_indikator_tujuan_penetapan_opd(kode_opd, tahun_aktif);



CREATE TABLE tb_target_indikator_tujuan_penetapan_opd (
    id BIGSERIAL PRIMARY KEY,

    indikator_tujuan_id BIGINT NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_indikator_tujuan
        FOREIGN KEY(indikator_tujuan_id)
        REFERENCES tb_indikator_tujuan_penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_tujuan_tahun
        UNIQUE(indikator_tujuan_id, tahun)

);

CREATE INDEX idx_target_indikator_tujuan_parent
ON tb_target_indikator_tujuan_penetapan_opd(indikator_tujuan_id);



CREATE TABLE tb_sasaran_penetapan_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    kode_sasaran_opd VARCHAR(255),

    sasaran_opd TEXT NOT NULL,

    periode VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100)

);

CREATE INDEX idx_sasaran_penetapan_opd_kode_tahun
ON tb_sasaran_penetapan_opd(kode_opd, tahun_aktif);



CREATE TABLE tb_indikator_sasaran_penetapan_opd (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    id_sasaran_opd BIGINT NOT NULL,

    kode_indikator VARCHAR(255),

    kode_opd VARCHAR(50) NOT NULL,

    indikator TEXT NOT NULL,

    rumus_perhitungan TEXT,

    sumber_data TEXT,

    definisi_operasional TEXT,

    tahun_aktif INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_sasaran_penetapan
        FOREIGN KEY(id_sasaran_opd)
        REFERENCES tb_sasaran_penetapan_opd(id)
        ON DELETE CASCADE

);

CREATE INDEX idx_indikator_sasaran_penetapan_parent
ON tb_indikator_sasaran_penetapan_opd(id_sasaran_opd);

CREATE INDEX idx_indikator_sasaran_penetapan_kode_tahun
ON tb_indikator_sasaran_penetapan_opd(kode_opd, tahun_aktif);



CREATE TABLE tb_target_indikator_sasaran_penetapan_opd (
    id BIGSERIAL PRIMARY KEY,

    indikator_sasaran_id BIGINT NOT NULL,

    tahun INTEGER NOT NULL,

    target NUMERIC(20,2) NOT NULL,

    satuan VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_target_indikator_sasaran
        FOREIGN KEY(indikator_sasaran_id)
        REFERENCES tb_indikator_sasaran_penetapan_opd(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_sasaran_tahun
        UNIQUE(indikator_sasaran_id, tahun)
);

CREATE INDEX idx_target_indikator_sasaran_parent
ON tb_target_indikator_sasaran_penetapan_opd(indikator_sasaran_id);
