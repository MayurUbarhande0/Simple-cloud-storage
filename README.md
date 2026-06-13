```markdown
# Simple Cloud Storage

This is an Simple File Storage system meant to be run locally 
It uses sql database for store file details 
it have jwt auth
Some part of the system are meant to build by you , like login and signup , but i will try to add that


---

## 🏗️ Architecture Layout

* **HTTP Gateway (Port `:8080`):** Acts as the reverse proxy entryway. Handles client authentication (JWT verification), it have routes of like /upload /download etc


```text
Client ──[ HTTP Multipart / JWT ]──> Gateway (Gzip)


```

---

## 🛠️ Installation & Setup

### 1. Configure the Environment

Clone the project repository and copy the environment template:

```bash
cp .env.example .env

```

Fill out the variables in `.env` with your secure cryptographic tokens and network endpoints.

### 2. Boot the Storage Node Engine

Start the storage engine listener from the root directory:

```bash
env $(cat .env | grep -v '^#' | xargs) go run main.go

```

### 3. Boot the HTTP Gateway Server

In a separate terminal window, launch the entry endpoint:

```bash
env $(cat .env | grep -v '^#' | xargs) go run gateway/main.go

```

---

## 🧪 How to Test End-to-End

To verify your complete, authenticated distributed storage pipeline, open a third terminal window alongside your running servers and fire the test client simulation:

```bash
go run ./tools/gen_token.go -user alice
```

Then use the browser frontend in `frontend/` or `curl` to POST to `http://127.0.0.1:8080/upload`.

### 🔍 Verification Checklist

1. **Client Console Output:** The test script should log an explicit HTTP `200 OK` transaction along with a tracking file token:
```text
Gateway Status Response Code: 200
Gateway Core JSON Body payload: {"file_id":"d8858bgsaf1p4k...", "status":"success"}

```


2. **Metadata Synchronization Verification:** Inspect the newly written `state.json` tracker file. It will render the mapped structure layout holding your verified runtime strings:
```json
{
 "User": {
  "user_test_99": {
   "user_id": "user_test_99",
   "file": [
    {
     "file_id": "d8858bgsaf1p4k7q7c2g",
     "filename": "verified_notes.txt",
     "storage_path": "/cloud/user_test_99/documents/production/test/d8858bgsaf1p4k7q7c2g",
     "size": 65
    }
   ]
  }
 }
}

```


3. **Data Compaction Realization:** Confirm that the target directory on your machine now contains the gzipped storage block written cleanly through the raw TCP stream handler.

---

## ⚡ WILL DO SOON

* **Dynamic Custom Binary Framing Protocol:** Uses 1-byte dynamic size headers to safely transmit structural payloads without buffer hanging, size constraints, or network deadlocks.
* **Zero-Memory Allocation Stream Copying:** Utilizes `io.Copy` to stream incoming multipart form uploads straight across network TCP sockets without buffering large chunks into system RAM.
* **Encrypted TCP Handshakes:** Gateway authentication tokens are encrypted using an AES-GCM helper and verified via SHA-256 signatures on the storage node before data channels unlock.
* **Thread-Safe Local State Engine:** Serializes asset metadata blocks dynamically to a local `state.json` registry file under a transactional `sync.RWMutex` concurrent lock pattern.
