# Ejemplo 08: La Activación Sigmoide

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 9 de 13: Diseccionando la Activación)*

---

### Línea Argumental
En `07-train` aplicamos ciegamente la función `sigmoid(x)` para que la red aprendiera. Las matemáticas puras nos dicen que es `1 / (1 + exp(-x))`. ¿Pero por qué? Las funciones de activación estrujan los infinitos números matemáticos entre rangos lógicos manejables (como 0 y 1 para la Sigmoide).

### ¿Qué hace este código?
- Aísla la función matemática de la Sigmoide.
- Ejecuta e imprime una salida basándose en valores simples.

### Siguiente Paso ->
[09-sigmoid-graph](../09-sigmoid-graph) - Los números en consola están bien, pero en ML "ver" los datos es crítico. Vamos a usar **Gonum Plot** para darle forma a esa fórmula y ver su clásica rampa en "S".
