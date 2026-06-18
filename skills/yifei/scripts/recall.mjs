#!/usr/bin/env node
'use strict';
// yifei-cli 经验记忆召回脚本。
// 按 frontmatter `aliases` 正则匹配 ~/.config/yifei-cli/memory/*.md 并输出正文。
//   node recall.mjs "<查询文本>"   -> 打印命中经验正文
//   node recall.mjs --list          -> 列出所有经验主题 + aliases
// 仅用 node 内置模块,零依赖,跨平台;无匹配/无目录静默退出 0。
import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import { pathToFileURL } from 'node:url';

export function memoryDir() {
  if (process.env.YIFEI_MEMORY_DIR) return process.env.YIFEI_MEMORY_DIR;
  if (process.platform === 'win32') {
    const base = process.env.APPDATA || path.join(os.homedir(), 'AppData', 'Roaming');
    return path.join(base, 'yifei-cli', 'memory');
  }
  const base = process.env.XDG_CONFIG_HOME || path.join(os.homedir(), '.config');
  return path.join(base, 'yifei-cli', 'memory');
}

export function parseFrontmatter(raw) {
  const m = raw.match(/^---\s*\r?\n([\s\S]*?)\r?\n---\s*\r?\n?/);
  if (!m) return { aliases: [], body: raw };
  const fm = m[1];
  const body = raw.slice(m[0].length);
  const aliasLine = fm.split(/\r?\n/).find((l) => l.trim().startsWith('aliases:')) || '';
  const aliases = aliasLine
    .replace(/^\s*aliases:\s*/, '')
    .replace(/^\[/, '')
    .replace(/\]$/, '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean);
  return { aliases, body };
}

function readMd(dir) {
  if (!fs.existsSync(dir)) return [];
  return fs
    .readdirSync(dir, { withFileTypes: true })
    .filter((e) => e.isFile() && e.name.endsWith('.md'))
    .map((e) => {
      const topic = e.name.replace(/\.md$/, '');
      const raw = fs.readFileSync(path.join(dir, e.name), 'utf8');
      return { topic, raw };
    })
    .sort((a, b) => a.topic.localeCompare(b.topic));
}

export function listTopics(dir) {
  return readMd(dir).map(({ topic, raw }) => ({ topic, aliases: parseFrontmatter(raw).aliases }));
}

function escapeRe(s) {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

export function recall(dir, query) {
  if (!query) return [];
  const out = [];
  for (const { topic, raw } of readMd(dir)) {
    const { aliases, body } = parseFrontmatter(raw);
    const terms = [topic, ...aliases].filter(Boolean).map(escapeRe);
    if (!terms.length) continue;
    if (new RegExp(terms.join('|'), 'i').test(query)) {
      out.push({ topic, body: body.trimEnd() });
    }
  }
  return out;
}

function main() {
  const arg = (process.argv[2] || '').trim();
  const dir = memoryDir();
  if (arg === '--list') {
    for (const t of listTopics(dir)) {
      process.stdout.write(`${t.topic}: ${t.aliases.join(', ')}\n`);
    }
    return;
  }
  if (!arg) return;
  for (const hit of recall(dir, arg)) {
    process.stdout.write(`--- 经验: ${hit.topic} ---\n${hit.body}\n\n`);
  }
}

if (import.meta.url === pathToFileURL(process.argv[1] || '').href) {
  main();
}
