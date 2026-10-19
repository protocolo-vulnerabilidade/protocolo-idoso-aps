import { describe, it, expect } from "vitest";
import {
  pontosIdade, dominiosIVCF, escoreIVCF, classificarIVCF, classificarEASFI, nivelMatriz,
  aplicarSentinelas, calcular, respostasVazias,
} from "./protocolo";

describe("IVCF-20", () => {
  it("idade", () => {
    expect(pontosIdade(74)).toBe(0);
    expect(pontosIdade(75)).toBe(1);
    expect(pontosIdade(84)).toBe(1);
    expect(pontosIdade(85)).toBe(3);
  });
  it("AIVD tem teto de 4 mesmo com três Sim", () => {
    const r = respostasVazias();
    r.aivd = [true, true, true];
    expect(dominiosIVCF(r).aivd).toBe(4);
  });
  it("ABVD: apenas banho vale 6", () => {
    const r = respostasVazias();
    r.abvd = true;
    expect(dominiosIVCF(r).abvd).toBe(6);
  });
  it("cognição: terceira questão vale 2 (máx 4)", () => {
    const r = respostasVazias();
    r.cognicao = [true, true, true];
    expect(dominiosIVCF(r).cognicao).toBe(4);
  });
  it("humor soma as duas sem teto (máx 4)", () => {
    const r = respostasVazias();
    r.humor = [true, true];
    expect(dominiosIVCF(r).humor).toBe(4);
  });
  it("marcha/continência soma 3 itens sem teto (máx 6)", () => {
    const r = respostasVazias();
    r.marcha = [true, true, true];
    expect(dominiosIVCF(r).marcha).toBe(6);
  });
  it("comorbidades é pergunta única de 4 pontos", () => {
    const r = respostasVazias();
    r.comorbidades = true;
    expect(dominiosIVCF(r).comorbidades).toBe(4);
  });
  it("escore máximo é exatamente 40", () => {
    const r = respostasVazias();
    r.idade = 90;
    r.autopercepcao = "ruim";
    r.aivd = [true, true, true];
    r.abvd = true;
    r.cognicao = [true, true, true];
    r.humor = [true, true];
    r.alcance = [true, true];
    r.aerobica = true;
    r.marcha = [true, true, true];
    r.comunicacao = [true, true];
    r.comorbidades = true;
    expect(escoreIVCF(r)).toBe(40);
  });
  it("faixas", () => {
    expect(classificarIVCF(6)).toBe("baixo");
    expect(classificarIVCF(7)).toBe("moderado");
    expect(classificarIVCF(15)).toBe("alto");
  });
});

describe("EASFI-BR", () => {
  it("faixas", () => {
    expect(classificarEASFI(9)).toBe("adequada");
    expect(classificarEASFI(10)).toBe("risco");
    expect(classificarEASFI(15)).toBe("problema");
  });
});

describe("Matriz e sentinelas", () => {
  it("matriz", () => {
    expect(nivelMatriz("baixo", "adequada")).toBe(1);
    expect(nivelMatriz("moderado", "problema")).toBe(4);
    expect(nivelMatriz("alto", "adequada")).toBe(3);
  });
  it("sentinela leve leva a 3", () => {
    expect(aplicarSentinelas(1, [false, false, true, false, false])).toBe(3);
  });
  it("sentinela grave leva a 4", () => {
    expect(aplicarSentinelas(1, [false, true, false, false, false])).toBe(4);
  });
  it("cálculo completo", () => {
    const r = respostasVazias();
    r.idade = 80; r.easfi = [1, 1, 1, 1, 1];
    expect(calcular(r).nivel).toBe(1);
  });
});
