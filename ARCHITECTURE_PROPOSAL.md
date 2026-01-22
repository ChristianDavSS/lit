# Propuesta de Arquitectura Hexagonal (Ports and Adapters)

## Problema Actual
El proyecto actual presenta un acoplamiento fuerte entre capas. Específicamente, módulos de bajo nivel (como `internals/ast`) dependen de la configuración global (`config`), la cual a su vez tiene lógica de interfaz de usuario (preguntar input al usuario). Esto crea ciclos de dependencias y hace que el código sea difícil de probar y mantener.

## Solución: Arquitectura Hexagonal

La Arquitectura Hexagonal (o Puertos y Adaptadores) busca aislar la lógica de negocio (el "Core") de los detalles de implementación (Frameworks, Bases de Datos, UI, Sistemas de Archivos).

### Estructura Propuesta

```
.
├── cmd/                # Puntos de entrada de la aplicación (Main)
│   └── lit/
│       └── main.go
├── internal/           # Código privado de la aplicación
│   ├── domain/         # Entidades y Lógica pura de negocio (Sin dependencias externas)
│   │   ├── model.go    # Definición de estructuras (FunctionData, Config, etc.)
│   │   └── ports.go    # Interfaces (Scanner, ConfigProvider, UI)
│   ├── service/        # Casos de uso (Orquestación)
│   │   ├── scanner.go  # Lógica de escaneo de archivos
│   │   └── fixer.go    # Lógica de corrección
│   └── adapter/        # Implementación de las interfaces (Puertos)
│       ├── config/     # Adaptador para leer/escribir config.json
│       ├── cli/        # Comandos de Cobra (UI)
│       ├── fs/         # Acceso a sistema de archivos (si es necesario abstraer)
│       └── analysis/   # Implementación del análisis AST (TreeSitter)
```

### Cambios Clave

1.  **Desacoplar el Core:** El analizador (`scanner`) no leerá variables globales. Recibirá la configuración necesaria (ej. convención de nombres) como argumentos o a través de una interfaz.
2.  **Separar Configuración de UI:** El paquete `config` solo se encargará de leer/escribir datos. La interacción con el usuario (preguntar convención) se moverá a la capa de UI (`cmd` o `adapter/cli`).
3.  **Inversión de Dependencias:** Los módulos de alto nivel (Core) definen interfaces (Puertos) que los módulos de bajo nivel (Adaptadores) implementan.

Esta estructura facilita:
- **Testabilidad:** Se puede probar el `scanner` pasando una configuración en memoria sin tocar archivos.
- **Mantenibilidad:** Cambiar la librería de CLI o el formato de configuración no afecta la lógica de análisis.
