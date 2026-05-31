import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "10s", target: 1 },
    { duration: "10s", target: 100 },
    { duration: "10s", target: 1 },
  ],
};

const endpoint = __ENV.TERRAPORT_FAKE_API_ENDPOINT || "http://localhost:8080";
const token = __ENV.TERRAPORT_TOKEN || "test-token";

export default function () {
  const res = http.get(`${endpoint}/v1/products/bankport`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  check(res, {
    "product metadata returns 200": (r) => r.status === 200,
  });
}
