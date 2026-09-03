# Panic Recover: Runtime Defender

> Submission for **GopherCon LATAM 2026 Mini Game Jam**  
> **Theme**: *Panic (or recover?)*

A fast-paced, retro-cyberpunk top-down space shooter written in Go with **Ebitengine v2**. Blast runtime errors and concurrency bugs across an endless scrolling call stack. Trigger high-octane **Panic Mode** when damaged and find `recover()` tokens before the stack trace unwinds!

---

## Autores
- **Cesar & Antigravity**

---

## Controles

| Ação | Teclado / Mouse | Controle (Gamepad) |
| :--- | :--- | :--- |
| **Mover Nave** | `W`, `A`, `S`, `D` ou Setas | Analógico Esquerdo / D-Pad |
| **Atirar** | Barra de Espaço, `J` ou Botão Esquerdo do Mouse | Botão Sul (`A` / `X`) |
| **Ativar Recover / Bomba** | `K`, `E` ou Botão Direito do Mouse | Botão Leste (`B` / `Círculo`) |
| **Reiniciar Jogo** | `R` ou Espaço (na tela de Game Over) | Botão Start |
| **Tela Cheia** | `F11` ou `Alt + Enter` | - |

---

## Instruções de Execução

### Pré-requisitos
- **Go 1.22+** (ou Go 1.26+)
- Bibliotecas gráficas padrão do Linux (libasound2-dev, libgl1-mesa-dev, libxcursor-dev, libxi-dev, libxinerama-dev, libxrandr-dev, libxxf86vm-dev se compilando nativo).

### Executar no Navegador (WebAssembly - Recomendado sem dependências de Cgo)
```bash
cd panic-recover
go run ./cmd/server
# Abra http://localhost:8080 no seu navegador!
```

### Executar diretamente no Desktop (Nativo)
```bash
cd panic-recover
go run ./cmd/game
```

### Compilar binário nativo
```bash
cd panic-recover
go build -o panic-recover ./cmd/game
./panic-recover
```

### Recompilar para WebAssembly
```bash
cd panic-recover
GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/game
```
