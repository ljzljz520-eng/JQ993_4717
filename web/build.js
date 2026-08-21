const fs=require('fs');
fs.mkdirSync('dist',{recursive:true});
fs.writeFileSync('dist/index.html','<!doctype html><html><head><meta charset="utf-8"><title>香薰产品维护台</title></head><body><main><h1>香薰产品维护台</h1><p>请使用本地 API 管理产品资料。</p></main></body></html>');
