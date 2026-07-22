# FileFlow

aria2 下载完成后的文件自动整理工具。将下载的文件从 SSD 移动到 HDD，再根据规则匹配后重命名并移动到最终目录。

## 使用方式

### 1. 编译

```bash
go build -o FileFlow .
```

### 2. 放置文件

将 `FileFlow`（二进制）和 `fileflow.yaml`（配置文件）放在同一目录下。

### 3. 配置 aria2

在 aria2 的配置文件中添加：

```
on-download-complete=/path/to/FileFlow
```

aria2 下载完成后会自动调用 FileFlow，传入参数：`GID file_count first_file_path`。

### 4. 配置规则

编辑 `fileflow.yaml`，定义文件匹配和移动规则（见下方）。

## 配置文件

### 文件位置

`fileflow.yaml` 与 `FileFlow` 二进制在同一目录。

### 配置格式

```yaml
rules:
  - name: "规则名称"
    pattern: "文件名匹配模式（glob）"
    target_dir: "目标目录"
    target_name: "目标文件名模板"
```

### 占位符

`target_name` 支持以下占位符：

| 占位符 | 说明 | 示例 |
|--------|------|------|
| `{source.0}` | 完整的源文件名 | `Rick.and.Morty.S09E09.1080p.x265-ELiTE.mkv` |
| `{source.1}` | 第一个通配符 `*` 匹配的内容 | `09` |
| `{source.2}` | 第二个通配符 `*` 匹配的内容 | ... |
| `{target.index}` | 目标目录中已存在文件的递增序号（01, 02...） | `03` |

### 配置示例

**场景一：剧集整理（保留原标题）**

```yaml
rules:
  - name: "瑞克和莫蒂 S09"
    pattern: "Rick.and.Morty.S09E*.1080p.x265-ELiTE.mkv"
    target_dir: "/mnt/ChaseShows/Rick.and.Morty.S09"
    target_name: "{source.0}"
```

匹配过程：

```
源文件: Rick.and.Morty.S09E09.1080p.x265-ELiTE.mkv
pattern:              *
         Rick.and.Morty.S09E09.1080p.x265-ELiTE.mkv
         ↓ source.0 = "Rick.and.Morty.S09E09.1080p.x265-ELiTE.mkv"
目标: /mnt/ChaseShows/Rick.and.Morty.S09/Rick.and.Morty.S09E09.1080p.x265-ELiTE.mkv
```


**场景二：使用目标目录序号**

```yaml
rules:
  - name: "剧集 - 自动编号"
    pattern: "*.mkv"
    target_dir: "/data/tv/unsorted"
    target_name: "Episode {target.index}.{ext}"
```

如果 `/data/tv/unsorted` 已有 `Episode 01.mkv`、`Episode 02.mkv`，下一个文件会被命名为 `Episode 03.mkv`。

**场景四：多个规则按优先级**

```yaml
rules:
  - name: "指定剧集"
    pattern: "Rick.and.Morty.S09E*.mkv"
    target_dir: "/mnt/ChaseShows/Rick.and.Morty.S09"
    target_name: "Rick.and.Morty.S09E{source.1}.mkv"

  - name: "其它视频"
    pattern: "*.mkv"
    target_dir: "/data/tv/unsorted"
    target_name: "{source.0}"
```

规则按定义顺序匹配，第一条命中即生效。`Rick.and.Morty` 会被第一条规则捕获，不会落到第二条。

## 工作流程

```
aria2 下载完成
  → FileFlow 被调用
  → 从 SSD 移动到 HDD（跨设备自动使用 copy+delete）
  → 加载 fileflow.yaml 规则
  → 判断 HDD 上的是文件还是目录
     ├── 文件 → 匹配规则 → 重命名 → 移动到 target_dir
     └── 目录 → 递归所有文件 → 逐个匹配 → 逐个移动
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `FILEFLOW_CONFIG` | 配置文件路径 | 与二进制同目录的 `fileflow.yaml` |
