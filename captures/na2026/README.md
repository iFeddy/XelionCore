# Capturas y Trazas de Paquetes NA2026

Este directorio contiene las capturas de tráfico de red y volcados de paquetes del cliente oficial **NA2026** conectando contra **XelionCore**.

## Estructura

- `pcaps/`: Capturas de red crudas (`.pcap` / `.pcapng`) generadas mediante `tcpdump` o Wireshark. Cada archivo se nombra con timestamp: `YYYYMMDD_HHMMSS_<descripcion>.pcap`.
- `raw/`: Trazas y volcados en texto plano (`hexdump` y opcodes oficiales) generados por el tracer de `XelionCore`.
- `analysis/`: Documentación técnica de estructuras, offsets y diferencias de NA2026 respecto a la documentación oficial.
