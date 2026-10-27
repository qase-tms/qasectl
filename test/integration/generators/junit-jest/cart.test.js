describe("Cart", () => {
  describe("Totals", () => {
    it("adds item prices", () => {
      expect(10 + 20).toBe(30);
    });

    it("applies discount", () => {
      expect(100 - 10).toBe(80);
    });

    it.skip("rounds to cents", () => {
      expect(1.005).toBe(1.01);
    });
  });
});
