this file is somewhere in `…/metrics/stories/publisher_uniqs.md`

```raw
type = "count(country, version)"
```

# queries

these are queries for unique users who published stories in a day

* per all countries and vesions
```sql
SELECT
    $Date AS date,
    uniqExact(userID) AS count,
    'WW' AS country,
    'all' AS version
FROM stories
WHERE date = $Date
```

* per version within all countries
```sql
SELECT
    $Date AS date,
    uniqExact(userID) AS count,
    'WW' AS country,
    version
FROM stories
WHERE date = $Date
GROUP BY version
```

* per country no matter the version
```sql
SELECT
    $Date AS date,
    uniqExact(userID) AS count,
    country,
    'all' AS version 
FROM stories
WHERE date = $Date
GROUP BY country
```

* per country and version
```sql
SELECT
    $Date AS date,
    uniqExact(userID) AS count,
    country,
    version
FROM stories
WHERE date = $Date
GROUP BY country, version
```
