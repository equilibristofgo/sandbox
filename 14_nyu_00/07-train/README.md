# Ejemplo 07-train: Aprendiendo de Verdad (XOR)

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 8 de 13: Backpropagation con Gonum)*

---

### Línea Argumental
Hasta ahora, todas las redes han tenido "pesos fijos". Eran tontas. 
En este capítulo, implementamos el corazón del Machine Learning: El **Backpropagation** (Propagación hacia atrás). Vamos a retar a nuestra red escrita en Go con el problema XOR, un hito que demuestra que una red multicapa puede aprender geometrías no lineales.

### ¿Qué hace este código?
1. **Inicializa aleatoriamente** los pesos (El cerebro de un bebé).
2. Para cada "Epoch" (Iteración de época):
   - **Forward Pass:** Intenta adivinar la solución usando Gonum y activaciones Sigmoide.
   - **Error:** Calcula cuánto se ha equivocado restando la salida de los Targets reales.
   - **Backward Pass:** Calcula los gradientes (cuánta culpa tiene cada peso usando la derivada de la Sigmoide).
   - **Weight Update:** Modifica los pesos un poquito en dirección contraria al error (`Learning Rate`).

### HUGOT
- Pequeña referencia a librerias como Hugot para inferencia, transformer y mas...
- Esto habre otra via, pero un pequeño ejemplo aqui puede ser interesante...
- Ver lo que ya he montado en otros proyectos..


### Siguiente Paso ->
[08-sigmoid](../08-sigmoid) - Hemos usado funciones de activación ("Sigmoide") para permitir el aprendizaje no lineal. Vamos a estudiar a fondo qué es geométricamente esa rampa.
