import logging
import io
import aiohttp
import re
import json
import uuid
from dotenv import load_dotenv
from PIL import Image
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
load_dotenv(".env.local")

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
    def __init__(self, room: rtc.Room) -> None:
        super().__init__(
            instructions="""Eres un experto en Go.
            REGLA DE ORO: Si el usuario dice "mira mi pantalla" o "analiza mi código", DEBES responder LLAMANDO a 'analyze_screen_code'.
            
            Ejemplo:
            Usuario: "Mira mi código"
            Asistente: [LLAMADA A analyze_screen_code]""",
        )
        self.vision = VisionManager()
        self.room = room

    @function_tool
    async def ping(self):
        """Usa esta función para comprobar si las herramientas están funcionando."""
        logger.info("--- TOOL CALL: ping ---")
        return "¡Herramientas funcionando correctamente!"

    @function_tool
    async def analyze_screen_code(self, context: RunContext, query: str = "Analiza el código en pantalla"):
        """Analiza la pantalla compartida actual para extraer y entender el código fuente."""
        logger.info(f"--- TOOL CALL: analyze_screen_code (query: {query}) ---")
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
    
    # 1. Creamos la instancia de tu clase Assistant
    # Esta instancia contiene las herramientas (@function_tool)
    mi_asistente = Assistant(room=ctx.room)
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

    await session.say("Hola Jose, ojos activados. Si compartes tu pantalla, puedo ayudarte a analizar el código.", allow_interruptions=True)
    
    # 3. Conectamos a la sala
    await ctx.connect()

if __name__ == "__main__":
    cli.run_app(server)
