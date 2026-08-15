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
                 ├─▶ ingest ─▶ BigQuery (raw) ─▶ features ─▶ year_features ─▶ JSON ─▶ dashboard
MAFF e-Stat API ─┘
```

- **Ingest (Go):** load daily weather and annual rice statistics into BigQuery.
- **Features (Go):** for each year, compute growth-stage risk features and rank
  past years by climate similarity (z-score standardized Euclidean distance).
- **Storage (BigQuery):** two raw tables (`weather_daily`, `rice_yield_annual`)
  and one derived table (`year_features`), exported to a static JSON.
- **Dashboard:** a static one-page view. Pick a target year; see the most
  climate-similar past years and their crop-situation index and yield.

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
| **MAFF** (農林水産省) | Rice crop-situation index, yield, harvest, area by prefecture | e-Stat API |
| **WAGRI** (農業データ連携基盤) | Pest/disease risk and growth-stage forecasts | API (planned, v1) |

v0 uses a single JMA station (Nagano) and roughly the last 30 years.

## Tech Stack

- **Pipeline:** Go
- **Infrastructure:** GCP — BigQuery, Cloud Run, Cloud Storage — provisioned with
  Terraform
- **Frontend:** static HTML + a lightweight charting library
- **CI:** GitHub Actions (lint, build, test)

## Roadmap

**v0**
- Ingest JMA (Nagano) weather and MAFF rice statistics
- Growth-stage risk features + similar-year ranking
- Static one-page dashboard
- Terraform for GCP resources, CI

**Later (v1+)**
- WAGRI integration (pest/disease risk, growth-stage forecasts)
- Multiple stations / elevation contrast (e.g. Iida vs. Karuizawa)
- Dynamic growth-stage windows based on accumulated temperature
- Deeper history and a live current-year feed

## Data & Licensing

- JMA data is used under the JMA Public Data License (公共データ利用規約 v1.0);
  attribution: 出典 気象庁.
- MAFF / e-Stat data is used under the Government of Japan Standard Terms of Use
  (政府標準利用規約 v2.0, CC BY 4.0 compatible).

