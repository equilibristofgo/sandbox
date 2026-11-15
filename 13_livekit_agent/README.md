# 🎙️ Jarvis LiveKit Agent (Vision & Multi-Avatar)

Agente inteligente para el podcast de **Equilibristas de Golang**, con capacidades de escucha (Whisper), síntesis de voz (Kokoro con OpenVINO en Intel Arc), visión en tiempo real (Qwen2-VL) y **presencia visual (Avatar)**.

---

## 🎨 Sistema de Avatares (3 Modos Soportados)

Puedes seleccionar el modo visual de Jarvis definiendo `AVATAR_MODE` en tu archivo `.env.local`:

### 1. Modo Orbe Sci-Fi Local (`AVATAR_MODE=orb`) — *Por defecto*
* **100% Local (0 MB VRAM extra):** Renderiza en CPU a 25 FPS con Pillow y emite un `rtc.LocalVideoTrack` nativo en la sala de LiveKit.
* **Estados reactivos automáticos:**
  * **Azul/Cian pulsante:** Estado en espera (*Listening/Idle*).
  * **Anillos giratorios:** Procesando o analizando pantalla (*Thinking*).
  * **Ondas dinámicas:** Modulando su radio con la voz de Kokoro (*Speaking*).
* **Uso:** El grabador de Egress y los participantes en la sala lo ven como una cámara web normal.

### 2. Modo Avatar 3D Web / OBS (`AVATAR_MODE=web`)
* **Ubicación:** `sandbox/13_livekit_agent/web/index.html`
* **Tecnología:** Three.js + LiveKit Web SDK + Web Audio API.
* **Uso:** 
  * Ábrelo en el navegador o añádelo como **Browser Source en OBS**.
  * Soporta fondo transparente pasando el parámetro `?transparent=true`.
  * Ejemplo para OBS: `http://localhost:8000/web/index.html?url=ws://localhost:7880&token=...&transparent=true`.

### 3. Modo Anam Cloud Avatar (`AVATAR_MODE=anam`)
* **Avatar Fotorrealista:** Utiliza el plugin oficial de LiveKit (`livekit-plugins-anam`) para lip-sync en la nube.
* **Requisitos:**
  ```bash
  uv pip install livekit-plugins-anam
  ```
  Añadir en `.env.local`:
  ```env
  AVATAR_MODE=anam
  ANAM_API_KEY=tu_api_key
  ANAM_PERSONA_ID=tu_persona_id
  ```

---

## 🚀 Puesta en Marcha

```bash
# 1. Configurar variables de entorno
mcedit .env.local

# 2. Ejecutar agente
uv run python src/agent_vision.py dev
```

### Ejemplo de `.env.local`:
```env
LIVEKIT_URL=ws://192.168.1.4:7880
LIVEKIT_API_KEY=devkey
LIVEKIT_API_SECRET=secret

# Selector de Avatar: "orb" (local), "web" (OBS), o "anam" (nube)
AVATAR_MODE=orb

# Servicios Locales (Intel Arc B580)
WHISPER_URL=http://192.168.1.4:10300/v1
LITELLM_URL=http://192.168.1.4:4001/v1
KOKORO_URL=http://192.168.1.4:8880/v1
OCR_URL=http://localhost:8001/ocr
```
