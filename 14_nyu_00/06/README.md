# Ejemplo 06: Gonum al Rescate

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 6 de 13: Multiplicación Matricial Segura)*

---

### Línea Argumental
Tras el desastre de los índices manuales del capítulo `05`, llega la caballería: **Gonum**.
Gonum nos proporciona tipos fuertemente tipados de matrices densas (`mat.Dense`) y vectores (`mat.VecDense`). Al utilizarlos, abstraemos todo el dolor algebraico y ganamos la velocidad de las rutinas BLAS subyacentes.

### ¿Qué hace este código?
- Declara las entradas reales, pesos y sesgos ya no como simples slices `[]float64`, sino como matrices matemáticas reales `mat.NewDense`.
- Delega la multiplicación ponderada (`Weights * Inputs`) a `result.MulVec()`, que internamente valida las dimensiones y asegura que no haya desbordamientos de meemoria antes de operar.

### Multiplicacion de matrices en GPU intel
- Ya vimos en el [02-ManualTensor](../02) como se podia hacer uso de las instrucciones AVX desde ensamblador para acelerar la multiplicacion de matrices... Y se que habia por ahi una libreria que te lo ocultaba, pero ahora no la encuentro (a ver si gemini me ayuda)
- El caso es que aqui haremos una prueba de llamada a traves de CGO en este caso, a funciones en C para usar las librerias de intel para multiplicacion de matrices en GPU.


### Siguiente Paso ->
[07-GonumBias](../07) - Una ecuación lineal no está completa sin el sesgo (Bias). Usaremos Gonum para finalizar el pase.
