import os
from ultralytics import YOLO

def train_fretboard_model():
    # 获取当前脚本所在目录 (vision_service)
    current_dir = os.path.dirname(os.path.abspath(__file__))
    data_yaml_path = os.path.join(current_dir, "../datasets/guitar_fretboard/data.yaml")

    if not os.path.exists(data_yaml_path):
        print(f"错误：找不到配置文件 {data_yaml_path}，请确认数据集解压路径！")
        return

    # 2. 初始化 YOLOv11 姿态估计模型
    # 使用 'yolo11n-pose.pt' (Nano版本) 保证在 CPU 服务器上的推理速度
    model = YOLO('yolo11n-pose.pt')

    model.train(
        data=data_yaml_path,    # 指向你的 data.yaml
        epochs=100,             # 训练 100 轮以达到 98% 左右的精度
        imgsz=640,              # 输入图像尺寸
        batch=8,               # RTX 4050 显存建议设为 16
        device=0,               # 强制使用第一块 GPU 训练
        project='runs/train',   # 训练结果保存路径
        name='guitar_fretboard_v3', # 实验名称
        save=True,              # 自动保存 best.pt
        exist_ok=True,          # 允许覆盖同名文件夹
        optimizer='auto',       # 自动选择优化器 (通常为 AdamW 或 SGD)
        lr0=0.01,               # 初始学习率
        patience=20,            # 如果 20 轮内精度不提升则提前停止，防止过拟合
        augment=True,            # 开启数据增强（翻转、亮度等）
        workers=2
    )

    print("训练完成！模型保存在: vision_service/runs/train/guitar_fretboard_v3/weights/best.pt")

if __name__ == "__main__":
    train_fretboard_model()