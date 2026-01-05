# Profile Ranking Algorithm

This document describes the algorithm used to rank and order profiles in the browse/discovery features.

## Overview

The ranking system uses a **hybrid scoring approach**:
- **Static scores**: Pre-calculated and stored in the database, updated every 15 minutes by a background job
- **Dynamic scores**: Calculated at query time based on real-time factors

```
Final Score = Static Score + Dynamic Modifiers
```

Profiles are ordered by `final_score DESC` (highest score first).

---

## Static Score Components

These factors are calculated by the `RankingService` and stored in `profiles.profile_score`.

| Factor | Points | Calculation |
|--------|--------|-------------|
| Email Verified | +10 | Boolean: user has verified their email |
| Person Verified | +25 | Boolean: profile has been identity-verified |
| Profile Completeness | 0-20 | See breakdown below |
| Popularity | 0-15 | `min(15, likes_received × 0.5)` |
| Response Rate | 0-15 | `response_rate_percent × 0.15` |
| New User Boost | 0-10 | `10 × max(0, 1 - days_since_creation/7)` |

### Profile Completeness Breakdown (0-20 points)

| Sub-factor | Points | Calculation |
|------------|--------|-------------|
| Photos | 0-10 | `min(10, photo_count/5 × 10)` |
| Bio Length | 0-5 | `min(5, bio_length/300 × 5)` |
| Salary Range | 0-5 | +5 if salary_range is provided |

### Response Rate Calculation

Response rate measures how often a user replies to conversations:

```
response_rate = (conversations_where_user_replied / conversations_where_user_received_messages) × 100
```

- If user has no conversations, defaults to 50%
- Rewards users who actively engage with matches

### New User Boost

New profiles get a temporary visibility boost that decays linearly over 7 days:

```
Day 0: +10 points
Day 1: +8.57 points
Day 3: +5.71 points
Day 7: +0 points (boost expires)
```

---

## Dynamic Score Components

These factors are calculated at query time in the `ListProfiles` SQL query.

| Factor | Points | Condition |
|--------|--------|-----------|
| Online Now | +15 | User is currently online (connected to WebSocket) |
| Recently Active | +5 | Last seen within the past hour (but not online) |
| Mutual Interest | +30 | They have already liked the viewing user |
| Distance Penalty | 0 to -20 | Applied for profiles >50km away |

### Distance Penalty Formula

```sql
CASE 
    WHEN distance_km < 50 THEN 0
    ELSE -min(20, distance_km × 0.1)
END
```

Examples:
- 30 km away: 0 penalty
- 50 km away: 0 penalty
- 100 km away: -10 penalty
- 200+ km away: -20 penalty (capped)

---

## Score Ranges

| Profile Type | Approximate Score Range |
|--------------|------------------------|
| New, unverified, incomplete | 10-25 |
| Complete, email verified | 30-50 |
| Verified, popular, active | 60-80 |
| Verified + they liked you + online | 90-120+ |

---

## Background Job

The `RankingService` runs as a background goroutine:

- **Interval**: Every 15 minutes
- **Location**: `backend/internal/services/ranking.go`
- **Started in**: `backend/cmd/server/main.go`

On each run:
1. Fetches all complete profiles
2. Calculates static score for each
3. Updates `profiles.profile_score` column

---

## Database Schema

```sql
-- Migration: 012_add_profile_score

ALTER TABLE profiles ADD COLUMN profile_score FLOAT NOT NULL DEFAULT 0;
CREATE INDEX idx_profiles_score ON profiles(profile_score DESC);
```

---

## Files

| File | Purpose |
|------|---------|
| `backend/internal/services/ranking.go` | Static score calculation + background job |
| `backend/internal/repository/profile.go` | Dynamic scoring in `ListProfiles` query |
| `backend/internal/models/profile.go` | `ProfileScore` field definition |
| `backend/migrations/012_add_profile_score.up.sql` | Database migration |

---

## Tuning

To adjust the algorithm, modify the weights in:

1. **Static factors**: `ranking.go` → `CalculateStaticScore()`
2. **Dynamic factors**: `profile.go` → `ListProfiles()` → `finalScoreCalc`

### Suggested A/B Tests

- Increase/decrease mutual interest bonus (+30)
- Adjust new user boost duration (7 days → 14 days)
- Modify distance penalty threshold (50km → 100km)
- Weight response rate higher for premium features



