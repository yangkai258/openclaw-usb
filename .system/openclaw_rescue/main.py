import os
import shutil
import sys

# 从当前脚本所在目录推算各路径
# E:\.system\openclaw_rescue\main.py
RESCUE_DIR = os.path.dirname(os.path.abspath(__file__))
SYSTEM_DIR = os.path.dirname(RESCUE_DIR)
NORMAL_DIR = os.path.join(SYSTEM_DIR, "openclaw_normal")

# 日常版用户配置（要修复的目标）
TARGET_CONFIG = os.path.join(NORMAL_DIR, "config.json")

# 出厂备份（在救援版目录自带）
DEFAULT_CONFIG = os.path.join(RESCUE_DIR, "config_default.json")

def repair_normal_config():
    """检查并修复日常版配置文件"""
    if not os.path.exists(DEFAULT_CONFIG):
        print(f"[错误] 出厂备份不存在: {DEFAULT_CONFIG}", file=sys.stderr)
        return False

    if not os.path.exists(TARGET_CONFIG) or os.path.getsize(TARGET_CONFIG) == 0:
        shutil.copy(DEFAULT_CONFIG, TARGET_CONFIG)
        print(f"[修复] 日常版配置已重置: {TARGET_CONFIG}")
        try:
            import ctypes
            ctypes.windll.user32.MessageBoxW(
                0,
                "软件配置已成功恢复为出厂默认状态！\n请重新双击运行主程序。",
                "OpenClaw 救援工具",
                0
            )
        except Exception:
            pass
        return True
    else:
        print(f"[OK] 日常版配置文件正常")
        try:
            import ctypes
            ctypes.windll.user32.MessageBoxW(
                0,
                "配置文件状态正常，无需修复。",
                "OpenClaw 救援工具",
                0
            )
        except Exception:
            pass
        return True

if __name__ == "__main__":
    repair_normal_config()
    # TODO: 接着启动救援版 UI / Agent ...