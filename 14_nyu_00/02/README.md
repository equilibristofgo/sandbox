# Ejemplo 02: El Lienzo en Blanco

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 2 de 13: Matrices Manuales)*

---

### Línea Argumental
Ya tenemos los datos paralelos de nuestro Worker (`01`), pero para que el Machine Learning ocurra necesitamos estructurarlo. El bloque fundamental es el **Tensor**, que en dos dimensiones (2D) no es más que una matriz. 

¿O no? ... [Aqui](https://www.youtube.com/watch?v=1GwAEnegaRs) un video donde cuentan que es una matriz y que es un tensor, porque ... no son exactamente lo mismo... [Aqui](https://www.youtube.com/watch?v=RxoZqQmHqbY) en que se diferencian... Pero merece tenerlo claro: escalar, vector, matriz ... el numero de "dimensiones"...

Pero como somos Go puristas al principio, probaremos a implementarlo a mano.

### ¿Qué hace este código?
- Instancia una pequeña Matriz de 2x2 usando Slices anidados de Go nativo (`[][]float64`).
- Demuestra que iterar en Go es sencillo para matrices pequeñas.
- **[Bonus] AVX Implementation**: Incluye una suma de vectores acelerada por hardware usando instrucciones SIMD (AVX) mediante `avo`.

Se ha incluido un test de rendimiento (`main_test.go`) para comparar la velocidad de Go nativo frente a la implementación optimizada (4x unrolled) en ensamblador AVX:

| Tamaño Vector | Go (ns/op) | AVX (ns/op) | Mejora |
|---------------|------------|-------------|--------|
| 1.024 elementos | 220.0 | 121.2 | **1.8x** |
| 65.536 elementos | 16.772 | 5.071 | **3.3x** |
| 1.048.576 elementos | 304.598 | 269.260 | **1.1x** |

> [!TIP]
> La implementación AVX ha sido optimizada mediante **Loop Unrolling (4x)**, procesando 32 elementos por cada salto del bucle, lo que maximiza el rendimiento de la pipeline de ejecución.

> [!NOTE]
> Para tamaños de datos pequeños/medianos, Go suele ser más eficiente debido a la "auto-vectorización" del compilador y el menor coste de llamada a función. El ensamblador manual brilla en operaciones masivas y complejas donde el compilador no puede optimizar más.

### Siguiente Paso ->
[03-FirstNeuron](../03) - Usando este concepto manual, montaremos la primera representación física del "Cerebro", un Perceptrón.
