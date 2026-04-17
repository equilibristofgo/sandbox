# Ejemplo 10: La Complejidad Funcional

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 11 de 13: El Poder del Sesgo y Peso)*

---

### Línea Argumental
Teníamos la función pura de `09`. Todo iba de 0 a 1 suavemente. Pero, ¿qué pasa cuándo el Perceptrón que diseñamos en capítulos anteriores aprende y modifica sus `pesos` y `sesgos` que vimos en la manualidad de `03` y `07`?

### ¿Qué hace este código?
- Vuelve a montar un struct `Neuron`.
- Genera un plot doble (`gonum/plot`): uno para la Sigmoide estándar y otro para nuestro cálculo Neuronal Ponderado.
- La línea Azul (Sigmoide) cortará el eje, mientras que nuestra línea Roja, modificada por `neurona([]float64{9.0}, 0.7)` nos enseñará cómo la gráfica entera "se desplaza" y "se escala". El peso cambia la pendiente. El sesgo la decanta.

### Siguiente Paso ->
[10-neuron-temperatura](../10-neuron-temperatura) - Todo esto está muy bien en abstracto. Vamos a darle sentido semántico a ese desplazamiento prediciendo el calor del día.
