# Nagano Rice Growth Risk Dashboard

A small data pipeline and dashboard that cross-references historical climate with
rice yields in Nagano, Japan, to answer a single question:

> **Given this year's climate pattern, which past years had similar conditions,
> and what were their yields?**

It is a climate-to-risk lookup prototype built on Japanese public agricultural
and meteorological open data.

## Motivation

I live in Nagano, one of Japan's important rice-producing regions. I got curious
about rice cultivation after playing *Sakuna: Of Rice and Ruin*
(天穂のサクナヒメ) — the game is unexpectedly detailed about how weather, water,
and timing shape a harvest. That curiosity turned into wanting to look at the
real thing: the open weather and crop data published for the region I actually
live in, and whether the patterns the game hints at show up in the numbers.

## Core Concept

Each year's weather is condensed into a small set of **growth-stage risk
features** — numbers that capture how stressful that year's climate was during
the rice plant's sensitive stages. Years are then compared by climate similarity,
and each similar year is shown alongside its actual crop outcome.

The output is presented as a analog lookup (instead of production prediction): it finds the past years most like a chosen year and lets you read across to how those years turned out — a data-driven way to reason about risk, closer to a decision-support prototype than a forecasting model.

## How It Works

```
JMA weather CSV ─┐
                 ├─▶ features (Go) ─▶ year_features.json ─▶ dashboard (Vercel)
MAFF yield JSON ─┘
```

- **Features (Go):** parse the committed raw files (JMA CSV, MAFF JSON), compute
  each year's growth-stage risk features, rank past years by climate similarity
  (z-score standardized Euclidean distance), and export a single static JSON with
  every year's ranking precomputed.
- **Serving:** a static one-page dashboard, deployed on Vercel, reads that JSON.
  Rankings are precomputed, so there is no query at runtime — pick a target year
  and read across to the most climate-similar past years and their crop-situation
  index and yield.
- **Storage:** the exported JSON is the pipeline's product. A private,
  Terraform-managed GCS bucket holds it as the derived-artifact store, out of the
  request path — the browser only talks to the static host.

## Rice Cultivation Risks (background)

The features model the main weather-driven risks for rice in Nagano:

- **High-temperature damage (高温障害):** temperatures ≥35°C during
  heading/flowering cause pollination failure; sustained heat during
  grain-filling produces chalky grains (白未熟粒) and downgrades quality.
- **Cold damage (冷害):** abnormally low summer temperatures prevent normal ear
  development — a traditional risk for highland areas like Nagano.
- **Pest & disease:** rice blast (いもち病) correlates strongly with temperature
  and humidity; the brown planthopper (トビイロウンカ) is migrating further north
  with warming. (Not modeled in v0.)
- **Water management:** irrigation timing affects yield. (Not modeled in v0.)

These weather-driven risks share a common trait: they can be reasoned about from
meteorological data (temperature, humidity, precipitation, sunshine) plus the
current growth stage — which is what this project builds on.

## Data Sources

| Source | What it provides | Access |
|--------|------------------|--------|
| **JMA** (気象庁) | Daily temperature, precipitation, sunshine per station (Nagano) | CSV download |
| **MAFF** (農林水産省) | Rice crop-situation index, yield, harvest, area by prefecture | e-Stat API (fetched once, committed) |
| **WAGRI** (農業データ連携基盤) | Pest/disease risk and growth-stage forecasts | API (planned, v1) |

v0 uses a single JMA station (Nagano): weather from 1996 and rice statistics
through 2020 (25 overlapping years). Both raw datasets are committed as versioned
fixtures — JMA has no public API, and the MAFF long-term table is a static
historical series fetched once from the e-Stat API — so the pipeline is fully
reproducible with no runtime API call or secret.

## Tech Stack

- **Pipeline:** Go
- **Infrastructure:** GCP Cloud Storage (private artifact store), provisioned with
  Terraform
- **Frontend:** static HTML with inline-SVG charts, deployed on Vercel
- **CI:** GitHub Actions (lint, build, test)

## Roadmap

**v0**
- Parse JMA (Nagano) weather and MAFF rice statistics
- Growth-stage risk features + similar-year ranking
- Static one-page dashboard on Vercel
- Terraform for GCP artifact storage, CI

**Later (v1+)**
- WAGRI integration (pest/disease risk, growth-stage forecasts)
- Multiple stations / elevation contrast (e.g. Iida vs. Karuizawa)
- Dynamic growth-stage windows based on accumulated temperature
- Deeper history and a live current-year feed

**Deliberate non-goals**
- No data warehouse. The dataset is yearly-aggregated for a single region (a few
  KB), far below the scale where a warehouse like BigQuery earns its place, and
  the transform already lives in tested Go — a warehouse would only duplicate
  that logic in SQL for no functional gain.

## Data & Licensing

- JMA data is used under the JMA Public Data License (公共データ利用規約 v1.0);
  attribution: 出典 気象庁.
- MAFF / e-Stat data is used under the Government of Japan Standard Terms of Use
  (政府標準利用規約 v2.0, CC BY 4.0 compatible).

