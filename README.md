# zeno-export

GO-WASM Export Engine for Zeno CY
Local Excel generation from normalized transactions.
Runs fully in the browser. Zero-retention. No server-side processing.

---

## 1. Mission

This module is responsible for converting normalized statement data
(from `zeno-engine` → `core.wasm`) into a downloadable Excel file.

- 100% local browser execution
- No data sent to backend
- Privacy-by-design
- Independent WASM module
- Loaded lazily only when needed

---

## 2. Architecture Overview
