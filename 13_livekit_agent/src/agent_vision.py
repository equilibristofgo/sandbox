import logging
import io
import aiohttp
import re
import json
import uuid
import os
import math
import time
import asyncio
from dotenv import load_dotenv
from PIL import Image, ImageDraw
from livekit import rtc
from livekit.agents import (
    Agent,
    AgentServer,
    AgentSession,
    JobContext,
    JobProcess,
    cli,
    room_io,
    function_tool,
    RunContext,
    llm,
    ModelSettings,
)
from livekit.plugins import openai, silero

logger = logging.getLogger("agent")
class JarvisOrbVisualizer:
    """Generador visual de un orbe holográfico reactivo publicado como rtc.LocalVideoTrack en LiveKit.
    100% local, corre en CPU con Pillow a 25 FPS sin consumir VRAM de la GPU."""
    def __init__(self, width: int = 512, height: int = 512, fps: int = 25):
        self.width = width
        self.height = height
        self.fps = fps
        self.state = "idle"  # "idle", "thinking", "speaking"
        self.running = False
        self.source = rtc.VideoSource(width=width, height=height)
        self.track = rtc.LocalVideoTrack.create_video_track("jarvis-orb", self.source)
        self._task = None

    def set_state(self, state: str):
        if state in ("idle", "thinking", "speaking") and state != self.state:
            self.state = state
            logger.info(f"[JarvisOrbVisualizer] Estado visual: {state.upper()}")

    async def start(self):
        self.running = True
        self._task = asyncio.create_task(self._render_loop())

    async def stop(self):
        self.running = False
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass

    async def _render_loop(self):
        frame_duration = 1.0 / self.fps
        t = 0.0
        cx, cy = self.width // 2, self.height // 2

        while self.running:
            start_time = time.monotonic()
            t += frame_duration

            # Fondo espacial oscuro / cian muy profundo
            img = Image.new("RGBA", (self.width, self.height), (10, 14, 23, 255))
            draw = ImageDraw.Draw(img)

            if self.state == "idle":
                # Estado de escucha: respiración suave azul/cian
                pulse = math.sin(t * 2.2)
                core_r = 75 + pulse * 7

                draw.ellipse([cx - (core_r + 40), cy - (core_r + 40), cx + (core_r + 40), cy + (core_r + 40)], outline=(0, 180, 220, 60), width=2)
                draw.ellipse([cx - (core_r + 18), cy - (core_r + 18), cx + (core_r + 18), cy + (core_r + 18)], outline=(0, 215, 255, 120), width=3)
                draw.ellipse([cx - core_r, cy - core_r, cx + core_r, cy + core_r], fill=(0, 160, 235, 230), outline=(140, 240, 255, 255), width=4)
                draw.ellipse([cx - 20, cy - 20, cx + 20, cy + 20], fill=(225, 248, 255, 255))

            elif self.state == "thinking":
                # Estado pensando: anillos orbitales giratorios de cálculo rápido
                pulse = math.sin(t * 6.5)
                core_r = 65 + pulse * 5

                draw.ellipse([cx - core_r, cy - core_r, cx + core_r, cy + core_r], fill=(0, 210, 230, 230), outline=(255, 255, 255, 255), width=3)

                for ring_idx, (r_x, r_y, speed) in enumerate([(135, 48, 3.8), (115, 38, -4.5), (155, 60, 2.3)]):
                    angle = t * speed
                    for dot in range(12):
                        th = angle + dot * (math.pi / 6)
                        dx = cx + r_x * math.cos(th)
                        dy = cy + r_y * math.sin(th)
                        draw.ellipse([dx - 4, dy - 4, dx + 4, dy + 4], fill=(0, 240, 255, 210))

                draw.ellipse([cx - 15, cy - 15, cx + 15, cy + 15], fill=(255, 255, 255, 255))

            elif self.state == "speaking":
                # Estado hablando: ondas de frecuencia dinámica y pulso sonoro
                audio_wave = abs(math.sin(t * 8.5)) * 26 + math.sin(t * 16.0) * 12
                core_r = 80 + audio_wave

                for i in range(1, 4):
                    ring_r = core_r + i * 28 + (math.sin(t * 11.0 + i) * 6)
                    alpha = max(20, 180 - i * 50)
                    draw.ellipse([cx - ring_r, cy - ring_r, cx + ring_r, cy + ring_r], outline=(0, 230, 255, alpha), width=3)

                draw.ellipse([cx - core_r, cy - core_r, cx + core_r, cy + core_r], fill=(0, 155, 255, 240), outline=(210, 250, 255, 255), width=5)
                draw.ellipse([cx - 28, cy - 28, cx + 28, cy + 28], fill=(245, 255, 255, 255))

            raw_bytes = img.tobytes()
            video_frame = rtc.VideoFrame(self.width, self.height, rtc.VideoBufferType.RGBA, raw_bytes)
            self.source.capture_frame(video_frame)

            elapsed = time.monotonic() - start_time
            sleep_time = max(0.001, frame_duration - elapsed)
            await asyncio.sleep(sleep_time)


