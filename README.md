> [!NOTE]  
> Here be dragons
>
> Needed to build this for something, and it served its purpose, but as of now, most of it remains unfinished. It was built in a rush, so the code isn't pretty and expect bugs. I do not recommend people use this. If you decide to use this, do not expect support or documentation.

# bigcord

Archive large discord servers via a bot

To export scraped data as parquet:

```shell
mkdir -p data
prefix="data/export-$(date +%s)"
for table in "default.messages" "default.channels"; do
  filename="$prefix-$table.parquet"
  echo "Exporting $table to $filename"
  compose exec clickhouse clickhouse-client -q "SELECT * FROM $table FINAL" --format Parquet > "$filename"
done

compose cp scraping:/data/media data/media
```