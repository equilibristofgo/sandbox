# Ejemplo Ebiten: El Simulator Interactivo XOR

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 13 de 13: La Inteligencia a la vista)*

---

### Línea Argumental
Aquí termina nuestra introducción al cerebro informático en Go. Nos despedimos abandonando las gráficas estáticas PNG de GonumPlot, dándole vida a una red **Multi-Layer Perceptron (MLP)** sobre el motor de juegos 2D `Ebiten`. Vamos a retomar la arquitectura con backpropagation de `07-train`, pero en tiempo real y modificable por el usuario.

### ¿Qué hace este código?
- Usa un esquema Input > Hidden(tanh) > Output(sigmoid) implementado puramente en Go sin matrices (para entender visualmente las conexiones del grafo real).
- Dibuja círculos (nodos) alterando su transparencia (alpha) según la intensidad de la activación.
- Dibuja líneas (pesos) y su grosor cambia según el valor absoluto del aprendizaje; los colores rojo/verde indican pesos neagtivos/positivos.
- **Tú eres el entrenador:** Pulsando *Arriba* provocas `Forward` y `Backwards` passes masivos para las puertas lógicas XOR. Pulsando *Abajo* desentrenas usando un Learning Rate (`lr`) que puedes acelerar/frenar con las flechas laterales, viendo la curva de error reaccionar en tiempo real a tus exigencias.

### Fin del Bloque 1
Con esto termina el bloque básico para la serie de Machine Learning. El siguiente paso ya requiere salir de nuestro sandbox en busca de tensores hardware y LLMs.