class VisionManager:
    def __init__(self, ocr_url: str = "http://localhost:8001/ocr"):
        self.ocr_url = ocr_url

    async def get_latest_screenshot(self, room: rtc.Room) -> io.BytesIO:
        """Busca el track de Screen Share y captura el último frame."""
        logger.info(f"Buscando tracks de SCREEN_SHARE en la sala con {len(room.remote_participants)} participantes.")
        for participant in room.remote_participants.values():
            for track_pub in participant.track_publications.values():
                if track_pub.source == rtc.TrackSource.SOURCE_SCREENSHARE:
                    logger.info(f"¡He contrado un track de SCREEN_SHARE del participante {participant.identity}!")
                    if not track_pub.subscribed:
                        logger.info("El track no estaba suscrito, suscribiendo...")
                        await track_pub.set_subscribed(True)
                    
                    track = track_pub.track
                    if isinstance(track, rtc.VideoTrack):
                        logger.info("Capturando frame del VideoTrack...")
                        # Capturamos el primer frame disponible
                        video_stream = rtc.VideoStream(track)
                        async for event in video_stream:
                            frame = event.frame
                            logger.info(f"Frame capturado ({frame.width}x{frame.height}). Procesando...")
                            rgba_frame = frame.convert(rtc.VideoBufferType.RGBA)
                            image = Image.frombytes("RGBA", (rgba_frame.width, rgba_frame.height), rgba_frame.data)
                            rgb_image = image.convert("RGB")
                            
                            buf = io.BytesIO()
                            rgb_image.save(buf, format="PNG")
                            buf.seek(0)
                            await video_stream.aclose()
                            logger.info("Frame capturado y convertido con éxito.")
                            return buf
        logger.warning("No se encontró ningún track de SCREEN_SHARE activo.")
        return None

class Assistant(Agent):
    def __init__(self, room: rtc.Room, visualizer: JarvisOrbVisualizer = None) -> None:
        super().__init__(
            instructions="""Eres un experto en Go.
            REGLA DE ORO: Si el usuario dice "mira mi pantalla" o "analiza mi código", DEBES responder LLAMANDO a 'analyze_screen_code'.
            
            Ejemplo:
            Usuario: "Mira mi código"
            Asistente: [LLAMADA A analyze_screen_code]""",
        )
        self.vision = VisionManager()
        self.room = room
        self.visualizer = visualizer

    @function_tool
    async def ping(self):
        """Usa esta función para comprobar si las herramientas están funcionando."""
        logger.info("--- TOOL CALL: ping ---")
        return "¡Herramientas funcionando correctamente!"

    @function_tool
    async def analyze_screen_code(self, context: RunContext, query: str = "Analiza el código en pantalla"):
        """Analiza la pantalla compartida actual para extraer y entender el código fuente."""
        logger.info(f"--- TOOL CALL: analyze_screen_code (query: {query}) ---")
        if self.visualizer:
            self.visualizer.set_state("thinking")

        screenshot = await self.vision.get_latest_screenshot(self.room)
        
        if not screenshot:
            return "No puedo ver tu pantalla. Asegúrate de estar compartiendo la pantalla (Screen Share) en LiveKit."

        async with aiohttp.ClientSession() as session:
            data = aiohttp.FormData()
            data.add_field("file", screenshot, filename="screenshot.png", content_type="image/png")
            
            try:
                async with session.post(self.vision.ocr_url, data=data) as resp:
                    if resp.status == 200:
                        result = await resp.json()
                        code = result.get("code", "")
                        if code:
                            return f"He extraído el siguiente código de tu pantalla:\n\n{code}\n\n¿Qué quieres que analice de él?"
                        return "He visto la pantalla pero no he podido extraer código claro."
                    return f"Error al conectar con el servicio de OCR (Status: {resp.status})."
            except Exception as e:
                return f"Error de conexión con el servicio de OCR: {e}"

    async def llm_node(
        self,
        chat_ctx: llm.ChatContext,
        tools: list[llm.Tool],
        model_settings: ModelSettings,
    ):
        if self.visualizer:
            self.visualizer.set_state("thinking")

        buffer = ""
        async for chunk in super().llm_node(chat_ctx, tools, model_settings):
            if isinstance(chunk, llm.ChatChunk) and chunk.delta:
                delta = chunk.delta
                if delta.content:
                    buffer += delta.content
                    
                    # Buscamos si hay una llamada XML completa
                    if "</tool_call>" in buffer:
                        match = re.search(r"<tool_call>(.*?)</tool_call>", buffer, re.DOTALL)
                        if match:
                            # Texto antes de la llamada (si hay)
                            before = buffer[:match.start()].strip()
                            if before:
                                yield llm.ChatChunk(
                                    id=chunk.id,
                                    delta=llm.ChoiceDelta(role="assistant", content=before)
                                )
                            
                            # Parseamos la llamada a función
                            fnc_json = match.group(1).strip()
                            try:
                                fnc_data = json.loads(fnc_json)
                                fnc_name = fnc_data.get("name")
                                fnc_args = json.dumps(fnc_data.get("arguments", {}))
                                
                                logger.info(f"--- PATCHED XML TOOL CALL: {fnc_name} ---")
                                yield llm.ChatChunk(
                                    id=chunk.id,
                                    delta=llm.ChoiceDelta(
                                        role="assistant",
                                        tool_calls=[llm.FunctionToolCall(
                                            name=fnc_name,
                                            arguments=fnc_args,
                                            call_id=f"xml_{uuid.uuid4().hex[:8]}"
                                        )]
                                    )
                                )
                            except Exception as e:
                                logger.error(f"Error parseando llamada XML: {e} | Content: {fnc_json}")
                            
                            buffer = buffer[match.end():]
                        continue
                    
                    # Si no hay etiqueta de apertura empezada, podemos emitir el buffer
                    if "<tool_call" not in buffer and "<" not in buffer:
                        yield llm.ChatChunk(
                            id=chunk.id,
                            delta=llm.ChoiceDelta(role="assistant", content=buffer)
                        )
                        buffer = ""
                else:
                    yield chunk
            else:
                yield chunk

