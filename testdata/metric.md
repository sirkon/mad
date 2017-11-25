this file is somewhere in `…/metrics/stories/publisher_uniqs.md`

```raw
type = "count(country, version)"
```

# queries

these are queries for unique users who published stories in a day

* for all countries and vesions
```sql
SELECT
    $Date AS date,
    uniqExact(userID) AS count,
    'WW' AS country,
    'all' AS version
FROM stories
WHERE date = $Date
```

* for version in all countries
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

* for country no matter the version
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

* for country and version
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
