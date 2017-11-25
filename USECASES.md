	t.Log(dest.Queries)
# Typical usecase.

Statistical systems based on storage with SQL interface, such as [Clickhouse](https://clickhouse.yandex), 
[Vertica](https://vertica.com), etc. Typically, you have a lot of queries computing values for them, this means
a lot of SQL. At my previous job (2015-2017) we have a tree of metrics like `<period>.<product name>.<feature>.<metric name>` (each metric can have several queries).
For instance `daily.chat.story.uniq_creators` which has a queury to compute amount of uniq users who created a story through
a given date. We ended by mapping this name into a file system path:
`${METRIC_DIR}/daily/chat/story/uniq_creators.md`
where the file itself was a piece of YAML which had a tendency to become unreadable at times. 

![example](usecase_example.png)

the query launcher will consume this file in a way similar to:

```go
…
input, err := ioutil.ReadFile("daily/story/uniq_creators.md")
if err != nil {
	return err
}
var metric struct {
	Type string `mad:"type"`
	Queries []mad.Source `mad:"queries,syntax=sql"`
}
if err := mad.Unmarshal(input, &metric); err != nil {
	return err
}
…
``` 



