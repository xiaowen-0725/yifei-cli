import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';

import { parseFrontmatter, listTopics, recall } from './recall.mjs';

const SCRIPT = fileURLToPath(new URL('./recall.mjs', import.meta.url));

function tmpMemory(files) {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'yfmem-'));
  for (const [name, content] of Object.entries(files)) {
    fs.writeFileSync(path.join(dir, name), content);
  }
  return dir;
}

const COPTC = `---
aliases: [COPTC, 销售订单, 销售订单单头]
tables: [COPTC, COPTD]
updated: 2026-06-18
---

## 单别含义
- \`221\` = 内销订单 (2026-06-18)
`;

test('parseFrontmatter extracts aliases and strips frontmatter', () => {
  const { aliases, body } = parseFrontmatter(COPTC);
  assert.deepEqual(aliases, ['COPTC', '销售订单', '销售订单单头']);
  assert.ok(body.startsWith('## 单别含义'));
  assert.ok(!body.includes('aliases:'));
});

test('recall matches by alias case-insensitively and returns body', () => {
  const dir = tmpMemory({ 'COPTC.md': COPTC });
  const hits = recall(dir, '帮我查 coptc 的单别');
  assert.equal(hits.length, 1);
  assert.equal(hits[0].topic, 'COPTC');
  assert.ok(hits[0].body.includes('内销订单'));
  assert.ok(!hits[0].body.includes('aliases:'));
});

test('recall matches by Chinese alias', () => {
  const dir = tmpMemory({ 'COPTC.md': COPTC });
  assert.equal(recall(dir, '销售订单怎么算金额').length, 1);
});

test('recall returns empty on no match', () => {
  const dir = tmpMemory({ 'COPTC.md': COPTC });
  assert.deepEqual(recall(dir, '完全无关的词'), []);
});

test('recall returns empty when dir missing', () => {
  assert.deepEqual(recall(path.join(os.tmpdir(), 'nope-yfmem-xyz'), 'COPTC'), []);
});

test('listTopics lists topics with aliases sorted', () => {
  const dir = tmpMemory({ 'COPTC.md': COPTC, 'INVMB.md': '---\naliases: [INVMB, 品号]\n---\n\n## 字段枚举\n- x\n' });
  const topics = listTopics(dir);
  assert.deepEqual(topics.map((t) => t.topic), ['COPTC', 'INVMB']);
  assert.deepEqual(topics[1].aliases, ['INVMB', '品号']);
});

test('CLI prints matching body and is silent on no match', () => {
  const dir = tmpMemory({ 'COPTC.md': COPTC });
  const env = { ...process.env, YIFEI_MEMORY_DIR: dir };
  const hit = execFileSync('node', [SCRIPT, '销售订单'], { env, encoding: 'utf8' });
  assert.ok(hit.includes('--- 经验: COPTC ---'));
  assert.ok(hit.includes('内销订单'));
  const miss = execFileSync('node', [SCRIPT, 'zzz无关'], { env, encoding: 'utf8' });
  assert.equal(miss.trim(), '');
});

test('CLI --list prints topics', () => {
  const dir = tmpMemory({ 'COPTC.md': COPTC });
  const env = { ...process.env, YIFEI_MEMORY_DIR: dir };
  const out = execFileSync('node', [SCRIPT, '--list'], { env, encoding: 'utf8' });
  assert.ok(out.includes('COPTC'));
  assert.ok(out.includes('销售订单'));
});

test('CLI exits 0 when memory dir does not exist', () => {
  const env = { ...process.env, YIFEI_MEMORY_DIR: path.join(os.tmpdir(), 'definitely-missing-yfmem') };
  const out = execFileSync('node', [SCRIPT, 'COPTC'], { env, encoding: 'utf8' });
  assert.equal(out.trim(), '');
});
