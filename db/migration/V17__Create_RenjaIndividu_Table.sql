CREATE TABLE renja_individu (
    id BIGSERIAL PRIMARY KEY,

    penetapan_individu_id BIGINT NOT NULL,

    pegawai_id VARCHAR(100) NOT NULL,

    kode_opd VARCHAR(50) NOT NULL,

    tahun_aktif INTEGER NOT NULL,

    kode_pk VARCHAR(255) NOT NULL,

    kode_program VARCHAR(50) NOT NULL,
    nama_program TEXT NOT NULL,
    pagu_program BIGINT NOT NULL,

    kode_kegiatan VARCHAR(50) NOT NULL,
    nama_kegiatan TEXT NOT NULL,
    pagu_kegiatan BIGINT NOT NULL,

    kode_subkegiatan VARCHAR(50) NOT NULL,
    nama_subkegiatan TEXT NOT NULL,
    pagu_subkegiatan BIGINT NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_renja_individu_penetapan
        FOREIGN KEY (penetapan_individu_id)
        REFERENCES penetapan_individu(id)
        ON DELETE CASCADE
);
