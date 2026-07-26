CREATE TABLE target_renja_individu (
    id BIGSERIAL PRIMARY KEY,

    indikator_renja_individu_id BIGINT NOT NULL,
    jenis_target VARCHAR(50) NOT NULL,

    kode_target_renja VARCHAR(255) NOT NULL,

    target NUMERIC(20,2) NOT NULL,
    satuan VARCHAR(100) NOT NULL,
    tahun INTEGER NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_target_renja_individu_indikator
        FOREIGN KEY (indikator_renja_individu_id)
        REFERENCES indikator_renja_individu(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_target_renja_individu
        UNIQUE (
            indikator_renja_individu_id,
            kode_target_renja
        )
);

CREATE INDEX idx_target_renja_individu_indikator
    ON target_renja_individu(indikator_renja_individu_id);

CREATE INDEX idx_target_renja_individu_tahun
    ON target_renja_individu(tahun);
