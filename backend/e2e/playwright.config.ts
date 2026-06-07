import { defineConfig } from "@playwright/test";

// API テストのみ（request API）なのでブラウザは使わない
export default defineConfig({
  testDir: "./tests",
  // WAF の geo 誤判定など一過性の失敗に備えてリトライ
  retries: 2,
  reporter: "list",
  use: {
    baseURL: process.env.API_BASE_URL ?? "http://localhost:8080",
  },
});
