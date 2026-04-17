# Matrix Multiplication with Level Zero & Go

Este ejemplo demuestra cómo realizar una multiplicación de matrices utilizando directamente la **Intel Level Zero API**, permitiendo un control total sobre el hardware sin pasar por abstracciones de alto nivel como SYCL.

## Requisitos

- **Intel Level Zero Loader**: Librería `libze_loader.so`.
- **Intel Level Zero Headers**: Instalados en `/usr/include/level_zero`.
- **ocloc**: Herramienta de compilación offline para kernels Intel (incluida en oneAPI).
- **GPU Compatible**: Intel Arc B580 o similar.

## Estructura del Proyecto

- `kernel.cl`: Kernel escrito en **OpenCL C**. Se compila a **SPIR-V** (formato binario intermedio) usando `ocloc`.
- `matrix_ze.cpp`: Wrapper en C++ que realiza la orquestación de bajo nivel:
  - Inicialización del driver y dispositivo.
  - Creación de contexto y cola de comandos.
  - Carga manual del binario SPIR-V.
  - Gestión de memoria **USM (Unified Shared Memory)**.
  - Lanzamiento del kernel y sincronización.
- `main.go`: Aplicación Go que interactúa con el wrapper mediante CGO.
- `Makefile`: Gestiona la compilación del kernel y la librería compartida.

## Cómo Ejecutar

1. **Cargar el entorno de oneAPI**:
   ```bash
   source /opt/intel/oneapi/setvars.sh
   ```

2. **Compilar y Ejecutar**:
   ```bash
   make && go build -o matrix_go main.go && ./matrix_go
   ```

## Diferencias con SYCL

A diferencia de SYCL (donde el runtime gestiona casi todo), en Level Zero:
- **Control Manual**: Tú decides exactamente cuándo se crea el contexto, cuándo se mueven los datos y cuándo se destruyen los recursos.
- **SPIR-V Directo**: El kernel se compila externamente a SPIR-V (`kernel.spv_bmg.spv`) y se carga explícitamente en el programa.
- **USM**: Se utiliza memoria compartida entre CPU y GPU de forma explícita (`zeMemAllocShared`).
