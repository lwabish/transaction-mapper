# CI/CD 流程优化迁移指南

## 概述

本指南说明如何从旧的 CI/CD 流程（生成 YAML 文件并提交到仓库）迁移到新的 Helm + ArgoCD 流程（无需代码提交）。

## 主要变化

### 1. 旧流程问题
- CI 每次构建后都会生成 `deploy/backend.yml` 和 `deploy/web.yml`
- CI 会自动提交这些文件到仓库
- 本地需要频繁 `git pull` 同步 CI 的提交
- 代码仓库包含大量生成的部署文件

### 2. 新流程优势
- ✅ CI 只负责构建和推送镜像
- ✅ ArgoCD 自动检测新镜像并更新部署
- ✅ 无需代码提交，保持代码仓库清洁
- ✅ 支持多环境配置（prod/staging）
- ✅ 使用 Helm 模板化配置管理

## 新架构

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    CI/CD    │    │   Docker    │    │   ArgoCD    │
│  (GitHub)   │───▶│   Hub       │───▶│             │
│             │    │             │    │             │
│ - 构建镜像   │    │ - 存储镜像   │    │ - 检测新镜像 │
│ - 推送镜像   │    │ - 标签管理   │    │ - 自动部署   │
└─────────────┘    └─────────────┘    └─────────────┘
                                             │
                                             ▼
                                    ┌─────────────┐
                                    │ Kubernetes  │
                                    │  Cluster    │
                                    │             │
                                    │ - Helm 部署 │
                                    │ - 滚动更新   │
                                    └─────────────┘
```

## 文件结构变化

### 新增文件
```
helm/
└── transaction-mapper/
    ├── Chart.yaml                    # Helm Chart 定义
    ├── values.yaml                   # 默认配置
    ├── templates/                    # K8S 资源模板
    │   ├── _helpers.tpl             # 模板辅助函数
    │   ├── backend-deployment.yaml  # 后端部署模板
    │   ├── backend-service.yaml     # 后端服务模板
    │   ├── web-deployment.yaml      # 前端部署模板
    │   ├── web-service.yaml         # 前端服务模板
    │   └── ingress.yaml             # Ingress 模板
    └── environments/                 # 环境配置
        ├── prod/                    # 生产环境
        │   └── kustomization.yaml   # Kustomize 配置
        └── staging/                 # 测试环境
            └── kustomization.yaml   # Kustomize 配置

argocd/
├── transaction-mapper-prod.yaml      # 生产环境 ArgoCD 应用
├── transaction-mapper-staging.yaml   # 测试环境 ArgoCD 应用
└── image-updater-config.yaml         # 镜像自动更新配置
```

### 废弃文件（可删除）
```
deploy/
├── backend-template.yml   # 可删除
├── web-template.yml       # 可删除
├── backend.yml           # 可删除
└── web.yml               # 可删除
```

## 迁移步骤

### 1. 部署 Helm Chart 和 ArgoCD

```bash
# 1. 创建命名空间
kubectl create namespace transaction-mapper
kubectl create namespace transaction-mapper-staging

# 2. 安装 ArgoCD 应用
kubectl apply -f argocd/transaction-mapper-prod.yaml
kubectl apply -f argocd/transaction-mapper-staging.yaml

# 3. 可选：配置 ArgoCD Image Updater
kubectl apply -f argocd/image-updater-config.yaml
```

### 2. 验证部署

```bash
# 检查 ArgoCD 应用状态
kubectl get applications -n argocd

# 检查部署状态
kubectl get pods -n transaction-mapper
kubectl get pods -n transaction-mapper-staging
```

### 3. 测试新流程

```bash
# 创建新标签触发 CI 构建
git tag v0.1.0
git push origin v0.1.0

# 观察 ArgoCD 自动更新部署
kubectl argocd app get transaction-mapper-prod -w
```

## 配置说明

### CI 配置变化

- **权限**: 从 `contents: write` 改为 `contents: read`（不再需要提交代码）
- **输出**: 只构建和推送镜像，不生成 YAML 文件
- **通知**: 输出镜像信息供 ArgoCD 使用

### ArgoCD 配置

- **源**: 使用 Helm Chart + Kustomize
- **镜像更新**: 支持自动检测新镜像标签
- **多环境**: 分别管理 prod 和 staging 环境

### Image Updater 配置

- **自动更新**: 检测新镜像标签并自动部署
- **语义版本**: 支持 semver 版本约束
- **Git 回写**: 可选将镜像更新提交到 Git

## 回滚方案

如果需要回滚到旧流程：

1. **保留旧的 deploy 文件**：暂时保留 `deploy/` 目录文件
2. **更新 ArgoCD 应用**：将 source.path 改回 `deploy/`
3. **恢复 CI 配置**：恢复生成 YAML 文件的步骤

## 故障排除

### 常见问题

1. **ArgoCD 无法检测新镜像**
   - 检查镜像仓库访问权限
   - 验证镜像标签格式是否符合配置
   - 查看 ArgoCD Image Updater 日志

2. **Helm Chart 渲染错误**
   - 检查模板语法
   - 验证 values.yaml 配置
   - 查看本地渲染：`helm template ./helm/transaction-mapper`

3. **权限问题**
   - 确保 ArgoCD 有访问目标命名空间的权限
   - 检查 service account 配置

### 调试命令

```bash
# 本地测试 Helm Chart
helm template transaction-mapper-prod ./helm/transaction-mapper \
  --values ./helm/transaction-mapper/values.yaml

# 测试 Kustomize
kubectl kustomize ./helm/transaction-mapper/environments/prod

# 查看 ArgoCD 同步状态
kubectl argocd app get transaction-mapper-prod
kubectl argocd app history transaction-mapper-prod

# 手动同步 ArgoCD 应用
kubectl argocd app sync transaction-mapper-prod
```

## 最佳实践

1. **版本管理**: 使用语义版本标签（如 `v1.2.3`）
2. **环境隔离**: 不同环境使用独立的命名空间和配置
3. **渐进部署**: 先部署到 staging，验证后再部署到 prod
4. **监控告警**: 配置部署状态和健康检查告警
5. **备份策略**: 定期备份 ArgoCD 配置和 Helm Chart