Simple Frontend for the Gateway

This simple frontend lets you upload and download files to the gateway server running on the same origin (http://localhost:8080).

How to use

1. Start the gateway server (from repo root):

```bash
export DB_URL=$(pwd)/identifier.sqlite
export SECRET_KEY=test-key
cd gateway
go run main.go
```

2. Serve the frontend files. The easiest is to run a static server from the `frontend` directory. For example:

```bash
cd frontend
# python3's simple HTTP server (serves on :8000)
python3 -m http.server 8000
```

3. Open http://localhost:8000 in your browser.

Notes
- The frontend expects the gateway to be reachable at `/upload` and `/download` on the same origin; if your gateway runs on a different host/port, update `app.js` fetch URLs accordingly.
- The upload sends a multipart/form-data POST with fields `path`, `filename` and `file` and the optional `Authorization: Bearer <token>` header.
- The download issues a POST to `/download?file_id=<id>` with optional `Authorization` header and triggers a download of the returned gzip file.

