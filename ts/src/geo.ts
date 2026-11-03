const EARTH_RADIUS_KM = 6378.1;

const toRad = (d: number) => (d * Math.PI) / 180;
const toDeg = (r: number) => (r * 180) / Math.PI;

export interface LatLng {
  lat: number;
  lng: number;
}

/** Point `km` away from (lat, lng) along `bearingDeg`. */
export function destination(lat: number, lng: number, km: number, bearingDeg: number): LatLng {
  const φ1 = toRad(lat);
  const λ1 = toRad(lng);
  const θ = toRad(bearingDeg);
  const δ = km / EARTH_RADIUS_KM;
  const φ2 = Math.asin(Math.sin(φ1) * Math.cos(δ) + Math.cos(φ1) * Math.sin(δ) * Math.cos(θ));
  const λ2 = λ1 + Math.atan2(Math.sin(θ) * Math.sin(δ) * Math.cos(φ1), Math.cos(δ) - Math.sin(φ1) * Math.sin(φ2));
  return { lat: toDeg(φ2), lng: toDeg(λ2) };
}

/**
 * The four `km`×`km` quadrants around a point, as Kakao `rect` strings ("x1,y1,x2,y2" = lng,lat,lng,lat).
 * Same four rectangles the Go RectSearch built from its N/E/S/W points.
 */
export function quadrantRects(lat: number, lng: number, km: number): string[] {
  const n = destination(lat, lng, km, 0);
  const e = destination(lat, lng, km, 90);
  const s = destination(lat, lng, km, 180);
  const w = destination(lat, lng, km, 270);
  const rect = (a: LatLng, b: LatLng) => [a.lng, a.lat, b.lng, b.lat].join(",");
  return [rect(e, s), rect(e, n), rect(w, s), rect(w, n)];
}

/** Haversine distance in whole metres. Arguments are Kakao-style strings: x = longitude, y = latitude. */
export function distanceMeters(x1: string, y1: string, x2: string, y2: string): number {
  const lon1 = parseFloat(x1) || 0;
  const lat1 = parseFloat(y1) || 0;
  const lon2 = parseFloat(x2) || 0;
  const lat2 = parseFloat(y2) || 0;
  const dLat = toRad(lat2 - lat1);
  const dLon = toRad(lon2 - lon1);
  const a = Math.sin(dLat / 2) ** 2 + Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) ** 2;
  return Math.trunc(EARTH_RADIUS_KM * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a)) * 1000);
}
