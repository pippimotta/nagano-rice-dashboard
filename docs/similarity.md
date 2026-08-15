# Similar-Year Ranking: How It Works

This explains the mechanism in `internal/riskmodel/similarity.go`: how the
project turns each year's weather features into a "which past years are most
alike" ranking.

## The question

> Given a target year, which other years had the most similar climate?

Each year is already reduced to a small feature vector by `ComputeYearFeatures`:

| Field | Meaning |
|-------|---------|
| `HeatDaysHeading` | days ≥35°C in the heading window |
| `GrainFillMeanTemp` | mean temperature in the grain-fill window |
| `ColdMeanTemp` | mean temperature in the cold-damage window |
| `GrainFillSunshine` | total sunshine in the grain-fill window |

Ranking means: turn "how different are two years" into a single **distance**
number, then sort by it. Small distance = similar.

## Step 1 — each year is a point in feature space

`vector(f)` flattens a `YearFeature` into 4 numbers, so each year is a point with
4 coordinates. Comparing years becomes measuring distance between points.

## Step 2 — put the features on a fair scale (z-score)

**Problem:** the 4 features use wildly different scales — heat days are 0–15,
sunshine is ~600. A raw distance would be dominated by sunshine; heat days would
barely count.

**Fix:** convert every value to a **z-score** — "how far from the average is
this, measured in the group's own spread":

```
z = (value − mean) / standard_deviation
```

- Subtracting the mean re-centers on 0 (above/below average).
- Dividing by the standard deviation converts that gap into a unit-free number
  ("how many typical spreads"), so all 4 features become comparable.

**Standard deviation** is computed the usual way (see `standardizeParams`):

1. mean of the values
2. for each value, squared deviation `(value − mean)²`
3. average those squared deviations → **variance**
4. square root of variance → **standard deviation**

The mean and standard deviation are computed **once over all years** (that is the
shared "ruler"), then applied to every year — target and candidates alike — by
`standardize`. Using one ruler is what makes the years comparable.

## Step 3 — measure distance (Euclidean)

For two years' z-score vectors, `euclidean` does the straight-line distance:
per dimension take the difference, square it, sum, take the square root — the
same idea as distance on a map, extended to 4 dimensions.

```
distance = sqrt( Σ (z_target[i] − z_candidate[i])² )
```

Squaring makes differences positive and penalizes big gaps more; the square root
brings the result back to a normal scale.

## Step 4 — rank

`RankSimilarYears` computes the distance from the target year to every other
year, then sorts ascending (closest first). Ties break on year number for a
stable, reproducible order.

## Worked example

Two features for clarity (`A` = heat days, `B` = cold mean temp):

| Year | A | B |
|------|---|---|
| 2000 (target) | 6 | 22 |
| 2001 | 6 | 22 |
| 2002 | 2 | 24 |
| 2003 | 12 | 20 |

Ruler (over all 4 years): A mean 6.5, std ≈ 3.57; B mean 22, std ≈ 1.41.

z-scores:

| Year | zA | zB |
|------|------|------|
| 2000 | −0.14 | 0 |
| 2001 | −0.14 | 0 |
| 2002 | −1.26 | +1.41 |
| 2003 | +1.54 | −1.41 |

Distances from 2000:

- 2001: `sqrt(0² + 0²)` = **0.00** (identical)
- 2002: `sqrt((−1.12)² + 1.41²)` ≈ **1.80**
- 2003: `sqrt(1.68² + (−1.41)²)` ≈ **2.20**

Ranking: **2001, 2002, 2003**.

## Edge cases and guards

- **Fewer than 2 years, or target absent** → empty result (nothing to compare).
- **Target excluded from its own ranking** → it would otherwise sit at distance 0
  and be meaningless.
- **A constant dimension (std = 0)** carries no information and would divide by
  zero, so `standardize` contributes 0 for it.

## Not covered here

- The z-score ruler is sensitive to outliers (one extreme year widens a
  dimension's spread and compresses everyone else on that axis).
- Standardizing equalizes *scale*, not *importance*; all features are weighted
  equally. Weighting would be a separate change.

## Map to the code

| Concept | Function |
|---------|----------|
| year → 4 numbers | `vector` |
| build the ruler (mean, std) | `standardizeParams` |
| apply the ruler → z-scores | `standardize` |
| distance between two years | `euclidean` |
| orchestrate + sort | `RankSimilarYears` |
