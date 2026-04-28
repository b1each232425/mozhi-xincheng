import cv2
import os
from pathlib import Path

def batch_extract_frames(input_dir, output_dir, frames_per_video=10):
    """
    遍历指定目录下的所有视频，并从中均匀抽取固定数量的帧。
    """
    # 1. 路径预处理
    input_path = Path(input_dir).resolve()
    output_path = Path(output_dir).resolve()

    # 支持的视频格式
    valid_extensions = ('.mp4', '.avi', '.mov', '.mkv', '.flv')

    # 确保输出目录存在
    if not output_path.exists():
        output_path.mkdir(parents=True, exist_ok=True)
        print(f"✅ 已创建输出目录: {output_path}")

    # 获取所有视频文件
    video_files = [f for f in os.listdir(input_path) if f.lower().endswith(valid_extensions)]

    if not video_files:
        print(f"错误: 在 {input_path} 中没找到任何视频文件！")
        print(f"请检查路径是否存在视频: {input_path}")
        return

    print(f"找到 {len(video_files)} 个视频，准备开始抽帧...")

    total_saved = 0
    for video_file in video_files:
        video_full_path = str(input_path / video_file)
        video_name = Path(video_file).stem

        cap = cv2.VideoCapture(video_full_path)
        if not cap.isOpened():
            print(f"⚠无法打开视频: {video_file}")
            continue

        # 获取总帧数
        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        if total_frames <= frames_per_video:
            print(f"视频 {video_file} 太短，跳过")
            cap.release()
            continue

        # 逻辑：均匀采样，避开视频开头和结尾 10% 的可能存在的转场或黑屏
        start_frame = int(total_frames * 0.1)
        end_frame = int(total_frames * 0.9)
        interval = (end_frame - start_frame) // frames_per_video

        saved_count = 0
        for i in range(frames_per_video):
            # 计算目标帧位置
            target_pos = start_frame + (i * interval)
            cap.set(cv2.CAP_PROP_POS_FRAMES, target_pos)

            ret, frame = cap.read()
            if ret:
                # 统一保存为 jpg，命名包含原视频名以防混淆
                img_filename = f"训练集{i}_idx{i:02d}.jpg"
                save_path = str(output_path / img_filename)
                cv2.imwrite(save_path, frame)
                saved_count += 1
                total_saved += 1
            else:
                break

        cap.release()
        print(f"  - 完成: {video_file} (提取 {saved_count} 帧)")

    print(f"\n 全部抽帧任务完成！")
    print(f"总计提取图片: {total_saved} 张")
    print(f"保存路径: {output_path}")

if __name__ == "__main__":
    # --- 配置区域 ---
    # 假设你的项目结构：
    # moxin-zhicheng_backend/
    # ├── uploads/videos/ (视频放这里)
    # ├── datasets/raw_images/ (图片会存到这里)
    # └── vision_service/extract_all_videos.py (脚本位置)

    # 使用相对路径自动推导
    BASE_DIR = Path(__file__).resolve().parent.parent
    INPUT_DIR = BASE_DIR / "uploads" / "videos"
    OUTPUT_DIR = BASE_DIR / "datasets" / "raw_images"

    # 执行
    batch_extract_frames(INPUT_DIR, OUTPUT_DIR, frames_per_video=10)