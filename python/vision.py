import cv2
import numpy as np
import time
import pyautogui


class Vision:
    def __init__(self):
        self.frame = None

    @staticmethod
    def take_screenshot():
        image = pyautogui.screenshot()

        image = cv2.cvtColor(np.array(image), cv2.COLOR_RGB2BGR)
        img_gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

        return img_gray

    def refresh_frame(self):
        self.frame = self.take_screenshot()

    def crop_frame(self, x, y, w, h):
        img = self.frame.copy()
        return img[y:y + h, x:x + w]
