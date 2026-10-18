# 回滚指南

## 自动回退机制
当新版本解析失败时，tars-update 自动回退到 `stable_versions.json` 中标记的稳定版本。

## 手动回滚
```bash
tars-update --rollback-to v1.0.1
```

## 验证
- 回滚后执行 `--version` 确认
- 检查 SHA256SUMS 校验通过
