import cv2
import numpy as np
import time


class Game:
    def __init__(self, vision, controller):
        self.vision = vision
        self.controller = controller
        self.state = 'not started'

    def run(self):
        cv2.namedWindow("screen", cv2.WINDOW_NORMAL)
        cv2.setMouseCallback('screen', self.mousePosition)

        while True:
            self.vision.refresh_frame()

            crop_img = self.vision.crop_frame(1195, 1265, 1400, 205)

            cv2.imshow("screen", self.vision.frame)
            cv2.imshow("crop_img", crop_img)

            ch = cv2.waitKey(0)
            if ch == 27 or ch == ord('q'):
                break

    def log(self, text):
        print('[%s] %s' % (time.strftime('%H:%M:%S'), text))

    @staticmethod
    def mousePosition(event, x, y, flags, param):
        print(x, y)
