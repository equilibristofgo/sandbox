# Ejemplo 09: Poniendo Cara a las Matemáticas

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 10 de 13: Gonum Plot y Gráficas 2D)*

---

### Línea Argumental
La ciencia de datos sin gráficas es un vuelo a ciegas. En el capítulo anterior vimos los números de la función de activación Sigmoide. Ahora usaremos otra herramienta bestial del ecosistema: `gonum/plot`.

### ¿Qué hace este código?
- Genera un rango de 100 pruebas desde entrada negativa a positiva `[-10, 10]`.
- Mapea esos valores usando nuestra función Sigmoide.
- Instancia un `plotter` y guarda directamente la famosa rampa en `sigmoid.png`. 

Aquí puedes visualizar físicamente cómo la función ahoga cualquier extremo gigantesco obligándolo a terminar entre 0 y 1.

### Siguiente Paso ->
[10-neuron](../10-neuron) - Ahora que sabemos cómo es la base pura, vamos a graficar qué pasa cuando le metemos pesos y sesgos, comparando la curva pura (Azul) contra la curva asimétrica de la neurona completa.
