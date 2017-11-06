# Typical usecase.

Statistical systems based on storage with SQL interface, such as [Clickhouse](https://clickhouse.yandex), 
[Vertica](https://vertica.com), etc. Typically, you have a lot of queries computing metric values for them, this means
a lot of SQL. At my previous job (2015-2017) we have a tree of queries like `<period>.<product name>.<feature>.<metric name>`.
For instance `daily.chat.story.uniq_creators` which has a queury to compute amount of uniq users who created a story through
a given date. We ended by mapping this name into a file system path:
`${METRIC_DIR}/daily/chat/story/uniq_creators.md`
where the file itself looked like

![example](usecase_example.png)

the query launcher will consume this file in a way similar to:

```go
…
input, err := ioutil.ReadFile("daily/chat/story/uniq_creators.md")
if err != nil {
	return err
}
var metric struct {
	Queries []mad.Source `mad:"queries,source=sql"`
}
if err := mad.Unmarshal(input, &metric); err != nil {
	return err
}
…
``` 
of course we could use other structured formats suitable for humans, such as YAML, but:

1. Its toolset is nowhere near as developed as Markdown support. Specifically, there's no syntax aware highlighting for 
code blocks and this is really matters in case such as ours.
2. YAML unmarshaller is a bit too allowing. We wanted to get an error when we declare a non-pointer field 
    ```go
    type Metric struct {
        …
        Field mad.Source `mad:"field,source=json"`
        …
    }
    ```
    and the Markdown input doesn't have `# field` header.
 


