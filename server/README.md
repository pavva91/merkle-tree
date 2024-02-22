## Bulk upload
```bash
curl --location --request POST 'http://localhost:8080/files' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/path/to/file/f1"' \
--form 'file=@"/path/to/file/f3"' \
--form 'file=@"/path/to/file/medium20MiB"' \
--form 'file=@"/path/to/file/small10MiB"'
```
