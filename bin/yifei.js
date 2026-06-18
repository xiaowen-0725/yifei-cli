#!/usr/bin/env node
'use strict';
// Launcher shim: execs the platform binary downloaded by scripts/postinstall.js.
const { spawnSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const binName = process.platform === 'win32' ? 'yifei.exe' : 'yifei';
const bin = path.join(__dirname, '..', 'binaries', binName);

if (!fs.existsSync(bin)) {
  console.error('yifei: 二进制未找到。请重新安装: npm install -g yifei-cli');
  console.error('或从 https://github.com/xiaowen-0725/yifei-cli/releases 手动下载。');
  process.exit(1);
}

const res = spawnSync(bin, process.argv.slice(2), { stdio: 'inherit' });
if (res.error) {
  console.error('yifei:', res.error.message);
  process.exit(1);
}
process.exit(res.status === null ? 1 : res.status);
