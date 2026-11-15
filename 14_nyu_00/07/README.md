# Ejemplo 07: Completando la Ecuación Lineal

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 7 de 13: Añadiendo el Sesgo con Gonum)*

---

### Línea Argumental
Sabiendo multiplicar tensores de forma segura con Gonum (`06`), ahora añadimos la inercia base de la neurona: el Sesgo (`Bias`). Matemáticamente, el paso frontal de una neurona (antes de activarse) es `y = Wx + b`. Gonum nos lo resuelve en dos líneas de código limpias.

### ¿Qué hace este código?
- Realiza el producto punto matricial `result.MulVec(weights, x)`.
- Inmediatamente le suma el vector de sesgos con `result.AddVec(result, bias)`.
- Todo ocurre usando punteros sobre el mismo contenedor `result`, optimizando las reservas de memoria (GC friendly), crucial en arquitecturas de Go sometidas a alta carga.

### Siguiente Paso ->
[07-train](../07-train) - Ya sabemos hacer el "Forward Pass". Ha llegado la hora de enseñarle a la red a pensar comprobando su error hacia atrás (Backpropagation) y resolviendo el clásico juego XOR.
