import shutil
import os
import sys

# fix_config.py 位于 .system/openclaw_rescue/
# 配置源位于 .system/openclaw_normal/
_rescue_dir = os.path.dirname(os.path.abspath(__file__))
_system_dir = os.path.dirname(_rescue_dir)
_normal_dir = _system_dir

config_src = os.path.join(_normal_dir, "openclaw_normal", "config_default.json")
config_dst = os.path.join(_normal_dir, "openclaw_normal", "config.json")

if not os.path.exists(config_src):
    print(f"[错误] 出厂备份不存在: {config_src}", file=sys.stderr)
    sys.exit(1)

if not os.path.exists(config_dst) or os.path.getsize(config_dst) == 0:
    shutil.copy(config_src, config_dst)
    print(f"[修复] 日常版配置已重置: {config_dst}")
else:
    print(f"[OK] 日常版配置文件正常，无需修复")