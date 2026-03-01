# Ejemplo 10: Mundo Real y Temperatura

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 12 de 13: Aplicación Teórica del Sesgo)*

---

### Línea Argumental
En la gráfica anterior `10` jugamos con el desplazamiento, pero no sabíamos muy bien por qué. Este ejemplo contextualiza esa idea introduciendo el tiempo (Horas de 0 a 24).

### ¿Qué hace este código?
- Creamos un eje X que representa las horas de un día.
- Creamos Neuronas donde la entrada son las horas. Una con `Bias=0` y otra con el `Bias=0.5`.
- Grafica el resultado simulando la temperatura. Aquí podemos ver físicamente qué ocurre cuando un sesgo "desplaza" el aprendizaje hacia unas horas más tempranas o tardías del día, forzando a aprender que "12:00" y "18:00" den salidas diferentes para los mismos pesos.

### Siguiente Paso ->
[10-neuron-js-golang](../10-neuron-js-golang) - ¡El final del viaje! Combinar todo lo aprendido (Activaciones, Backward/Forward pasess manuales y pesos) y montarlo en un simulador visual con inercia en vivo usando **Ebiten**, donde veremos cómo se comportan al intentar simular el difícil juego de resolver XOR al aire.
