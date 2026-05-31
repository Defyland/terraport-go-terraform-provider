import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 10,
  duration: "30s",
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
  sleep(1);
}
