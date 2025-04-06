# Web Crawler CLI

A fast, concurrent, domain-aware website crawler written in Go - outputs a clean sitemap of all internal links.

![Go](https://img.shields.io/badge/built%20with-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![CLI Tool](https://img.shields.io/badge/type-CLI-2ea44f?style=for-the-badge)

---

### Usage

```bash
go run main.go [flags] <url>
```

### Flags

```text
--depth        Maximum crawl depth (default: 3)
--domain-only  Only crawl links on the same domain (default: true)
```

### Examples

```bash
# Crawl a website up to 3 levels deep
go run main.go https://example.com

# Crawl up to 5 levels, allowing external domains
go run main.go --depth=5 --domain-only=false https://example.com
```

---

### Output

After crawling, you'll get a `sitemap.txt` file with a sorted list of all internal URLs:

```text
https://example.com
https://example.com/about
https://example.com/blog/post-1
...
```
