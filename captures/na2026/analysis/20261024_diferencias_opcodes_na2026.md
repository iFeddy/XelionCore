# Mapeo Exacto de Opcodes NA2026 vs Legacy en XelionCore

| Función / Paquete | OpCode Legacy (PDB 2008/2012) | OpCode NA2026 (Captura Oficial 2026) | Tamaño Wire / Payload | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| **Semilla XOR (C2S/S2C)** | `0x0807` | `0x0807` | 5B / 4B | Handshake inicial idéntico (`PROTO_NC_MISC_SEED_CMD`) |
| **Check Versión (Req)** | `0x0C65` | `0x0C2C` | 35B / 34B | `PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ` (Hash MD5 32B) |
| **Check Versión (Ack)** | `0x0C67` | `0x0C2D` | 4B / 3B | `PROTO_NC_USER_CLIENT_VERSION_CHECK_ACK` (`00` = OK) |
| **Login Req (Credenciales)** | `0x0C5A` / `0x0C03` | `0x0C01` | 354B / 351B | `PROTO_NC_USER_LOGIN_REQ` (32B pad, 260B user, 32B pass, spawnApps) |
| **Login Fail (Clave errónea)** | `0x0C09` | `0x0C07` | 6B / 4B | `PROTO_NC_USER_LOGINFAIL_ACK` (`0x0045` = Clave inválida) |
| **Login Ack (OK pre-mundos)** | - | `0x0C40` | 4B / 3B | Retorna `0x00` (Éxito de autenticación) |
| **Login Ack (Mundos)** | `0x0C0A` | `0x0C06` | Dinámico (295B) | `PROTO_NC_USER_LOGIN_ACK` (1B count + 4B header + N * 18B mundos) |
| **Sesión / Token de Cuenta** | - | `0x0C3E` | 67B / 66B | Token de sesión para TOS y transferencias |
| **Ping de Sesión TOS** | - | `0x0C34` / `0x0C35` | 11B / 10B | Ping / Pong de sesión en ventana de TOS |
| **Aceptación TOS** | - | `0x0C3F` / `0x0C41` | 68B / 3B | Envío de hash de aceptación y confirmación |
| **Refresco Mundos (Req)** | `0x0C1B` | `0x0C18` | 3B / 2B | Solicitud de listado tras TOS |
| **Refresco Mundos (Ack)** | `0x0C1C` | `0x0C19` | Dinámico (295B) | `PROTO_NC_USER_WORLDSTATUS_ACK` (mismo formato que 0x0C06) |
| **Selección de Mundo (Req)** | `0x0C0B` | `0x0C0A` | 4B / 3B | `PROTO_NC_USER_WORLDSELECT_REQ` (`0A 0C [WorldID]`) |
| **Selección de Mundo (Ack)** | `0x0C0C` | `0x0C0B` | 87B / 86B | `PROTO_NC_USER_WORLDSELECT_ACK` (Status, IP:Port 9310, Token 64B) |
