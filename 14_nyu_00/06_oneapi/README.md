# Matrix Multiplication with oneAPI (SYCL) & Go

Este proyecto demuestra cómo realizar una multiplicación de matrices acelerada por GPU utilizando **Intel oneAPI (SYCL)** e integrada con una aplicación **Go** mediante **CGO**.

## Requisitos

- **Intel oneAPI Base Toolkit**: Específicamente el compilador `icpx` (Intel C++ Compiler).
- **GPU Compatible**: El ejemplo está probado en una **Intel Arc B580**, pero debería funcionar en cualquier GPU compatible con SYCL (Level Zero u OpenCL).
- **Go**: Instalado y configurado.

## Estructura del Proyecto

- `matrix_sycl.cpp`: Contiene el kernel de SYCL que realiza la multiplicación de matrices en la GPU.
- `matrix_sycl.h`: Define la interfaz C necesaria para conectar C++ con Go.
- `main.go`: Aplicación principal en Go que:
  - Genera matrices aleatorias.
  - Llama al kernel de la GPU.
  - Realiza una multiplicación en CPU (Go puro) para comparar.
  - Verifica que los resultados coincidan.
- `Makefile`: Automatiza la compilación del kernel SYCL en una librería compartida (`.so`).

## Cómo Ejecutar

1. **Cargar el entorno de oneAPI**:
   ```bash
   source /opt/intel/oneapi/setvars.sh
   ```

2. **Compilar y Ejecutar**:
   ```bash
   make && go build -o matrix_go main.go && ./matrix_go
   ```

## Detalles Técnicos: oneAPI y `setvars.sh`

Para que este ejemplo funcione, es fundamental entender qué hace `source /opt/intel/oneapi/setvars.sh`:

### 1. Disponibilidad de Instrucciones y Compilador
`setvars.sh` configura la variable de entorno `PATH` para que el compilador `icpx` sea accesible. `icpx` es el compilador de C++ de Intel que soporta **SYCL**. 
- Al usar el flag `-fsycl` en el `Makefile`, el compilador realiza una "compilación en dos pasos" (Fat Binary): genera código para la CPU y kernels específicos para la GPU (en formato SPIR-V).

### 2. Visibilidad de Librerías (Headers y Runtime)
- **Headers**: `setvars.sh` configura `CPATH` o variables internas de `icpx` para encontrar `<sycl/sycl.hpp>`.
- **Runtime**: Configura `LD_LIBRARY_PATH`. Esto es crítico para que:
    - El ejecutable encuentre las librerías de tiempo de ejecución de SYCL (`libsycl.so`).
    - El "Unified Runtime" de Intel pueda localizar los *loaders* de **Level Zero** o **OpenCL** que se comunican con el driver de la GPU.

### 3. Integración con Go (CGO)
En `main.go`, el bloque de CGO:
```go
#cgo LDFLAGS: -L. -lmatrix_sycl -Wl,-rpath,.
```
Le dice a Go que busque nuestra librería `libmatrix_sycl.so` en el directorio actual. Sin embargo, `libmatrix_sycl.so` depende de las librerías de Intel. Por eso, al ejecutar `./matrix_go`, el sistema operativo necesita que `LD_LIBRARY_PATH` esté correctamente configurado (lo cual hace `setvars.sh`).

## Glosario de Conceptos Clave

Para profundizar en la tecnología utilizada:

- **`icpx`**: Es el compilador de C++ de Intel basado en LLVM. A diferencia de `g++`, entiende las extensiones de SYCL y puede generar código para múltiples arquitecturas (CPU, GPU, FPGA) simultáneamente.
- **`-fsycl`**: Este flag es "mágico". Activa el soporte de SYCL en el compilador. Le indica que debe buscar kernels (bloques de código para GPU), compilarlos para el dispositivo objetivo y empaquetarlos dentro del binario final (lo que se conoce como *Fat Binary*).
- **`SPIR-V` (Standard Portable Intermediate Representation)**: Es un formato binario intermedio (similar al bytecode de Java) que utiliza el compilador. En lugar de generar código máquina para una GPU específica al compilar, genera SPIR-V. Luego, el driver de la GPU convierte ese SPIR-V en instrucciones reales para tu hardware específico al momento de ejecutar.
- **Extensión `.hpp`**: SYCL es una especificación "Single-Source", lo que significa que el código de la CPU y la GPU están en el mismo archivo. Gran parte de SYCL está implementado en cabeceras (`headers`), por lo que verás mucho uso de `.hpp`. El compilador usa estas cabeceras para entender los tipos de datos y la sintaxis de los kernels.
- **`libsycl.so`**: Es la librería de tiempo de ejecución (runtime). Es la responsable de gestionar las colas de ejecución (`sycl::queue`), mover datos entre la RAM y la memoria de la GPU, y lanzar los kernels.
- **Level Zero**: Es la interfaz de bajo nivel de Intel para sus GPUs. Es más moderna y eficiente que OpenCL, diseñada específicamente para extraer el máximo rendimiento de hardware como la **Intel Arc B580**. SYCL usa Level Zero por debajo por defecto en GPUs Intel modernas.
- **OpenCL**: Un estándar abierto más antiguo para computación paralela. SYCL también puede usar OpenCL como "backend" si Level Zero no está disponible o para compatibilidad con hardware de otros fabricantes.

## Resultados Esperados (Tamaño 2048x2048)

En una **Intel Arc B580**, se obtienen los siguientes resultados aproximados:

- **GPU Time**: ~0.8s
- **CPU Time**: ~24.5s
- **Speedup**: **~30x**

El programa imprimirá el nombre del dispositivo GPU detectado por SYCL antes de comenzar.
