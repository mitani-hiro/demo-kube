import { test, expect } from "@playwright/test";

test("GET /api/health returns 200", async ({ request }) => {
  const res = await request.get("/api/health");
  expect(res.status()).toBe(200);
});

test("POST /api/hoge returns 200 with user payload", async ({ request }) => {
  const res = await request.post("/api/hoge");
  expect(res.status()).toBe(200);

  const body = await res.json();
  expect(body.id).toBe(999);
  expect(body.name).toBe("Taro Yamada");
  // hostname は応答した producer Pod 名（LB 分散検証用フィールド）
  expect(typeof body.hostname).toBe("string");
  expect(body.hostname.length).toBeGreaterThan(0);
});
