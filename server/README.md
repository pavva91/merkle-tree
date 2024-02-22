## Bulk upload

```bash
curl --location --request POST 'http://localhost:8080/files' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/path/to/file/f1"' \
--form 'file=@"/path/to/file/f2"' \
--form 'file=@"/path/to/file/f3"'
```
## Download a file

```bash
curl -v -X GET 'http://localhost:8080/files/f1' -o tmp.out
```