def prewarm(proc: JobProcess):
    proc.userdata["vad"] = silero.VAD.load()

server = AgentServer(setup_fnc=prewarm)

@server.rtc_session()
async def my_agent(ctx: JobContext):
    ctx.log_context_fields = {"room": ctx.room.name}

    # Configuración del modo de avatar: "orb", "anam", "web", o "none"
    avatar_mode = os.getenv("AVATAR_MODE", "orb").lower()
    logger.info(f"--- Modo de Avatar seleccionado: {avatar_mode.upper()} ---")

    visualizer = None
    if avatar_mode == "orb":
        visualizer = JarvisOrbVisualizer(width=512, height=512, fps=25)
        # Publicamos el track de vídeo en la sala como webcam nativa de Jarvis
        await ctx.room.local_participant.publish_track(visualizer.track)
        await visualizer.start()
        ctx.add_shutdown_callback(visualizer.stop)
        logger.info("VideoTrack de Jarvis Orb publicado con éxito en la sala.")

    # 1. Creamos la instancia de tu clase Assistant
    mi_asistente = Assistant(room=ctx.room, visualizer=visualizer)
    logger.info(f"Agente instanciado. Herramientas detectadas: {[t.info.name for t in mi_asistente.tools]}")

    session = AgentSession(
        stt=openai.STT(
            base_url="http://192.168.1.4:10300/v1", 
            api_key="none",
            model="Systran/faster-whisper-small",
            language="es"
        ),
        llm=openai.LLM(
            base_url="http://192.168.1.4:4001/v1", 
            api_key="xxxx-xxxx-xxxx-xxxx",
            model="qwen",
            tool_choice={"type": "function", "function": {"name": "analyze_screen_code"}},
            temperature=0,
            parallel_tool_calls=False,
            _strict_tool_schema=False
        ),
        tts=openai.TTS(
            base_url="http://192.168.1.4:8880/v1",
            api_key="not-needed",
            voice="im_nicola",
            model="kokoro"
        ),
        turn_detection="vad",
        vad=ctx.proc.userdata["vad"],
        min_endpointing_delay=0.6,
        preemptive_generation=False,
    )

    # Añadimos un callback para cerrar la sesión limpiamente al terminar el job
    ctx.add_shutdown_callback(session.aclose)

    await session.start(
        agent=mi_asistente,
        room=ctx.room,
        room_options=room_io.RoomOptions(
            audio_input=room_io.AudioInputOptions(noise_cancellation=None),
        ),
    )

    # Sincronización de eventos de habla con el visualizador
    if visualizer:
        @session.on("agent_started_speaking")
        def _on_agent_started_speaking():
            visualizer.set_state("speaking")

        @session.on("agent_stopped_speaking")
        def _on_agent_stopped_speaking():
            visualizer.set_state("idle")

        @session.on("user_started_speaking")
        def _on_user_started_speaking():
            visualizer.set_state("idle")

        @session.on("user_stopped_speaking")
        def _on_user_stopped_speaking():
            visualizer.set_state("thinking")

    # Modo 3: Avatar fotorrealista en la nube (Anam) si está habilitado
    if avatar_mode == "anam":
        anam_api_key = os.getenv("ANAM_API_KEY")
        anam_persona_id = os.getenv("ANAM_PERSONA_ID")
        if anam_api_key and anam_persona_id:
            try:
                from livekit.plugins import anam
                avatar = anam.Avatar(
                    persona_id=anam_persona_id,
                    api_key=anam_api_key,
                )
                await avatar.start(session, room=ctx.room)
                logger.info("Anam Avatar iniciado con éxito en la sala.")
            except ImportError:
                logger.warning("livekit-plugins-anam no está instalado. Ejecuta: uv pip install livekit-plugins-anam")
            except Exception as e:
                logger.error(f"Error iniciando Anam Avatar: {e}")

    await session.say("Hola Jose, ojos activados. Si compartes tu pantalla, puedo ayudarte a analizar el código.", allow_interruptions=True)
    
    # 3. Conectamos a la sala
    await ctx.connect()

if __name__ == "__main__":
    cli.run_app(server)
