# iaa 

一个用于快速构建 Go Web 项目框架的小工具。

## 安装

```bash
go install github.com/q1mi/iaa@latest
```

## 使用

### 基本使用

```bash
# 创建基础项目（默认）
iaa new project_name
```

### 高级使用

```bash
# 使用进阶模板（包含依赖注入等高级特性）
iaa new project_name --advanced

# 使用自定义模板仓库
iaa new project_name --repo https://github.com/your/custom-template.git

# 查看所有可用参数
iaa new --help
```

## 参数说明

- `--advanced`: 使用进阶模板，包含依赖注入、更完善的项目结构
- `--repo <url>`: 指定自定义模板仓库URL（优先级最高）

**参数优先级**: `--repo` > `--advanced` > 默认基础模板

## 模板类型

### 基础模板
最简单朴素的基于gin的web项目框架，全部依赖均使用全局变量，没有使用依赖注入。

- 仓库：[gin-base-layout](https://github.com/q1mi/gin-base-layout)
- 适用：快速原型开发、简单项目

### 进阶模板  
包含依赖注入、Wire框架、更完善的项目架构的gin项目模板。

- 仓库：[gin-advanced-layout](https://github.com/q1mi/gin-advanced-layout) 
- 适用：企业级开发、复杂项目

### 自定义模板
支持使用任何兼容的Git仓库作为项目模板。

## 使用示例

```bash
# 创建基础项目
iaa new my-api

# 创建进阶项目
iaa new my-enterprise-api --advanced

# 使用公司内部模板
iaa new my-project --repo https://github.com/mycompany/gin-template.git

# 使用GitHub上的其他模板
iaa new my-project --repo https://github.com/someone/awesome-gin-template.git
```
