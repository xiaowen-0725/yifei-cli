#!/usr/bin/env node
'use strict';
// Downloads the correct platform binary from the GitHub release matching this
// package version, into ../binaries/. Soft-fails so a single unsupported
// platform never breaks `npm install`; the bin shim warns if the binary is absent.
const https = require('https');
const fs = require('fs');
const path = require('path');

const VERSION = require('../package.json').version;
const REPO = 'xiaowen-0725/yifei-cli';

const ASSETS = {
  'darwin-x64': 'yifei-darwin-amd64',
  'darwin-arm64': 'yifei-darwin-arm64',
  'linux-x64': 'yifei-linux-amd64',
  'linux-arm64': 'yifei-linux-arm64',
  'win32-x64': 'yifei-windows-amd64.exe',
  'win32-arm64': 'yifei-windows-arm64.exe',
};

const key = `${process.platform}-${process.arch}`;
const asset = ASSETS[key];
if (!asset) {
  console.error(`yifei-cli: 暂不支持的平台 ${key},跳过二进制下载。`);
  process.exit(0);
}

const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${asset}`;
const outDir = path.join(__dirname, '..', 'binaries');
fs.mkdirSync(outDir, { recursive: true });
const outFile = path.join(outDir, process.platform === 'win32' ? 'yifei.exe' : 'yifei');

function download(u, dest, cb, redirects) {
  redirects = redirects || 0;
  if (redirects > 10) return cb(new Error('重定向次数过多'));
  https
    .get(u, { headers: { 'User-Agent': 'yifei-cli-postinstall' } }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        res.resume();
        return download(res.headers.location, dest, cb, redirects + 1);
      }
      if (res.statusCode !== 200) {
        res.resume();
        return cb(new Error(`下载失败: HTTP ${res.statusCode}`));
      }
      const file = fs.createWriteStream(dest);
      res.pipe(file);
      file.on('finish', () => file.close(() => cb(null)));
      file.on('error', cb);
    })
    .on('error', cb);
}

console.log(`yifei-cli: 正在下载 ${asset} (v${VERSION})...`);
download(url, outFile, (err) => {
  if (err) {
    console.error(`yifei-cli: ${err.message}`);
    console.error(`可手动下载: https://github.com/${REPO}/releases/tag/v${VERSION}`);
    process.exit(0); // soft-fail: don't break npm install
  }
  if (process.platform !== 'win32') {
    try { fs.chmodSync(outFile, 0o755); } catch (e) { /* ignore */ }
  }
  console.log('yifei-cli: 安装完成。运行 `yifei --help` 开始使用。');
});
