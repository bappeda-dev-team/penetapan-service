-- =========================================================
-- SNAPSHOT PAGU URUSAN
-- =========================================================

CREATE TABLE pagu_renja_urusan (
    id BIGSERIAL PRIMARY KEY,

    urusan_id BIGINT NOT NULL,

    kode_pagu VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    pagu BIGINT NOT NULL,

    jenis_pagu VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_pagu_urusan
        FOREIGN KEY(urusan_id)
        REFERENCES renja_urusan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pagu_urusan
        UNIQUE(
            urusan_id,
            kode_pagu
        )
);

CREATE INDEX idx_pagu_urusan_source
ON pagu_renja_urusan(
    kode_pagu
);

CREATE INDEX idx_pagu_urusan_tahun
ON pagu_renja_urusan(
    tahun
);


-- =========================================================
-- SNAPSHOT PAGU BIDANG URUSAN
-- =========================================================

CREATE TABLE pagu_renja_bidang_urusan (
    id BIGSERIAL PRIMARY KEY,

    bidang_urusan_id BIGINT NOT NULL,

    kode_pagu VARCHAR(255) NOT NULL,

    tahun INTEGER NOT NULL,

    pagu BIGINT NOT NULL,

    jenis_pagu VARCHAR(50) NOT NULL,

    created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by VARCHAR(100),

    CONSTRAINT fk_pagu_bidang_urusan
        FOREIGN KEY(bidang_urusan_id)
        REFERENCES renja_bidang_urusan(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pagu_bidang_urusan
        UNIQUE(
            bidang_urusan_id,
            kode_pagu
        )
);

CREATE INDEX idx_pagu_bidang_urusan_source
ON pagu_renja_bidang_urusan(
    kode_pagu
);

CREATE INDEX idx_pagu_bidang_urusan_tahun
ON pagu_renja_bidang_urusan(
    tahun
);
