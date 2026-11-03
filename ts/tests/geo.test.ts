import { test } from "node:test";
import assert from "node:assert/strict";
import { destination, distanceMeters, quadrantRects } from "../src/geo.js";
import { categoryLabel } from "../src/kakao.js";

test("distanceMeters: Seoul City Hall to Gangnam Station is about 8.7 km", () => {
  const d = distanceMeters("126.9780", "37.5665", "127.0276", "37.4979");
  assert.ok(d > 8500 && d < 8900, String(d));
  assert.equal(distanceMeters("127", "37", "127", "37"), 0);
});

test("distanceMeters uses both latitudes (Go version reused lat2)", () => {
  const d = distanceMeters("127", "37.0", "127", "37.01");
  assert.ok(d > 1100 && d < 1120, String(d));
});

test("destination moves along the bearing", () => {
  const n = destination(37.5, 127, 1, 0);
  assert.ok(n.lat > 37.5 && Math.abs(n.lng - 127) < 1e-9);
  const e = destination(37.5, 127, 1, 90);
  assert.ok(e.lng > 127 && Math.abs(e.lat - 37.5) < 1e-4);
});

test("quadrantRects covers the four quadrants around the point", () => {
  const rects = quadrantRects(37.5, 127, 0.5).map((r) => r.split(",").map(Number));
  assert.equal(rects.length, 4);
  for (const [x1, y1, x2, y2] of rects) {
    // each rect has the centre as one corner coordinate on each axis
    assert.ok(Math.min(Math.abs(x1 - 127), Math.abs(x2 - 127)) < 1e-9);
    assert.ok(Math.min(Math.abs(y1 - 37.5), Math.abs(y2 - 37.5)) < 1e-4);
  }
});

test("category labels match the Go parseCategory table", () => {
  assert.equal(categoryLabel("korea"), "한식");
  assert.equal(categoryLabel("meat"), "육류,고기");
  assert.equal(categoryLabel("anything"), null);
  assert.equal(categoryLabel("FD6"), null);
  assert.equal(categoryLabel(undefined), null);
});
