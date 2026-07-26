CREATE TABLE indikator_renja_individu (
    id BIGSERIAL PRIMARY KEY,

    renja_individu_id BIGINT NOT NULL,
    jenis_indikator VARCHAR(50) NOT NULL,

    kode_indikator_renja VARCHAR(255) NOT NULL,
    indikator TEXT NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_indikator_renja_individu_renja
        FOREIGN KEY (renja_individu_id)
        REFERENCES renja_individu(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_indikator_renja_individu
        UNIQUE (
            renja_individu_id,
            jenis_indikator,
            kode_indikator_renja
        )
);

CREATE INDEX idx_indikator_renja_individu_renja
    ON indikator_renja_individu(renja_individu_id);
