import { readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";

/**
 * Renders the original Gin templates unchanged. They only use `{{ define }}`, `{{ template "x" . }}`,
 * `{{ .title }}` and `{{ .key }}`, so a tiny renderer covers them; anything else is rejected at load time.
 */
export interface PageData {
  title: string;
  key?: string;
}

const DEFINE_RE = /\{\{\s*define\s+"([^"]+)"\s*\}\}([\s\S]*?)\{\{\s*end\s*\}\}/g;
const DIRECTIVE_RE = /\{\{\s*(?:template\s+"([^"]+)"\s*\.|\.(title|key))\s*\}\}/g;

export function escapeHtml(s: string): string {
  return s.replace(/[&<>"']/g, (c) => `&#${c.charCodeAt(0)};`);
}

export function loadTemplates(dir: string): Map<string, string> {
  const defs = new Map<string, string>();
  for (const sub of readdirSync(dir, { withFileTypes: true })) {
    if (!sub.isDirectory()) continue;
    for (const f of readdirSync(join(dir, sub.name))) {
      if (!f.endsWith(".html")) continue;
      const src = readFileSync(join(dir, sub.name, f), "utf8");
      for (const m of src.matchAll(DEFINE_RE)) defs.set(m[1], m[2]);
    }
  }
  for (const [name, body] of defs) {
    const leftover = body.replace(DIRECTIVE_RE, "").match(/\{\{[\s\S]*?\}\}/);
    if (leftover) throw new Error(`template ${name}: unsupported directive ${leftover[0]}`);
  }
  return defs;
}

export function render(defs: Map<string, string>, name: string, data: PageData, depth = 0): string {
  const body = defs.get(name);
  if (body === undefined) throw new Error(`template not found: ${name}`);
  if (depth > 10) throw new Error(`template recursion: ${name}`);
  return body.replace(DIRECTIVE_RE, (_, tpl: string | undefined, field: "title" | "key" | undefined) => {
    if (tpl) return render(defs, tpl, data, depth + 1);
    // `key` sits in a script src query string, `title` in element text.
    return field === "key" ? encodeURIComponent(data.key ?? "") : escapeHtml(data.title);
  });
}
