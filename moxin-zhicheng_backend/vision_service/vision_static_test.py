import cv2
import mediapipe as mp
import numpy as np
from ultralytics import YOLO  # 引入你训练好的 YOLO

# 1. 初始化模型
# 请确保 best.pt 放在脚本同级目录下
yolo_model = YOLO('best.pt')

mp_hands = mp.solutions.hands
hands = mp_hands.Hands(static_image_mode=True, max_num_hands=1, min_detection_confidence=0.5)
mp_draw = mp.solutions.drawing_utils

def process_guitar_with_yolo(image_path):
    # 读取图片
    image = cv2.imread(image_path)
    if image is None:
        print(f"错误：无法读取图片: {image_path}")
        return

    h, w, _ = image.shape

    # --- 步骤 1: 使用你的 YOLOv11 自动识别指板 4 个顶点 ---
    results_yolo = yolo_model(image)
    src_pts = None

    for r in results_yolo:
        if r.keypoints is not None and len(r.keypoints.xy) > 0:
            # 提取 4 个关键点坐标 (x, y)
            # 顺序通常是你标注时的顺序：左上, 右上, 右下, 左下
            src_pts = r.keypoints.xy[0].cpu().numpy().astype(np.float32)
            break

    if src_pts is None or len(src_pts) < 4:
        print("警告：YOLO 未能完整识别出指板的 4 个顶点，请检查图片或模型！")
        return

    # --- 步骤 2: 透视变换（拉直指板） ---
    dst_w, dst_h = 1200, 300
    dst_pts = np.array([[0, 0], [dst_w, 0], [dst_w, dst_h], [0, dst_h]], dtype=np.float32)
    M = cv2.getPerspectiveTransform(src_pts, dst_pts)
    warped_fretboard = cv2.warpPerspective(image, M, (dst_w, dst_h))

    # --- 步骤 3: MediaPipe 手部识别 ---
    image_rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    results_hands = hands.process(image_rgb)

    if results_hands.multi_hand_landmarks:
        for hand_landmarks in results_hands.multi_hand_landmarks:
            # 绘制手部骨架用于可视化
            mp_draw.draw_landmarks(image, hand_landmarks, mp_hands.HAND_CONNECTIONS)

            # 获取食指指尖坐标
            index_tip = hand_landmarks.landmark[mp_hands.HandLandmark.INDEX_FINGER_TIP]
            px, py = int(index_tip.x * w), int(index_tip.y * h)

            # --- 步骤 4: 坐标转换（核心业务逻辑） ---
            finger_point = np.array([[[px, py]]], dtype=np.float32)
            transformed_point = cv2.perspectiveTransform(finger_point, M)
            tx, ty = transformed_point[0][0]

            # 判定手指是否在指板范围内
            if 0 <= tx <= dst_w and 0 <= ty <= dst_h:
                string_num = int(ty / (dst_h / 6)) + 1
                fret_num = calculate_fret(tx, dst_w)

                # 可视化：在拉直后的图上绘制识别结果
                cv2.putText(warped_fretboard, f"String: {string_num} Fret: {fret_num}",
                            (int(tx), int(ty)-10), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 0), 2)
                cv2.circle(warped_fretboard, (int(tx), int(ty)), 10, (0, 0, 255), -1)
                print(f"识别结果：第 {string_num} 弦, 第 {fret_num} 品")

    # --- 步骤 5: 结果展示 ---
    # 在原图上画出 YOLO 找到的 4 个点，方便你核对
    for i, pt in enumerate(src_pts):
        cv2.circle(image, (int(pt[0]), int(pt[1])), 8, (255, 0, 255), -1)
        cv2.putText(image, str(i), (int(pt[0]), int(pt[1])), cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 0, 255), 2)

    cv2.imshow("Step 1: AI Recognition (YOLO + MediaPipe)", cv2.resize(image, (w//2, h//2)))
    cv2.imshow("Step 2: Normalized Fretboard (Warped)", warped_fretboard)
    cv2.waitKey(0)
    cv2.destroyAllWindows()

def calculate_fret(tx, total_width):
    # (保持你原来的 17.817 物理公式逻辑)
    scale_length = total_width * 1.5
    current_pos, remaining_length = 0, scale_length
    for fret in range(1, 23):
        fret_width = remaining_length / 17.817
        if current_pos <= tx < (current_pos + fret_width):
            return fret
        current_pos += fret_width
        remaining_length -= fret_width
    return 0

if __name__ == "__main__":
    # 替换为你本地真实的吉他图片路径
    test_img = "your_guitar_image.jpg"
    process_guitar_with_yolo(test_img)