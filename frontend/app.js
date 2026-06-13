// Configure the API base URL for the gateway. When serving the frontend
// with a static server (python -m http.server) you MUST use the gateway
// host:port here (e.g. http://127.0.0.1:8080) because the static server
// doesn't proxy POST requests to the API.
const API_BASE = 'http://127.0.0.1:8080';

document.getElementById('upload-btn').addEventListener('click', async () => {
  const token = document.getElementById('upload-token').value.trim();
  const path = document.getElementById('upload-path').value.trim();
  const filename = document.getElementById('upload-filename').value.trim();
  const fileElem = document.getElementById('upload-file');
  const result = document.getElementById('upload-result');
  result.innerText = '';

  if (!fileElem.files || fileElem.files.length === 0) {
    result.innerText = 'Select a file first';
    return;
  }
  const file = fileElem.files[0];
  const form = new FormData();
  form.append('path', path);
  form.append('filename', filename);
  form.append('file', file, file.name);

  try {
    const res = await fetch(API_BASE + '/upload', {
      method: 'POST',
      headers: token ? { 'Authorization': 'Bearer ' + token } : {},
      body: form
    });
    const text = await res.text();
    result.innerText = `HTTP ${res.status}\n${text}`;
  } catch (err) {
    result.innerText = 'Error: ' + err;
  }
});


document.getElementById('download-btn').addEventListener('click', async () => {
  const token = document.getElementById('download-token').value.trim();
  const fileId = document.getElementById('download-file-id').value.trim();
  const result = document.getElementById('download-result');
  result.innerText = '';
  if (!fileId) { result.innerText = 'Enter file id'; return; }

  try {
    const res = await fetch(API_BASE + '/download?file_id=' + encodeURIComponent(fileId), {
      method: 'POST',
      headers: token ? { 'Authorization': 'Bearer ' + token } : {}
    });
    if (!res.ok) {
      const txt = await res.text();
      result.innerText = `HTTP ${res.status}\n${txt}`;
      return;
    }
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = fileId + '.gz';
    a.click();
    URL.revokeObjectURL(url);
    result.innerText = 'Download started';
  } catch (err) {
    result.innerText = 'Error: ' + err;
  }
});

