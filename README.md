# GoCh8p

My attempt at a CHIP-8 "emulator"

![example](working.png)


References:
* https://www.youtube.com/watch?v=V3ORJUS9YP4
  
![Memory Map](memory_map.png)
## Status
- [X] Tick Rate (60Hz)
- [X] Memory (4096 bytes)
- [X] Stack (16 levels)
- [X] Registers (V0-F)
- [X] Counters (PC, SP, I)
- [ ] Keyboard (16)
- [X] Display (64x32)
- [X] Fontset (5x8, 0-F)
- [ ] Timers (Sound, Delay)
- [ ] Opcodes ![opcodes](opcodes.png)
  - [ ] 0x00E0: "CLS",
  - [ ] 0x00EC: "YIELD",
  - [ ] 0x00EE: "RET",
  - [ ] 0x1000: "1NNN",
  - [ ] 0x2000: "2NNN",
  - [ ] 0x3000: "3XNN",
  - [ ] 0x4000: "4XNN",
  - [ ] 0x5000: "5XY0",
  - [X] 0x6000: "6XNN",
  - [ ] 0x7000: "7XNN",
  - [ ] 0x8000: "8XY0",
  - [ ] 0x8001: "8XY1",
  - [ ] 0x8002: "8XY2",
  - [ ] 0x8003: "8XY3",
  - [ ] 0x8004: "8XY4",
  - [ ] 0x8005: "8XY5",
  - [ ] 0x8006: "8XY6",
  - [ ] 0x8007: "8XY7",
  - [ ] 0x800E: "8XYE",
  - [ ] 0x9000: "9XY0",
  - [ ] 0xA000: "ANNN",
  - [ ] 0xB000: "BNNN",
  - [ ] 0xC000: "CXNN",
  - [X] 0xD000: "DXYN",
  - [ ] 0xE09E: "EX9E",
  - [ ] 0xE0A1: "EXA1",
  - [ ] 0xF007: "FX07",
  - [ ] 0xF00A: "FX0A",
  - [ ] 0xF015: "FX15",
  - [ ] 0xF018: "FX18",
  - [ ] 0xF01E: "FX1E",
  - [ ] 0xF029: "FX29",
  - [ ] 0xF033: "FX33",
  - [ ] 0xF055: "FX55",
  - [ ] 0xF065: "FX65",
