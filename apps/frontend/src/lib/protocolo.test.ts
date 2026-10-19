import { describe, it, expect } from "vitest";
import { respostasVazias, NIVEIS, SENTINELAS } from "./protocolo";

// O cálculo dos escores e da matriz MIV-CFS pertence ao backend (Go) e é
// testado lá. Aqui guardamos o contrato do formato que o frontend envia.

describe("Contrato de Respostas", () => {
  it("respostasVazias tem exatamente os campos esperados pelo backend", () => {
    const r = respostasVazias();
    expect(Object.keys(r).sort()).toEqual([
      "abvd",
      "aerobica",
      "aivd",
      "alcance",
      "autopercepcao",
      "cognicao",
      "comorbidades",
      "comunicacao",
      "easfi",
      "humor",
      "idade",
      "marcha",
      "sentinelas",
    ]);
    expect(r.idade).toBeNull();
    expect(r.autopercepcao).toBeNull();
    expect(r.abvd).toBeNull();
    expect(r.aerobica).toBeNull();
    expect(r.comorbidades).toBeNull();
    expect(r.aivd).toHaveLength(3);
    expect(r.cognicao).toHaveLength(3);
    expect(r.humor).toHaveLength(2);
    expect(r.alcance).toHaveLength(2);
    expect(r.marcha).toHaveLength(3);
    expect(r.comunicacao).toHaveLength(2);
    expect(r.easfi).toHaveLength(5);
    expect(r.sentinelas).toHaveLength(5);
    expect(
      [
        ...r.aivd,
        ...r.cognicao,
        ...r.humor,
        ...r.alcance,
        ...r.marcha,
        ...r.comunicacao,
        ...r.sentinelas,
      ].every((v) => v === null),
    ).toBe(true);
    expect(r.easfi.every((v) => v === null)).toBe(true);
  });
});

describe("Dados de exibição", () => {
  it("níveis 1–4 completos", () => {
    expect(Object.keys(NIVEIS).sort()).toEqual(["1", "2", "3", "4"]);
    for (const n of Object.values(NIVEIS)) {
      expect(n.nome).toBeTruthy();
      expect(n.conduta).toBeTruthy();
      expect(n.reavaliacao).toBeTruthy();
    }
  });
  it("5 marcadores sentinela", () => {
    expect(SENTINELAS).toHaveLength(5);
    expect(SENTINELAS.every((s) => s.length > 10)).toBe(true);
  });
});
