# tars-update

独立 Go module：Tars 系列项目的自动更新引擎。

## 关于作者

> 你好！我是这个项目的作者，一名正在读小学的学生。
> 这个项目从规划到代码实现，**完全由我自己主导设计，并借助 AI 工具完成**。每一行代码、每一个功能模块，都凝聚了我对技术的热爱与坚持。
> 如果你认可这个项目，或者觉得它对你有帮助，希望能给我一个 **Star** 支持，这对我意义重大！

## 欢迎支持

- 如果你在使用过程中遇到问题，欢迎提交 [Issue](https://github.com/HP-jh) 或参与讨论
- 如果你觉得项目不错，请点一个 **Star**，这是对我最大的鼓励
- 也欢迎分享给身边有需要的朋友，让更多人看到这个项目
- 后续功能迭代会继续推进，你的支持是我持续更新的动力

## ⚠️ 维护说明

> **本项目依赖 AI 服务套餐，若因套餐额度耗尽或服务商调整，可能导致更新暂时停滞。**
> 如遇较长时间未更新，请理解这是 AI 服务限制所致，并非项目放弃维护。
> 恢复服务后会第一时间继续推进，感谢你的耐心与包容。

---

## 功能

- **版本检查**：查询 GitHub Releases API，semver 比对
- **下载校验**：下载对应平台二进制 + SHA256SUMS.txt，哈希校验
- **原子替换**：rename 原子替换，旧二进制备份保留 7 天
- **失败回退**：GitHub API 不可用时回退到本地 `stable_versions.json`
- **跨平台**：linux / darwin / windows / android（amd64 / arm64）
- **多仓库**：同一模块支持 TarsSecureGuard / TarsWingGuard / TarsEnterprise 等 6 仓库

## 使用

```go
import "github.com/HP-jh/tars-update"

u := update.NewUpdater("twg", "v0.1.0", "HP-jh", "TarsWingGuard")

// Check
result := u.Check()
if result.HasUpdate {
    fmt.Printf("New version available: %s\n", result.Latest.String())
}

// Apply
if err := u.Apply(result); err != nil {
    log.Fatal(err)
}

// Rollback
if err := u.Rollback(); err != nil {
    log.Fatal(err)
}
```

## CLI 集成

在 `wingctl` 中集成：

```bash
wing update check    # 检查新版本
wing update apply    # 下载并应用更新
wing update rollback # 回退到上一版本
```

## 安全

- SHA256 校验必做
- 旧二进制备份保留，支持手动 rollback
- 升级日志写 `~/.tars/update.log`
- 极致档用户可配置 2FA 确认（由调用方实现）

## License

Apache-2.0

## Changelog

- v1.0.1: routine maintenance and reliability improvements
